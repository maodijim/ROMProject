package gameConnection

import (
	"bufio"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"

	log "github.com/sirupsen/logrus"
)

type authJson struct {
	Status int    `json:"code"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

var (
	cmdQueueInterval    = 300 * time.Millisecond
	TradeProtoCmdId     = Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"]
	LogInUserProtoCmdId = Cmd.Command_value["LOGIN_USER_PROTOCMD"]
)

type GameConnection struct {
	Authed             bool
	cmdQueue           [][]byte
	lastCmdSend        time.Time
	currentIndex       uint32
	ShouldHeartBeat    bool
	ShouldChangeScene  bool
	Configs            *config.ServerConfigs
	conn               net.Conn
	Role               *utils.RoleInfo
	AvailableRoles     map[uint32]*utils.RoleInfo
	DebugMsg           bool
	quit               chan bool
	shouldQuit         bool
	enteringMap        bool
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
	ExchangeItems      map[uint32]utils.ExchangeItem
	SkillItems         map[uint32]utils.SkillItem
	BuffItems          map[uint32]utils.BuffItem
	BuffItemsByName    map[string]utils.BuffItemByName
	Items              map[uint32]utils.Items
	ItemsByName        map[string]utils.ItemsByName
	notifier           map[string]chan interface{}
	MonsterItems       map[uint32]utils.MonsterInfo
	MonsterItemsByName map[string]utils.MonsterInfo
	lastAttack         time.Time
}

func (g *GameConnection) IsAuthed() bool {
	return g.Authed
}

func (g *GameConnection) SetAuthed(Authed bool) {
	g.Authed = Authed
}

const (
	queryTimeout              = 8 * time.Second
	printHeartBeatLogInterval = 60 * time.Second
	maxRetry                  = 0
)

var (
	ErrQueryTimeout             = errors.New("query timeout")
	ErrUseClosedConnection      = errors.New("use of closed network connection")
	ErrConnectionClosedByRemote = errors.New("an existing connection was forcibly closed by the remote host")
)

func (g *GameConnection) GameServerLogin() {
	if g.Configs.AccId == 0 {
		err := g.getAccId()
		if err != nil {
			log.Fatalf("get accId failed: %v", err)
		}
	}
	g.connectGameServer()
	g.sendReqUserLoginParamCmd()
	ticker := time.NewTicker(5 * time.Second)
	g.quit = make(chan bool)

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
			case <-g.quit:
				log.Infof("%s quit", g.Role.GetRoleName())
				ticker.Stop()
				return
			}
		}
	}()

loginLoop:
	for {
		select {
		case <-time.After(15 * time.Second):
			log.Infof("Login timeout")
			break loginLoop
		case <-time.Tick(5 * time.Second):
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
				g.Configs.CharacterName = utils.RandomString(7)
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
				log.Error(err)
			}
			log.Infof("Created character %s", g.Configs.CharacterName)
		}
	}

	if g.IsAuthed() && g.conn != nil {
		g.SelectRole()
	}
	if g.conn != nil && g.Role.GetMapId() != 0 && g.Role.GetInGame() && !g.enteringMap && g.Role.GetLoginResult() == 0 {
		g.enterGameMap()
	}
}

func (g *GameConnection) handleConnection() {
	// defer g.conn.Close()
	// listen for reply
	scanner := bufio.NewReader(g.conn)
	buf := make([]byte, 512000)

	for {
		select {
		case <-g.quit:
			return
		default:
			if g.conn == nil {
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
						time.Sleep(8 * time.Second)
						g.Reconnect()
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
	log.Infof("%s Reconnecting", g.Role.GetRoleName())
	if g.conn != nil {
		g.Close()
		g.enteringMap = false
		g.Role = utils.NewRole()
	}
	g.shouldQuit = false
	g.GameServerLogin()
}

func (g *GameConnection) Close() {
	g.quit <- true
	g.shouldQuit = true
	g.SetAuthed(false)
	_ = g.conn.Close()
}

func (g *GameConnection) connectGameServer() {
	log.Infof("Trying to connect to game server: %s for token %s", g.Configs.IpPort, g.Configs.AccessToken)
	conn, err := net.Dial("tcp", g.Configs.IpPort)
	if err != nil {
		log.Errorf("Failed to connect to server %s: %s", g.Configs.IpPort, err)
	}
	g.conn = conn
}

func (g *GameConnection) httpAuth(authHost string) (*authJson, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, authHost, nil)
	if err != nil {
		log.Fatalf("failed to create http newRequest: %s", err)
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
	req.URL.RawQuery = q.Encode()
	log.Printf("Sending request to: %s", req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to request %s: %s", req.URL.String(), err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%s", err)
	}
	result := &authJson{}
	err = json.Unmarshal(body, result)
	return result, err
}

func (g *GameConnection) SendServerTimeUserCmd(par Cmd.LoginCmdParam) {
	serverTimeUserCmd := g.getServerTimeUserCmd(0, 0, par)
	data, _ := proto.Marshal(serverTimeUserCmd)
	log.Debug(data)
	out := utils.ConstructBody(1, 11, utils.TcpFlag[1], data, g.getNonce(false), utils.CipherKey)
	log.Infof("sending %v bytes serverTimeUserCmd", len(out))
	_ = g.sendCmd(utils.TcpFlag[1], out, 0)
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
		err = g.sendCmd(utils.TcpFlag[1], body, 0)
	}
	return err
}

func (g *GameConnection) sendCmd(flag, body []byte, delay time.Duration) (err error) {
	var newBody []byte
	if len(body) > 0 {
		newBody = make([]byte, 3)
		newBody[0] = flag[0]
		binary.LittleEndian.PutUint16(newBody[1:], uint16(len(body)))
		newBody = append(newBody[:], body[:]...)
		time.Sleep(delay)
	}
	if time.Since(g.lastCmdSend) <= cmdQueueInterval {
		g.cmdQueue = append(g.cmdQueue, newBody)
		go func() {
			select {
			case <-time.After(cmdQueueInterval):
				if len(g.cmdQueue) > 0 {
					err = g.sendCmd(flag, body, 0)
				}
			}
		}()
		return
	} else if len(g.cmdQueue) > 0 {
		var oldBody []byte
		for _, b := range g.cmdQueue {
			oldBody = append(oldBody, b...)
		}
		newBody = append(oldBody, newBody...)
		g.cmdQueue = [][]byte{}
		g.currentIndex = 0
	}
	g.lastCmdSend = time.Now()
	if g.conn != nil {
		log.Debugf("sending %v bytes at %v", len(newBody), time.Now())
		writeLen, err := g.conn.Write(newBody)
		log.Debugf("sent %d bytes", writeLen)
		if err != nil {
			log.Errorf("%s failed to send command: %v", g.Role.GetRoleName(), err)
			// g.Close()
		}
	}
	return err
}

func (g *GameConnection) getServerTimeUserCmd(curTime uint64, timeZone int32, params Cmd.LoginCmdParam) *Cmd.ServerTimeUserCmd {
	pb := &Cmd.ServerTimeUserCmd{}
	if params != 0 {
		pb.Param = &params
	}
	if curTime > 0 {
		pb.Time = &curTime
	}
	if timeZone > 0 {
		pb.TimeZone = &timeZone
	}
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
		Deviceid: &deviceId,
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
}

func (g *GameConnection) ExitMapWait(mapId uint32) {
	g.ExitMap(mapId)
	for {
		select {
		case <-time.Tick(2 * time.Second):
			if g.Role.GetMapId() != mapId {
				continue
			}
			g.ChangeMap(mapId)
			return
		}
	}
}

func (g *GameConnection) ExitMap(targetMapId uint32) {
	cmd := &Cmd.GoToExitPosUserCmd{
		Mapid: &targetMapId,
	}
	_ = g.sendProtoCmd(cmd, sceneUserCmdId, Cmd.CmdParam_value["GOTO_EXIT_POS_USER_CMD"])
}

func (g *GameConnection) ChangeMap(mId uint32) {
	cmd := &Cmd.ChangeSceneUserCmd{
		MapID: &mId,
	}
	log.Infof("%s is sending change scene cmd: %v", g.Role.GetRoleName(), cmd)
	g.sendProtoCmd(cmd, 5, 23)
	g.enteringMap = true
	g.inMap = true
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
	reqLoginParamCmd := Cmd.ReqLoginUserCmd{
		Accid: &g.Configs.AccId,
	}
	_ = g.sendProtoCmd(&reqLoginParamCmd, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["REQ_LOGIN_PARAM_USER_CMD"])
}

func (g *GameConnection) SelectRole() {
	if len(g.AvailableRoles) == 0 {
		log.Error("no available roles")
		return
	}
	role, ok := g.AvailableRoles[uint32(g.Configs.Char)]
	if !ok {
		log.Errorf("no available role with at %d", g.Configs.Char)
		g.Close()
	}
	g.Role = role
	log.Infof("selecting role with id: %d name: %s", g.Role.GetRoleId(), g.Role.GetRoleName())
	g.doSelectRole()
}

func (g *GameConnection) doSelectRole() {
	g.currentIndex = 0
	tp := Cmd.EOptionType_EOPTIONTYPE_USE_SLIM
	cmd := &Cmd.NewSetOptionUserCmd{
		Type: &tp,
	}
	_ = g.sendProtoCmd(cmd, sceneUser2CmdId, Cmd.User2Param_value["USER2PARAM_NEW_SET_OPTION"])

	k := "0000"
	cp := uint32(13510102)
	rid := g.Role.GetRoleId()
	cmd1 := &Cmd.SelectRoleUserCmd{
		Id:       &rid,
		Deviceid: &g.Configs.DeviceId,
		ExtraData: &Cmd.ExtraData{
			System: &g.Configs.Device,
			Model:  &g.Configs.Model,
		},
		Pushkey:  &k,
		Clickpos: &cp,
	}
	log.Infof("sending select role command: %v", cmd1)
	_ = g.sendProtoCmd(cmd1, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["SELECT_ROLE_USER_CMD"])
	g.Role.SetRoleSelected(true)
}

func (g *GameConnection) IsEnteringGameMap() bool {
	return g.enteringMap
}

func (g *GameConnection) InGameMap() bool {
	return g.inMap
}

func (g *GameConnection) AddNotifier(notifierType string) {
	g.Mutex.Lock()
	g.notifier[notifierType] = make(chan interface{})
	g.Mutex.Unlock()
}

func (g *GameConnection) removeNotifier(notifierType string) {
	g.Mutex.Lock()
	g.notifier[notifierType] = nil
	g.Mutex.Unlock()
}

func (g *GameConnection) getAccId() (err error) {
	if g.Configs.Username != "" && g.Configs.Password != "" {
		res, err := g.httpAuth(g.Configs.AuthServer)
		if err != nil {
			err = fmt.Errorf("failed to login to auth server for user %s: %w", g.Configs.Username, err)
			return err
		}
		g.Configs.AccId, _ = strconv.ParseUint(res.Data, 10, 64)
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

func NewConnection(config *config.ServerConfigs, skillItems map[uint32]utils.SkillItem, items *utils.ItemsLoader) *GameConnection {
	si := map[uint32]utils.SkillItem{}
	if skillItems != nil {
		si = skillItems
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

	gc := &GameConnection{
		Configs:            config,
		Role:               utils.NewRole(),
		AvailableRoles:     map[uint32]*utils.RoleInfo{},
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
		BuffItems:          bi,
		BuffItemsByName:    bin,
		ExchangeItems:      ei,
		Items:              allItems,
		ItemsByName:        allItemsByName,
		MonsterItems:       map[uint32]utils.MonsterInfo{},
		MonsterItemsByName: map[string]utils.MonsterInfo{},
		notifier:           map[string]chan interface{}{},
	}
	if gc.MonsterItemsByName == nil {
		gc.MonsterItemsByName = map[string]utils.MonsterInfo{}
	}
	if gc.MonsterItems == nil {
		gc.MonsterItems = map[uint32]utils.MonsterInfo{}
	}
	return gc
}

func (g *GameConnection) IsTCPConnected() bool {
	return g.conn != nil
}
