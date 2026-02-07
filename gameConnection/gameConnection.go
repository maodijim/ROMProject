package gameConnection

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"google.golang.org/protobuf/proto"

	log "github.com/sirupsen/logrus"
)

type authJson struct {
	Status int    `json:"code"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

var (
	cmdQueueInterval     = 75 * time.Millisecond
	fixedItemCDSubtract  = uint64(40000) // in milliseconds
	fixedSkillCDSubtract = uint64(50000) // in milliseconds
	TradeProtoCmdId      = Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"]
	LogInUserProtoCmdId  = Cmd.Command_value["LOGIN_USER_PROTOCMD"]
)

// logWriter implements io.Writer interface
type logWriter struct {
	gc *GameConnection
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	message := string(p)
	lw.gc.AddLog(message)
	return len(p), nil
}

type chatMessage struct {
	MsgChannel Cmd.EGameChatChannel `json:"msgChannel"`
	SenderId   uint64               `json:"senderId"`
	SenderName string               `json:"senderName"`
	Content    string               `json:"content"`
	Timestamp  uint64               `json:"timestamp"`
	IsSent     bool                 `json:"isSent"`
}

type GameConnection struct {
	Authed             bool
	cmdQueue           [][]byte
	lastCmdSend        time.Time
	currentIndex       uint32
	ShouldHeartBeat    bool
	ShouldChangeScene  bool
	Configs            *config.ServerConfigs
	conn               net.Conn
	Role               *RoleInfo
	AvailableRoles     map[uint32]*RoleInfo
	DebugMsg           bool
	quitContext        context.Context
	quitCancel         context.CancelFunc
	shouldQuit         bool
	enteringMap        bool
	mails              []*Cmd.MailData
	inMap              bool
	tradeDetail        map[uint32]*Cmd.DetailPendingListRecordTradeCmd
	tradeBrief         map[uint32]*Cmd.BriefPendingListRecordTradeCmd
	sellInfo           map[uint32]*Cmd.ItemSellInfoRecordTradeCmd
	tradeHistory       *Cmd.MyTradeLogRecordTradeCmd
	pendingSells       *Cmd.MyPendingListRecordTradeCmd
	buyItem            map[uint32]*Cmd.BuyItemRecordTradeCmd
	sellItem           map[uint32]*Cmd.SellItemRecordTradeCmd
	reqServerPrice     map[uint32]*Cmd.ReqServerPriceRecordTradeCmd
	MapNpcs            map[uint64]*Cmd.MapNpc
	MapUsers           map[uint64]*Cmd.MapUser
	GotoList           *Cmd.GoToListUserCmd
	Mutex              *sync.RWMutex
	lastHeartBeat      time.Time
	retries            map[string]uint
	cancelAtkCtx       context.CancelFunc
	ExchangeItems      map[uint32]utils.ExchangeItem
	SkillItems         map[uint32]utils.SkillItem
	SkillItemsByName   map[string][]utils.SkillItem
	BuffItems          map[uint32]utils.BuffItem
	BuffItemsByName    map[string]utils.BuffItemByName
	Items              map[uint32]utils.Items
	ItemsByName        map[string]utils.ItemsByName
	notifier           map[gameTypes.NotifierType]chan interface{}
	MonsterItems       map[uint32]utils.MonsterInfo
	MonsterItemsByName map[string]utils.MonsterInfo
	AtkStat            AttackMonsterStat
	BossInfo           *Cmd.BossListUserCmd
	logger             *log.Logger
	logBuffer          []string
	logMutex           sync.RWMutex
	LogNotify          chan string
	chatHistory        []chatMessage
	chatMutex          sync.RWMutex
	reconnecting       bool
}

func (g *GameConnection) GetServerTimeDelayMilli() uint64 {
	return fixedSkillCDSubtract
}

func (g *GameConnection) GetServerTimeDelaySec() uint32 {
	return uint32(fixedSkillCDSubtract / 1000)
}

func (g *GameConnection) SetQueryTimeout(timeout time.Duration) {
	queryTimeout = timeout
}

func (g *GameConnection) GetItemCat(itemId uint32) uint32 {
	if g.ExchangeItems != nil {
		if _, ok := g.ExchangeItems[itemId]; ok {
			catInt, _ := strconv.ParseUint(g.ExchangeItems[itemId].Category, 10, 32)
			return uint32(catInt)
		}
		log.Warnf("item id %d not found", itemId)
	}
	return 0
}

func (g *GameConnection) SetLogLevel(level log.Level) {
	g.logger.SetLevel(level)
}

func (g *GameConnection) GetLogger() *log.Logger {
	return g.logger
}

func (g *GameConnection) addChatMessage(msg chatMessage) {
	g.chatMutex.Lock()
	defer g.chatMutex.Unlock()

	g.chatHistory = append(g.chatHistory, msg)
	if len(g.chatHistory) > g.Configs.GetChatMaxSize() {
		g.chatHistory = g.chatHistory[len(g.chatHistory)-g.Configs.GetChatMaxSize():]
	}
}

func (g *GameConnection) GetChatHistory() []chatMessage {
	g.chatMutex.RLock()
	defer g.chatMutex.RUnlock()

	history := make([]chatMessage, len(g.chatHistory))
	copy(history, g.chatHistory)
	return history
}

func (g *GameConnection) ClearChatHistory() {
	g.chatMutex.Lock()
	defer g.chatMutex.Unlock()

	g.chatHistory = []chatMessage{}

}

func (g *GameConnection) SentChatMessage(channelId int32, content string, destId uint64) error {
	id := Cmd.EGameChatChannel(channelId)
	msg := &Cmd.ChatCmd{
		Channel: &id,
		Str:     &content,
	}
	if destId != 0 {
		msg.DesID = &destId
	}
	_ = g.sendProtoCmd(msg,
		Cmd.Command_value["CHAT_PROTOCMD"],
		Cmd.ChatParam_value["CHATPARAM_CHAT"],
	)
	return nil
}

// LogWriter returns an io.Writer that writes to the log buffer
func (g *GameConnection) LogWriter() io.Writer {
	return &logWriter{gc: g}
}

func (g *GameConnection) AddLog(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Split message by newlines and add each line separately
	// lines := strings.Split(strings.TrimRight(message, "\n"), "\n")

	g.logMutex.Lock()
	defer g.logMutex.Unlock()

	logLine := fmt.Sprintf("[%s] %s", timestamp, message)
	g.logBuffer = append(g.logBuffer, logLine)

	if len(g.logBuffer) > 1000 {
		g.logBuffer = g.logBuffer[len(g.logBuffer)-1000:]
	}
	select {
	case g.LogNotify <- message:
	default:
	}
}

func (g *GameConnection) GetLogs() []string {
	g.logMutex.RLock()
	defer g.logMutex.RUnlock()

	logs := make([]string, len(g.logBuffer))
	copy(logs, g.logBuffer)
	return logs
}

func (g *GameConnection) ClearLogs() {
	g.logMutex.Lock()
	defer g.logMutex.Unlock()

	g.logBuffer = []string{}
}

func (g *GameConnection) IsAuthed() bool {
	return g.Authed
}

func (g *GameConnection) SetAuthed(Authed bool) {
	g.Authed = Authed
}

const (
	printHeartBeatLogInterval = 60 * time.Second
	maxRetry                  = 1
)

var (
	queryTimeout                = 3 * time.Second
	ErrQueryTimeout             = errors.New("query timeout")
	ErrUseClosedConnection      = errors.New("use of closed network connection")
	ErrConnectionClosedByRemote = errors.New("an existing connection was forcibly closed by the remote host")
)

func (g *GameConnection) GameServerLogin() {
	if g.Configs.AccId == 0 || g.reconnecting {
		err := g.getAccId()
		if err != nil {
			g.logger.Errorf("get accId failed: %v", err)
			return
		}
	}
	g.logger.Infof("Account Id: %d", g.Configs.AccId)
	err := g.connectGameServer()
	if err != nil {
		g.logger.Errorf("connect game server failed: %v", err)
		return
	}
	g.quitContext, g.quitCancel = context.WithCancel(context.Background())
	go g.sendHandler()
	time.Sleep(time.Second)
	g.sendReqUserLoginParamCmd()
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		// handle/parse TCP response
		go g.handleConnection()

		// handle connection heartbeat
		// connection will be closed without sending heartbeat
		for {
			select {
			case <-ticker.C:
				if g.ShouldHeartBeat {
					timeNow := time.Now().Unix()
					if timeNow > g.lastHeartBeat.Add(printHeartBeatLogInterval).Unix() && uint(timeNow) > g.retries["heartbeat"]+60 {
						log.Printf("last heart beat received at %v", g.lastHeartBeat.UTC())
						g.retries["heartbeat"] = uint(timeNow)
					}
					g.sendHeartBeat()
				}
			case <-g.quitContext.Done():
				g.logger.Infof("%s quit", g.Role.GetRoleName())
				ticker.Stop()
				return
			}
		}
	}()
	loginTicker := time.NewTicker(5 * time.Second)
	defer loginTicker.Stop()
loginLoop:
	for {
		select {
		case <-g.quitContext.Done():
			g.logger.Infof("Login loop quit")
			return
		case <-time.After(15 * time.Second):
			g.logger.Infof("Login timeout")
			g.Reconnect()
			return
		case <-loginTicker.C:
			if !g.IsAuthed() {
				continue
			} else {
				break loginLoop
			}
		}
	}
	if g.Configs.AutoCreateChar {
		if _, ok := g.AvailableRoles[uint32(g.Configs.Char)]; !ok {
			// random character name
			if g.Configs.CharacterName == "" {
				g.Configs.CharacterName = utils.RandomZhCharacterName(6)
			}
			err := g.CreateCharacter(
				g.Configs.CharacterName,
				2,
				42,
				12,
				2,
				0,
				1,
			)
			if err != nil {
				g.logger.Error(err)
			}
			g.logger.Infof("Created character %s", g.Configs.CharacterName)
			time.Sleep(3 * time.Second)
		}
	}
	if g.IsAuthed() && g.conn != nil {
		g.SelectRole()
	}
	enterMapTick := time.NewTicker(2 * time.Second)
	defer enterMapTick.Stop()
	enterMapTimeout := time.After(60 * time.Second)
enterMapLoop:
	for {
		select {
		case <-g.quitContext.Done():
			g.logger.Infof("Enter map loop quit")
			return
		case <-enterMapTimeout:
			g.logger.Infof("Enter map timeout")
			g.Reconnect()
			return
		case <-enterMapTick.C:
			if g.conn != nil && g.Role.GetMapId() != 0 && g.Role.GetInGame() && !g.enteringMap && g.Role.GetLoginResult() == 0 {
				g.enterGameMap()
				break enterMapLoop
			}
			g.logger.Infof("Waiting for enter map")
			time.Sleep(time.Second * 2)
		}
	}

	g.WaitForInGame()
}

func (g *GameConnection) WaitForInGame() {
	ticker := time.NewTicker(2 * time.Second)
	timeout := time.After(60 * time.Second)
	for {
		select {
		case <-g.quitContext.Done():
			g.logger.Infof("Wait for in game quit")
			return
		case <-timeout:
			log.Errorf("Wait for in game timeout")
			ticker.Stop()
			g.Reconnect()
			return
		case <-ticker.C:
			if g.Role.GetInGame() {
				// If not moved strange things will happen
				g.MoveChart(g.Role.GetPos())
				if g.Configs.AuthPass != "" {
					log.Infof("sending auth pass")
					g.AuthConfirm(g.Configs.AuthPass)
				}
				ticker.Stop()
				return
			} else {
				log.Warn("Waiting for in game")
				time.Sleep(time.Second * 2)
			}
		}
	}
}

func (g *GameConnection) handleConnection() {
	// defer g.conn.Close()
	// listen for reply
	scanner := bufio.NewReader(g.conn)
	buf := make([]byte, 512000)
	g.logger.Infof("connection handler started")
	g.reconnecting = false
	for {
		select {
		case <-g.quitContext.Done():
			g.logger.Infof("connection handler quit")
			return
		default:
			if g.conn == nil {
				g.logger.Infof("connection handler stopped due to nil connection")
				return
			}
			data, msgFlag, err := g.parseRawTCP(scanner, buf)
			if err != nil {
				switch err {
				default:
					if g.shouldQuit {
						g.Close()
						return
					} else {
						time.Sleep(10 * time.Second)
						g.Reconnect()
						time.Sleep(15 * time.Second)
						return
					}
				}
			}
			contentSize := len(data)
			out := utils.ParseBody(data, utils.CipherKey)

			if contentSize > 512 {
				log.Debugf("received tcpFlag %d, %d bytes from server", int(msgFlag), contentSize+1)
			} else {
				log.Debugf("received tcpFlag %d, %d bytes from server: %x", int(msgFlag), contentSize+1, out)
			}

			if g.DebugMsg {
				utils.TranslateMsg(out)
			}

			g.HandleMsg(out)
		}
	}
}

func (g *GameConnection) parseRawTCP(scanner *bufio.Reader, buf []byte) ([]byte, byte, error) {
	flag, err := scanner.ReadByte()
	if err != nil {
		log.Errorf("tcp read tcpFlag byte err: %s", err)
		return []byte{}, 0, err
	}
	validFlag := utils.IsValidFlag(flag)
	invalidSize := 0
	for {
		if validFlag {
			if invalidSize > 0 {
				log.Debugf("%d bytes drop due to invalid flag", invalidSize)
			}
			break
		}
		invalidSize += 1
		flag, err = scanner.ReadByte()
		validFlag = utils.IsValidFlag(flag)
	}

	bodySize, err := scanner.ReadByte()
	if err != nil {
		log.Errorf("tcp read size byte err: %s", err)
		return []byte{}, 0, err
	}

	body, err := scanner.ReadByte()
	if err != nil {
		log.Errorf("tcp read body byte err: %s", err)
		return []byte{}, 0, err
	}

	dFlag := []byte{flag}
	size := []byte{bodySize, body}
	contentSize := utils.GetContentSize(dFlag, size)
	_, err = io.ReadFull(scanner, buf[:contentSize])
	content := [][]byte{
		dFlag,
		size,
		buf[:contentSize],
	}

	data := make([]byte, contentSize+3)
	var i int
	for _, val := range content {
		i += copy(data[i:], val)
	}
	return data, flag, nil
}

func (g *GameConnection) Reconnect() {
	if g.reconnecting {
		return
	}
	g.reconnecting = true
	defer func() {
		g.reconnecting = false
		g.shouldQuit = false
	}()
	g.logger.Infof("%s Reconnecting", g.Role.GetRoleName())
	if g.conn != nil {
		if g.cancelAtkCtx != nil {
			g.cancelAtkCtx()
			g.AtkStat = AttackMonsterStat{}
		}
		g.Close()
		g.enteringMap = false
		g.Authed = false
		roleOptions := RoleTeamOption(g.Configs.TeamConfig)
		g.Role = NewRole(roleOptions)
		if g.Mutex.TryLock() {
			g.Mutex.Unlock()
		} else {
			g.Mutex.Unlock()
		}
	}
	g.GameServerLogin()
}

func (g *GameConnection) Close() {
	g.shouldQuit = true
	g.SetAuthed(false)
	if g.quitCancel != nil {
		g.quitCancel()
	}
	if g.conn != nil {
		_ = g.conn.Close()
	}
	g.conn = nil
}

func (g *GameConnection) connectGameServer() error {
	log.Infof("Trying to connect to game server: %s for token %s", g.Configs.IpPort, g.Configs.AccessToken)
	conn, err := net.Dial("tcp", g.Configs.IpPort)
	if err != nil {
		return err
	}
	g.conn = conn
	return nil
}

func (g *GameConnection) getResourceVersion() string {
	host := g.Configs.IpPort
	// remove port from host if exists
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	req, err := http.Get(fmt.Sprintf("http://%s/update/Android/version.php", host))
	if err != nil {
		g.logger.Errorf("failed to get resource version: %s", err)
		return "0"
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		g.logger.Errorf("failed to read resource version response body: %s", err)
		return "0"
	}
	var result any
	err = json.Unmarshal(body, &result)
	if err != nil {
		g.logger.Errorf("failed to unmarshal resource version response body: %s", err)
		return "0"
	}
	versionMap, ok := result.(map[string]interface{})
	if !ok {
		g.logger.Errorf("failed to parse resource version response body: %s", string(body))
		return "0"
	}
	okMsg, ok := versionMap["message"].(string)
	if !ok {
		g.logger.Errorf("failed to parse resource version message: %s", string(body))
		return "0"
	}
	if okMsg != "ok" {
		g.logger.Errorf("resource version response not ok: %s", string(body))
		return "0"
	}
	dataMap, ok := versionMap["data"].(map[string]interface{})
	if !ok {
		g.logger.Errorf("failed to parse resource version data: %s", string(body))
		return "0"
	}
	clientVer, ok := dataMap["client"].(string)
	return clientVer
}

func (g *GameConnection) httpAuth(authHost string) (*authJson, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, authHost, nil)
	if err != nil {
		log.Errorf("failed to create http newRequest: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range g.Configs.AuthParams {
		gPtr := *g.Configs
		val := reflect.ValueOf(gPtr).FieldByName(v).Interface()
		if u, ok := val.(uint32); ok {
			q.Add(k, strconv.Itoa(int(u)))
		} else {
			q.Add(k, val.(string))
		}
	}
	resVer := g.getResourceVersion()
	g.logger.Infof("Resource version: %s", resVer)
	q.Add("res", resVer)
	req.URL.RawQuery = q.Encode()
	g.logger.Debugf("Sending request to: %s", req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		g.logger.Errorf("failed to request %s: %s", req.URL.String(), err)
	}
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		g.logger.Errorf("%s", err)
	}
	result := &authJson{}
	err = json.Unmarshal(body, result)
	if err != nil {
		if res.StatusCode == 200 {
			result.Data = strings.TrimSpace(string(body))
			err = nil
		}
	}
	return result, err
}

func (g *GameConnection) SendServerTimeUserCmd(par Cmd.LoginCmdParam) {
	curTime := uint64(time.Now().UnixMilli())
	serverTimeUserCmd := g.getServerTimeUserCmd(curTime, 0, par)
	data, _ := proto.Marshal(serverTimeUserCmd)
	log.Debug(data)
	out := utils.ConstructBody(1, 11, utils.TcpFlag[1], data, g.getNonce(false), utils.CipherKey)
	g.logger.Infof("sending %v bytes serverTimeUserCmd", len(out))
	g.sendCmd(utils.TcpFlag[1], out, 0)
}

func (g *GameConnection) sendProtoCmd(protoCmd proto.Message, cmdId, cmdParId int32) (err error) {
	data, err := proto.Marshal(protoCmd)
	if err != nil {
		log.Errorf("failed to marshal sell info query: %s", err)
	} else {

		body := utils.ConstructBody(
			cmdId,
			cmdParId,
			utils.TcpFlag[1],
			data,
			g.getNonce(true),
			utils.CipherKey,
		)
		g.sendCmd(utils.TcpFlag[1], body, 0)
	}
	return err
}

func (g *GameConnection) sendProtoCmdIndex(protoCmd proto.Message, cmdId, cmdParId int32, index uint32) (err error) {
	data, err := proto.Marshal(protoCmd)
	if err != nil {
		log.Errorf("failed to marshal sell info query: %s", err)
	} else {

		body := utils.ConstructBody(
			cmdId,
			cmdParId,
			utils.TcpFlag[1],
			data,
			g.getNonceIndex(true, index),
			utils.CipherKey,
		)
		g.sendCmd(utils.TcpFlag[1], body, 0)
	}
	return err
}

func (g *GameConnection) sendCmd(flag, body []byte, delay time.Duration) {
	var newBody []byte
	if len(body) > 0 {
		newBody = make([]byte, 3)
		newBody[0] = flag[0]
		binary.LittleEndian.PutUint16(newBody[1:], uint16(len(body)))
		newBody = append(newBody[:], body[:]...)
		time.Sleep(delay)
	}
	// g.Mutex.Lock()
	g.cmdQueue = append(g.cmdQueue, newBody)
	// g.Mutex.Unlock()
}

func (g *GameConnection) sendHandler() {
	ticker := time.NewTicker(cmdQueueInterval)
	defer ticker.Stop()
	g.logger.Infof("send handler started")
	g.reconnecting = false
	for {
		select {
		case <-g.quitContext.Done():
			g.logger.Infof("send handler quit")
			return
		case <-ticker.C:
			g.Mutex.Lock()
			if len(g.cmdQueue) > 0 {
				body := []byte{}
				for _, b := range g.cmdQueue {
					body = append(body, b...)
				}
				g.currentIndex = 0
				if g.conn != nil {
					g.logger.Debugf("sending %v bytes at %v", len(body), time.Now())
					writeLen, err := g.conn.Write(body)
					g.logger.Debugf("sent %d bytes", writeLen)
					if err != nil {
						g.logger.Errorf("%s failed to send command: %v", g.Role.GetRoleName(), err)
						g.cmdQueue = [][]byte{}
						g.Mutex.Unlock()
						time.Sleep(15 * time.Second)
						g.Reconnect()
						time.Sleep(time.Second * 15)
						return
					}
				}
			}
			g.cmdQueue = [][]byte{}
			g.Mutex.Unlock()
		}
	}
}

func (g *GameConnection) getServerTimeUserCmd(curTime uint64, timeZone int32, params Cmd.LoginCmdParam) *Cmd.ServerTimeUserCmd {
	pb := &Cmd.ServerTimeUserCmd{}
	if params != 0 {
		pb.Param = &params
	}
	if curTime > 0 {
		pb.Time = &curTime
	}

	// if timeZone > 0 {
	// 	pb.TimeZone = &timeZone
	// }
	return pb
}

func (g *GameConnection) getReqLoginUserCmd(zId, serverId, language, langZone, clientVer, timestamp uint32, accId uint64, userSha1Str, version, domain, ip, device, phone, safeDevice, authoriz, deviceId string) *Cmd.ReqLoginUserCmd {
	pb := &Cmd.ReqLoginUserCmd{
		Sha1:    &userSha1Str,
		Accid:   &accId,
		Zoneid:  &zId,
		Version: &version,
		// Domain:        &domain,
		// Ip:            &ip,
		// Device:        &device,
		Language: &language,
		// Deviceid: &deviceId,
		// Clientversion: &clientVer,
		Timestamp: &timestamp,
		// Langzone:      &langZone,
	}

	loginSite := uint32(0)
	pb.Site = &loginSite

	// if serverId != 0 {
	// 	pb.Serverid = &serverId
	// }
	//
	// if authoriz != "" {
	// 	pb.Authorize = &authoriz
	// }
	//
	// if safeDevice != "" {
	// 	pb.SafeDevice = &safeDevice
	// }
	//
	// if phone != "" {
	// 	pb.Phone = &phone
	// }
	return pb
}

func (g *GameConnection) enterGameMap() {
	// mid := g.Role.GetMapId()
	// g.ChangeMap(mid)
	g.enteringMap = true
	if g.Role.GetMapId() == 0 {
		log.Warnf("%s map id is 0 when entering game map", g.Role.GetRoleName())
		return
	}
	g.ChangeMap(g.Role.GetMapId())
	// If not moved strange things will happen
	// g.MoveChart(g.Role.GetPos())
}

func (g *GameConnection) QueryCat(catId uint32) (results *Cmd.BriefPendingListRecordTradeCmd) {
	results = &Cmd.BriefPendingListRecordTradeCmd{}
	cmd := &Cmd.BriefPendingListRecordTradeCmd{
		Category: &catId,
		Charid:   g.Role.RoleId,
	}
	err := g.sendProtoCmd(cmd, TradeProtoCmdId, Cmd.RecordUserTradeParam_value["BRIEF_PENDING_LIST_RECORDTRADE"])
	err = g.waitForQueryResponse(catId)
	for retry := 0; err != nil && retry < 3; retry++ {
		log.Infof("retrying #%d query category %d", retry, catId)
		_ = g.sendProtoCmd(cmd, TradeProtoCmdId, Cmd.RecordUserTradeParam_value["BRIEF_PENDING_LIST_RECORDTRADE"])
		err = g.waitForQueryResponse(catId)
	}
	if err != nil {
		log.Error(err)
	} else {
		results = g.tradeBrief[catId]
	}
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	delete(g.tradeBrief, catId)
	return results
}

func (g *GameConnection) waitForQueryResponse(key uint32) (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("Checking for query response")
		g.Mutex.Lock()
		if g.tradeDetail[key] == nil && g.tradeBrief[key] == nil && g.sellInfo[key] == nil {
			g.Mutex.Unlock()
			time.Sleep(time.Second)
			continue
		} else {
			g.Mutex.Unlock()
			break
		}
	}
	if time.Since(startQueryTime) > queryTimeout {
		err = ErrQueryTimeout
	}
	return err
}

func (g *GameConnection) QueryItemPrice(itemId uint32, pageIndex uint32) (results []*Cmd.TradeItemBaseInfo) {
	results = []*Cmd.TradeItemBaseInfo{}
	tt := Cmd.ETradeType_ETRADETYPE_ALL
	cmd := &Cmd.DetailPendingListRecordTradeCmd{
		SearchCond: &Cmd.SearchCond{
			ItemId:    &itemId,
			TradeType: &tt,
			PageIndex: &pageIndex,
		},
		Charid: g.Role.RoleId,
	}
	_ = g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["DETAIL_PENDING_LIST_RECORDTRADE"],
	)

	err := g.waitForQueryResponse(itemId)
	for retry := 0; err != nil && retry < 3; retry++ {
		log.Infof("retrying #%d query item %d", retry, itemId)
		_ = g.sendProtoCmd(
			cmd,
			Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
			Cmd.RecordUserTradeParam_value["DETAIL_PENDING_LIST_RECORDTRADE"],
		)
		err = g.waitForQueryResponse(itemId)
	}
	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query for item %d return error: %s", itemId, err)
	} else if g.tradeDetail != nil && g.tradeDetail[itemId] != nil {
		results = g.tradeDetail[itemId].GetLists()
	}
	pageCount := g.tradeDetail[itemId].GetTotalPageCount()
	delete(g.tradeDetail, itemId)
	g.Mutex.Unlock()

	if pageCount > 1 && pageIndex == 0 {
		for i := uint32(1); i <= pageCount; i++ {
			r1 := g.QueryItemPrice(itemId, i)
			results = append(results, r1...)
		}
	}

	return results
}

func (g *GameConnection) QueryItemSellInfo(itemId, publicityId uint32) (result *Cmd.ItemSellInfoRecordTradeCmd) {
	id64 := uint64(publicityId)
	cmd := &Cmd.ItemSellInfoRecordTradeCmd{
		Charid:      g.Role.RoleId,
		Itemid:      &itemId,
		PublicityId: &publicityId,
		OrderId:     &id64,
	}
	_ = g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["ITEM_SELL_INFO_RECORDTRADE"],
	)

	err := g.waitForQueryResponse(itemId)
	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query sell info for item %d, publicity %d return error: %s", itemId, publicityId, err)
	} else if g.sellInfo != nil && g.sellInfo[itemId] != nil {
		result = g.sellInfo[itemId]
	}
	delete(g.sellInfo, itemId)
	g.Mutex.Unlock()
	return result
}

func (g *GameConnection) sendHeartBeat() {
	curTime := uint64(utils.GetTimeNow(true))
	cmd := &Cmd.HeartBeatUserCmd{
		Time: &curTime,
	}
	_ = g.sendProtoCmd(
		cmd,
		LogInUserProtoCmdId,
		Cmd.LoginCmdParam_value["HEART_BEAT_USER_CMD"],
	)
	g.currentIndex = 1
}

func (g *GameConnection) getNonce(includeTime bool) []byte {
	currentTime := int64(0)
	if includeTime {
		currentTime = utils.GetTimeNow(false)
	}
	g.currentIndex += 1
	sign := fmt.Sprintf("%d_%d_!^ro&", currentTime, g.currentIndex)
	signStr := fmt.Sprintf("%x", sha1.Sum([]byte(sign)))
	nonce := &Cmd.Nonce{
		Index: &g.currentIndex,
		Sign:  &signStr,
	}
	if includeTime {
		newT := uint32(currentTime)
		nonce.Timestamp = &newT
	}
	pOut, _ := proto.Marshal(nonce)
	return pOut
}

// getNonceIndex generates a nonce with a specific index
func (g *GameConnection) getNonceIndex(includeTime bool, index uint32) []byte {
	currentTime := int64(0)
	if includeTime {
		currentTime = utils.GetTimeNow(false)
	}
	sign := fmt.Sprintf("%d_%d_!^ro&", currentTime, index)
	signStr := fmt.Sprintf("%x", sha1.Sum([]byte(sign)))
	nonce := &Cmd.Nonce{
		Index: &index,
		Sign:  &signStr,
	}
	if includeTime {
		newT := uint32(currentTime)
		nonce.Timestamp = &newT
	}
	pOut, _ := proto.Marshal(nonce)
	return pOut
}

func (g *GameConnection) SendReqUserLoginCmd(timeStamp uint32) {
	reqLoginCmd := g.getReqLoginUserCmd(
		g.Configs.ZoneId,
		g.Configs.ServerId,
		g.Configs.Lang,
		g.Configs.LangZone,
		g.Configs.AppPreVer,
		timeStamp,
		g.Configs.AccId,
		g.Configs.Sha1Str,
		g.Configs.Version,
		g.Configs.Domain,
		g.Configs.Ip,
		g.Configs.Device,
		g.Configs.Phone,
		g.Configs.SafeDevice,
		g.Configs.Authoriz,
		g.Configs.DeviceId,
	)
	time.Sleep(1 * time.Second)
	_ = g.sendProtoCmd(reqLoginCmd, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["REQ_LOGIN_USER_CMD"])
}

func (g *GameConnection) sendReqUserLoginParamCmd() {
	reqLoginParamCmd := Cmd.ReqLoginParamUserCmd{
		Accid: &g.Configs.AccId,
		Pwd:   &g.Configs.Password,
	}
	_ = g.sendProtoCmd(&reqLoginParamCmd, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["REQ_LOGIN_PARAM_USER_CMD"])
}

func (g *GameConnection) SelectRole() {
	// retry 3 times
	roleMaxRetry := 3
	for retry := 0; retry < roleMaxRetry; retry++ {
		if len(g.AvailableRoles) == 0 {
			log.Error("no available roles waiting")
			if retry >= roleMaxRetry-1 {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
	role, ok := g.AvailableRoles[uint32(g.Configs.Char)]
	if !ok {
		log.Errorf("no available role with at %d", g.Configs.Char)
		g.Close()
	}
	g.Role = role
	g.logger.Infof("selecting role with id: %d name: %s", g.Role.GetRoleId(), g.Role.GetRoleName())
	g.doSelectRole()
}

func (g *GameConnection) doSelectRole() {
	g.currentIndex = 0
	tp := Cmd.EOptionType_EOPTIONTYPE_USE_SLIM
	cmd := &Cmd.NewSetOptionUserCmd{
		Type: &tp,
	}
	_ = g.sendProtoCmd(cmd, sceneUser2CmdId, Cmd.User2Param_value["USER2PARAM_NEW_SET_OPTION"])

	// k := "0000"
	// cp := uint32(13510102)
	rid := g.Role.GetRoleId()
	cmd1 := &Cmd.SelectRoleUserCmd{
		Id:        &rid,
		Deviceid:  &g.Configs.DeviceId,
		ExtraData: &Cmd.ExtraData{
			// System: &g.Configs.Device,
			// Model:  &g.Configs.Model,
		},
		// Pushkey:  &k,
		// Clickpos: &cp,
	}
	log.Infof("sending select role command: %v", cmd1)
	_ = g.sendProtoCmd(cmd1, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["SELECT_ROLE_USER_CMD"])
	g.Role.SetRoleSelected(true)
	_ = g.GetAllPackItems()
}

func (g *GameConnection) IsEnteringGameMap() bool {
	return g.enteringMap
}

func (g *GameConnection) InGameMap() bool {
	return g.inMap
}

func (g *GameConnection) getAccId() (err error) {
	if g.Configs.Username != "" && g.Configs.Password != "" {
		res, err := g.httpAuth(g.Configs.AuthServer)
		if err != nil {
			err = fmt.Errorf("failed to login to auth server for user %s: %w", g.Configs.Username, err)
			return err
		}
		g.Configs.AccId, err = strconv.ParseUint(res.Data, 10, 64)
		if err != nil {
			err = fmt.Errorf("failed to parse account id %s: %w", res.Data, err)
			return err
		}
	} else {
		err = fmt.Errorf("no account id nor username and password provided for login: %w", err)
	}
	return err
}

func (g *GameConnection) LoadMonster(monsterJsonPath string) *GameConnection {
	monsters := utils.MonsterParser(monsterJsonPath)
	g.Mutex.Lock()
	for _, monster := range monsters {
		g.MonsterItems[uint32(monster.Id)] = monster
		g.MonsterItemsByName[monster.NameZh] = monster
	}
	g.Mutex.Unlock()
	return g
}

func (g *GameConnection) SetEnteringMap() {
	g.enteringMap = true
}

func NewConnection(Config *config.ServerConfigs, skillItems map[uint32]utils.SkillItem, items *utils.ItemsLoader) *GameConnection {
	si := map[uint32]utils.SkillItem{}
	siName := map[string][]utils.SkillItem{}
	if skillItems != nil {
		si = skillItems
		if Config.HuntConfig.PrepEliteCD > 0 {
			if skill, ok := si[uint32(50057001)]; ok {
				skill.CD = strconv.Itoa(Config.HuntConfig.PrepEliteCD)
				si[uint32(50057001)] = skill
			}
		}

		for _, skill := range skillItems {
			if _, ok := siName[skill.NameZh]; !ok {
				siName[skill.NameZh] = []utils.SkillItem{}
			} else {
				siName[skill.NameZh] = append(siName[skill.NameZh], skill)
			}
		}
	}
	bi := map[uint32]utils.BuffItem{}
	if items != nil && items.BuffItems != nil {
		bi = items.BuffItems
	}
	bin := map[string]utils.BuffItemByName{}
	if items != nil && items.BuffItemsByName != nil {
		bin = items.BuffItemsByName
	}
	ei := map[uint32]utils.ExchangeItem{}
	if items != nil && items.ExchangeItems != nil {
		ei = items.ExchangeItems
	}

	allItems := map[uint32]utils.Items{}
	if items != nil && items.Items != nil {
		allItems = items.Items
	}

	allItemsByName := map[string]utils.ItemsByName{}
	if items != nil && items.ItemsByName != nil {
		allItemsByName = items.ItemsByName
	}
	roleOption := RoleTeamOption(Config.TeamConfig)
	gc := &GameConnection{
		Configs:            Config,
		Role:               NewRole(roleOption),
		AvailableRoles:     map[uint32]*RoleInfo{},
		currentIndex:       1,
		tradeBrief:         map[uint32]*Cmd.BriefPendingListRecordTradeCmd{},
		tradeDetail:        map[uint32]*Cmd.DetailPendingListRecordTradeCmd{},
		sellInfo:           map[uint32]*Cmd.ItemSellInfoRecordTradeCmd{},
		tradeHistory:       &Cmd.MyTradeLogRecordTradeCmd{},
		buyItem:            map[uint32]*Cmd.BuyItemRecordTradeCmd{},
		sellItem:           map[uint32]*Cmd.SellItemRecordTradeCmd{},
		pendingSells:       &Cmd.MyPendingListRecordTradeCmd{},
		reqServerPrice:     map[uint32]*Cmd.ReqServerPriceRecordTradeCmd{},
		Mutex:              &sync.RWMutex{},
		retries:            map[string]uint{},
		MapNpcs:            map[uint64]*Cmd.MapNpc{},
		MapUsers:           map[uint64]*Cmd.MapUser{},
		SkillItems:         si,
		SkillItemsByName:   siName,
		BuffItems:          bi,
		BuffItemsByName:    bin,
		ExchangeItems:      ei,
		Items:              allItems,
		ItemsByName:        allItemsByName,
		MonsterItems:       map[uint32]utils.MonsterInfo{},
		MonsterItemsByName: map[string]utils.MonsterInfo{},
		notifier:           map[gameTypes.NotifierType]chan interface{}{},
		logMutex:           sync.RWMutex{},
		LogNotify:          make(chan string, 100),
		chatHistory:        []chatMessage{},
		logger:             log.New(),
	}
	gc.logger.SetOutput(io.MultiWriter(os.Stdout, gc.LogWriter()))
	if Config.Loglines > 0 {
		gc.logBuffer = make([]string, 0, Config.Loglines)
	} else {
		gc.logBuffer = make([]string, 0, 1000)
	}
	if gc.MonsterItemsByName == nil {
		gc.MonsterItemsByName = map[string]utils.MonsterInfo{}
	}
	if gc.MonsterItems == nil {
		gc.MonsterItems = map[uint32]utils.MonsterInfo{}
	}
	if gc.Configs.HuntConfig.HuntBossConfig == nil {
		gc.Configs.HuntConfig.HuntBossConfig = &config.HuntBossConfig{}
	}
	if gc.Configs.HuntConfig.HuntBossConfig == nil {
		gc.Configs.HuntConfig.HuntBossConfig = &config.HuntBossConfig{}
	}
	return gc
}

func (g *GameConnection) IsTCPConnected() bool {
	return g.conn != nil
}

func (g *GameConnection) EquipCard(cardGuid, EquipGuid string, slot uint32, operation Cmd.ECardOper) error {
	cmd := &Cmd.EquipCard{
		Cardguid:  &cardGuid,
		Equipguid: &EquipGuid,
		Pos:       &slot,
		Oper:      &operation,
	}
	return g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"],
		Cmd.ItemParam_value["ITEMPARAM_EQUIPCARD"],
	)
}

func (g *GameConnection) EquipCardOn(cardGuid, EquipGuid string, slot uint32) error {
	op := Cmd.ECardOper_ECARDOPER_EQUIPON
	return g.EquipCard(cardGuid, EquipGuid, slot, op)
}

func (g *GameConnection) EquipCardOff(cardGuid, EquipGuid string, slot uint32) error {
	op := Cmd.ECardOper_ECARDOPER_EQUIPOFF
	return g.EquipCard(cardGuid, EquipGuid, slot, op)
}

func (g *GameConnection) CheckDraculaBuff() {
	if g.GetBuffByID(51551).BuffName != "" {
		g.logger.Trace("德古拉男爵卡片已激活")
		return
	}

	draculaCard := g.FindPackItemByName("德古拉男爵卡片", Cmd.EPackType_EPACKTYPE_MAIN)
	if draculaCard == nil {
		g.logger.Trace("没有找到德古拉男爵卡片")
		return
	}

	g.logger.Info("使用德古拉男爵卡片")

	// get current equip card
	var weapon *Cmd.ItemData
	equipItems := g.Role.GetPackItemsByType(Cmd.EPackType_EPACKTYPE_EQUIP)
	for _, equipItem := range equipItems {
		if equipItem.GetBase().GetEquipType() == Cmd.EEquipType_EEQUIPTYPE_WEAPON {
			weapon = equipItem
			break
		}
	}
	if weapon == nil {
		g.logger.Warn("没有找到装备的武器，无法装备德古拉男爵卡片")
		return
	}
	firstCard := &Cmd.CardData{}
	equipedCards := weapon.GetCard()
	for i, card := range equipedCards {
		if i == 0 {
			firstCard = card
		}
		if card.GetId() == draculaCard.GetBase().GetId() {
			g.logger.Info("德古拉男爵卡片已装备")
			return
		}
	}
	// equip card
	err := g.EquipCardOn(draculaCard.GetBase().GetGuid(), weapon.GetBase().GetGuid(), 1)
	if err != nil {
		g.logger.Errorf("装备德古拉男爵卡片失败: %s", err)
		return
	}
	time.Sleep(time.Millisecond * 500)

	// add go routine to reequip previous card after buff is applied
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		timeOut := time.After(time.Minute * 2)
		for {
			select {
			case <-timeOut:
				g.logger.Info("德古拉男爵卡片监控协程超时，停止监控")
				ticker.Stop()
				return
			case <-ticker.C:
				if g.GetBuffByID(51551).BuffName != "" {
					g.logger.Info("德古拉男爵卡片buff已应用，重新装备之前的卡片")
					err := g.EquipCardOn(firstCard.GetGuid(), weapon.GetBase().GetGuid(), 1)
					if err != nil {
						g.logger.Errorf("重新装备之前的卡片失败: %s", err)
					}
					ticker.Stop()
					return
				}
			case <-g.quitContext.Done():
				g.logger.Info("停止德古拉男爵卡片监控协程")
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *GameConnection) UsesPlayDead() {
	g.logger.Infof("使用装死!")
	num := int32(1)
	dir := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DIR))
	pData := &Cmd.PhaseData{
		Number: &num,
		Pos:    g.Role.Pos,
		Dir:    &dir,
	}
	g.SkillCmd(10020001, pData, true)
}

func (g *GameConnection) EnableGodMode() {
	if g.GetBuffByName("装死(无敌)").BuffName == "" {
		time.Sleep(time.Millisecond * 1000)
		g.CheckuseFlyWing()
		time.Sleep(time.Millisecond * 3200)
		g.UsesPlayDead()
		time.Sleep(time.Millisecond * 1000)
		g.UsesPlayDead()
		time.Sleep(time.Millisecond * 1000)
	}
}

func (g *GameConnection) CheckuseFlyWing() {
	g.buyFlyWing()
	g.UseFlyWing()
	item := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil && item.GetBase().GetCount() > 0 {
		g.logger.Infof("使用苍蝇翅膀 还有%d个", item.GetBase().GetCount())
	} else {
		g.logger.Warn("没有找到苍蝇翅膀")
		_ = g.GetMainPackItems()
	}
}

func (g *GameConnection) buyFlyWing() {
	if item := g.FindPackItemByName("苍蝇翅膀", Cmd.EPackType_EPACKTYPE_MAIN); item == nil || item.GetBase().GetCount() > 1000 {
		return
	}
	shopConfig, err := g.QueryShopConfig(gameTypes.ShopType_Item, 1)
	if err != nil {
		g.logger.Errorf("查询商店配置失败 %s", err)
		return
	}
	for _, item := range shopConfig.GetGoods() {
		if item.GetItemid() == 5024 {
			g.logger.Infof("购买999苍蝇翅膀")
			g.BuyShopItem(item, 999)
		}
	}
}

func (g *GameConnection) InMap(MapID uint32, CarryTeam bool) {
	if g.Role.GetMapId() != MapID {
		if MapID == gameTypes.MapId_LhzDun03.Uint32() {
			g.GoToMap(gameTypes.MapId_LhzDun01.Uint32())
			time.Sleep(time.Millisecond * 500)
			g.MoveChartWait(g.ParsePos(21948, -583, 43399))
			time.Sleep(time.Millisecond * 500)
			g.ExitMapPos(gameTypes.MapId_LhzDun01.Uint32(), 2, g.Role.GetPos())
			time.Sleep(time.Millisecond * 500)
			// g.EnableGodMode()
			g.MoveChartWait(g.ParsePos(-14088, 357, -56505))
			time.Sleep(time.Millisecond * 500)
			g.ExitMapPos(gameTypes.MapId_LhzDun02.Uint32(), 3, g.Role.GetPos())
			time.Sleep(time.Millisecond * 500)
		} else if MapID == gameTypes.MapId_LhzDun02.Uint32() {
			g.GoToMap(gameTypes.MapId_LhzDun01.Uint32())
			time.Sleep(time.Millisecond * 500)
			g.MoveChartWait(g.ParsePos(21948, -583, 43399))
			time.Sleep(time.Millisecond * 500)
			g.ExitMapPos(gameTypes.MapId_LhzDun01.Uint32(), 2, g.Role.GetPos())
			time.Sleep(time.Millisecond * 500)
		} else if MapID == gameTypes.MapId_LhzDun02West.Uint32() {
			g.GoToMap(gameTypes.MapId_LhzDun01.Uint32())
			time.Sleep(time.Millisecond * 500)
			g.MoveChartWait(g.ParsePos(21948, -583, 43399))
			time.Sleep(time.Millisecond * 500)
			g.ExitMapPos(gameTypes.MapId_LhzDun01.Uint32(), 2, g.Role.GetPos())
			time.Sleep(time.Millisecond * 500)
			// g.EnableGodMode()
			g.MoveChartWait(g.ParsePos(46754, 357, 857))
			time.Sleep(time.Millisecond * 500)
			g.ExitMapPos(gameTypes.MapId_LhzDun02.Uint32(), 2, g.Role.GetPos())
		} else {
			if CarryTeam {
				if g.GetCurrentTeamName() == "" || g.AllMemberOffline() {
					g.logger.Warn("没有队伍/没有队员在线，无法传送队伍, 仅自己传送")
					g.GoToMap(MapID)
					return
				}
				g.TeamGoToMap(MapID)
			} else {
				g.GoToMap(MapID)
			}
		}
		g.DelBuffByName("装死(无敌)")
	}

	// g.EnableGodMode()

}

func (g *GameConnection) UseElementStone(stoneType gameTypes.ElementArrowType) {
	item := g.FindPackItemByName(string(stoneType), Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		g.logger.Tracef("%s没有找到", string(stoneType))
	} else {
		if g.GetBuffByName(string(stoneType)).BuffName == "" {
			g.logger.Infof("使用%s", string(stoneType))
			g.UseItem(item.GetBase().GetGuid(), 1)
			time.Sleep(time.Millisecond * 1000)
		}

	}
}

func (g *GameConnection) UseElementArrow(arrowType gameTypes.ElementArrowType) {
	item := g.FindPackItemByName(string(arrowType), Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		g.logger.Tracef("%s没有找到", arrowType)
	} else {
		if item.GetBase().GetIsactive() {
			g.logger.Tracef("%s已装备", arrowType)
			return
		}
		g.logger.Infof("使用%s", arrowType)
		g.UseItem(item.GetBase().GetGuid(), 0)
		time.Sleep(time.Millisecond * 1000)
	}
}
