package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/utils"
	"bufio"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type authJson struct {
	Status int                    `json:"status"`
	Msg    string                 `json:"message"`
	Data   map[string]interface{} `json:"data"`
}

var (
	cmdQueueInterval    = 500 * time.Millisecond
	TradeProtoCmdId     = Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"]
	LogInUserProtoCmdId = Cmd.Command_value["LOGIN_USER_PROTOCMD"]
)

type GameConnection struct {
	cmdQueue          [][]byte
	lastCmdSend       time.Time
	currentIndex      uint32
	ShouldHeartBeat   bool
	ShouldChangeScene bool
	Configs           *config.ServerConfigs
	conn              net.Conn
	Role              *utils.RoleInfo
	DebugMsg          bool
	quit              chan bool
	shouldQuit        bool
	enteringMap       bool
	inMap             bool
	tradeDetail       map[uint32]*Cmd.DetailPendingListRecordTradeCmd
	tradeBrief        map[uint32]*Cmd.BriefPendingListRecordTradeCmd
	sellInfo          map[uint32]*Cmd.ItemSellInfoRecordTradeCmd
	tradeHistory      *Cmd.MyTradeLogRecordTradeCmd
	pendingSells      *Cmd.MyPendingListRecordTradeCmd
	buyItem           map[uint32]*Cmd.BuyItemRecordTradeCmd
	sellItem          map[uint32]*Cmd.SellItemRecordTradeCmd
	reqServerPrice    map[uint32]*Cmd.ReqServerPriceRecordTradeCmd
	MapNpcs           map[uint64]*Cmd.MapNpc
	MapUsers          map[uint64]*Cmd.MapUser
	GotoList          *Cmd.GoToListUserCmd
	Mutex             *sync.RWMutex
	lastHeartBeat     time.Time
	retries           map[string]uint
	ExchangeItems     map[uint32]utils.ExchangeItem
	SkillItems        map[uint32]utils.SkillItem
	BuffItems         map[uint32]utils.BuffItem
	BuffItemsByName   map[string]utils.BuffItemByName
	Items             map[uint32]utils.Items
	ItemsByName       map[string]utils.ItemsByName
	notifier          map[string]chan interface{}
}

const (
	queryTimeout              = 15 * time.Second
	printHeartBeatLogInterval = 60 * time.Second
	maxRetry                  = 0
)

var (
	ErrQueryTimeout             = errors.New("query timeout")
	ErrUseClosedConnection      = errors.New("use of closed network connection")
	ErrConnectionClosedByRemote = errors.New("an existing connection was forcibly closed by the remote host")
)

func (g *GameConnection) GameServerLogin() {
	res, err := g.httpAuth(g.Configs.AuthServer)
	if err != nil {
		log.Errorf("failed to get http auth: %s", err)
		log.Exit(4)
	}

	if res.Status == 1001 {
		log.Errorf("failed to get game server access token %v", res)
		log.Exit(1)
	}

	regions := res.Data["regions"].([]interface{})
	regionNum := g.Configs.Region
	id := regions[regionNum].(map[string]interface{})["accid"].(float64)
	sid, err := strconv.ParseUint(regions[regionNum].(map[string]interface{})["serverid"].(string), 10, 32)
	lineGroup := regions[regionNum].(map[string]interface{})["linegroup"].(float64)
	gwPort := regions[regionNum].(map[string]interface{})["gateway_ports"].([]interface{})[0]
	shaStr := regions[regionNum].(map[string]interface{})["sha1"].(string)
	timeStamp, err := strconv.ParseUint(regions[regionNum].(map[string]interface{})["timestamp"].(string), 10, 32)
	g.Configs.AccId = uint64(id)
	g.Configs.ServerId = uint32(lineGroup)
	g.Configs.Sha1Str = shaStr
	g.Configs.Authoriz = res.Data["authorize_state"].(string)
	g.Configs.IpPort = fmt.Sprintf("%s:%s", g.Configs.GameServer, gwPort)
	fmt.Printf("%v, %d ,%s", regions, sid, gwPort)
	g.connectGameServer()
	g.sendServerTimeUserCmd(0)
	g.sendReqUserLoginCmd(uint32(timeStamp))
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

}

func (g *GameConnection) handleConnection() {
	//defer g.conn.Close()
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
					if strings.ContainsAny(err.Error(), ErrUseClosedConnection.Error()) && g.shouldQuit {
						g.Close()
						return
					} else {
						time.Sleep(5 * time.Second)
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
	req, err := http.NewRequest("GET", authHost, nil)
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

func (g *GameConnection) sendServerTimeUserCmd(par Cmd.LoginCmdParam) {
	serverTimeUserCmd := g.getServerTimeUserCmd(0, 0, par)
	data, _ := proto.Marshal(serverTimeUserCmd)
	log.Debug(data)
	out := utils.ConstructBody(1, 11, utils.TcpFlag[1], data, g.getNonce(false), utils.CipherKey)
	log.Infof("sending %v bytes serverTimeUserCmd", len(out))
	_ = g.sendCmd(utils.TcpFlag[1], out, 0)
}

func (g GameConnection) sendProtoCmd(protoCmd proto.Message, cmdId, cmdParId int32) (err error) {
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
	g.Mutex.Lock()
	if time.Since(g.lastCmdSend) <= cmdQueueInterval {
		g.Mutex.Unlock()
		g.cmdQueue = append(g.cmdQueue, newBody)
		go func() {
			time.Sleep(cmdQueueInterval)
			err = g.sendCmd(flag, body, 0)
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
	g.Mutex.Unlock()
	if g.conn != nil {
		writeLen, err := g.conn.Write(newBody)
		log.Debugf("sent %d bytes", writeLen)
		if err != nil {
			log.Errorf("%s failed to send command: %v", g.Role.GetRoleName(), err)
			//g.Close()
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
		Sha1:          &userSha1Str,
		Accid:         &accId,
		Zoneid:        &zId,
		Version:       &version,
		Domain:        &domain,
		Ip:            &ip,
		Device:        &device,
		Language:      &language,
		Deviceid:      &deviceId,
		Clientversion: &clientVer,
		Timestamp:     &timestamp,
		Langzone:      &langZone,
	}
	if serverId != 0 {
		pb.Serverid = &serverId
	}

	if authoriz != "" {
		pb.Authorize = &authoriz
	}

	if safeDevice != "" {
		pb.SafeDevice = &safeDevice
	}

	if phone != "" {
		pb.Phone = &phone
	}
	return pb
}

func (g *GameConnection) enterGameMap() {
	//mid := g.Role.GetMapId()
	//g.ChangeMap(mid)
	g.enteringMap = true
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
	}
	_ = g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["DETAIL_PENDING_LIST_RECORDTRADE"],
	)

	err := g.waitForQueryResponse(itemId)
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
	_ = g.sendProtoCmd(cmd, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["HEART_BEAT_USER_CMD"])
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

func (g *GameConnection) sendReqUserLoginCmd(timeStamp uint32) {
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
	time.Sleep(3 * time.Second)
	_ = g.sendProtoCmd(reqLoginCmd, LogInUserProtoCmdId, Cmd.LoginCmdParam_value["REQ_LOGIN_USER_CMD"])
}

func (g *GameConnection) SelectRole() {
	if g.Role.GetRoleId() != 0 && g.Role.GetRoleName() != "" && g.Role.GetAuthenticated() && g.conn != nil && !g.Role.GetRoleSelected() {
		log.Infof("selecting role with id: %d name: %s", g.Role.GetRoleId(), g.Role.GetRoleName())
		g.doSelectRole()
	}
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

func (g *GameConnection) addNotifier(notifierType string) {
	g.Mutex.Lock()
	g.notifier[notifierType] = make(chan interface{})
	g.Mutex.Unlock()
}

func (g *GameConnection) removeNotifier(notifierType string) {
	g.Mutex.Lock()
	g.notifier[notifierType] = nil
	g.Mutex.Unlock()
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
		Configs:         config,
		Role:            utils.NewRole(),
		currentIndex:    1,
		tradeBrief:      map[uint32]*Cmd.BriefPendingListRecordTradeCmd{},
		tradeDetail:     map[uint32]*Cmd.DetailPendingListRecordTradeCmd{},
		sellInfo:        map[uint32]*Cmd.ItemSellInfoRecordTradeCmd{},
		tradeHistory:    &Cmd.MyTradeLogRecordTradeCmd{},
		buyItem:         map[uint32]*Cmd.BuyItemRecordTradeCmd{},
		sellItem:        map[uint32]*Cmd.SellItemRecordTradeCmd{},
		pendingSells:    &Cmd.MyPendingListRecordTradeCmd{},
		reqServerPrice:  map[uint32]*Cmd.ReqServerPriceRecordTradeCmd{},
		Mutex:           &sync.RWMutex{},
		retries:         map[string]uint{},
		MapNpcs:         map[uint64]*Cmd.MapNpc{},
		MapUsers:        map[uint64]*Cmd.MapUser{},
		SkillItems:      si,
		BuffItems:       bi,
		BuffItemsByName: bin,
		ExchangeItems:   ei,
		Items:           allItems,
		ItemsByName:     allItemsByName,
		notifier:        map[string]chan interface{}{},
	}
	return gc
}
