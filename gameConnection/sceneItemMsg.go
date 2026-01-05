package gameConnection

import (
	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"github.com/golang/protobuf/proto"
)

func (g *GameConnection) HandleSceneItemProtoCmd(cmdParamId int32, rawData []byte) (err error) {
	var param proto.Message
	switch cmdParamId {
	case Cmd.ItemParam_value["ITEMPARAM_NTF_HIGHTREFINE_DATA"]:
		param = &Cmd.NtfHighRefineDataCmd{}
		err = utils.ParseCmd(rawData, param)

	case Cmd.ItemParam_value["ITEMPARAM_PACKSLOTNTF"]:
		param = &Cmd.PackSlotNtfItemCmd{}
		err = utils.ParseCmd(rawData, param)

	case Cmd.ItemParam_value["ITEMPARAM_ITEMSHOW"]:
		param = &Cmd.ItemShow{}
		err = utils.ParseCmd(rawData, param)

	case Cmd.ItemParam_value["ITEMPARAM_PACKAGEUPDATE"]:
		param = &Cmd.PackageUpdate{}
		err = utils.ParseCmd(rawData, param)
		packUpdate := param.(*Cmd.PackageUpdate)
		packType := packUpdate.GetType()
		g.Role.Mutex.Lock()
		for _, item := range packUpdate.GetUpdateItems() {
			guid := item.GetBase().GetGuid()
			if g.Role.PackItems[packType] == nil {
				g.Role.PackItems[packType] = map[string]*Cmd.ItemData{}
			}
			g.Role.PackItems[packType][guid] = item
		}
		for _, item := range packUpdate.GetDelItems() {
			guid := item.GetBase().GetGuid()
			delete(g.Role.PackItems[packType], guid)
		}
		g.Role.Mutex.Unlock()

	case Cmd.ItemParam_value["ITEMPARAM_BROWSEPACK"]:
		param = &Cmd.BrowsePackage{}
		err = utils.ParseCmd(rawData, param)

	case Cmd.ItemParam_value["ITEMPARAM_PACKAGEITEM"]:
		param = &Cmd.PackageItem{}
		err = utils.ParseCmd(rawData, param)
		g.Role.Mutex.Lock()
		items := param.(*Cmd.PackageItem)
		if len(items.GetData()) == 0 {
			g.Role.Mutex.Unlock()
		} else {
			if g.Role.PackItems[items.GetType()] == nil {
				g.Role.PackItems[items.GetType()] = map[string]*Cmd.ItemData{}
			}
			for _, data := range items.GetData() {
				g.Role.PackItems[items.GetType()][data.GetBase().GetGuid()] = data
			}
		}
		g.Role.Mutex.Unlock()

	case Cmd.ItemParam_value["ITEMPARAM_QUERY_LOTTERYINFO"]:
		param = &Cmd.QueryLotteryInfo{}
		err = utils.ParseCmd(rawData, param)
		lotteryInfo := param.(*Cmd.QueryLotteryInfo)
		if g.Notifier(gameTypes.NtfType_LotteryQueryInfo) != nil {
			go func() {
				g.Notifier(gameTypes.NtfType_LotteryQueryInfo) <- lotteryInfo
			}()
		}

	case Cmd.ItemParam_value["ITEMPARAM_LOTTERY"]:
		param = &Cmd.LotteryCmd{}
		err = utils.ParseCmd(rawData, param)
		lotteryCmd := param.(*Cmd.LotteryCmd)
		if g.Notifier(gameTypes.NtfType_LotteryCmd) != nil && lotteryCmd.GetCharid() == g.Role.GetRoleId() {
			go func() {
				g.Notifier(gameTypes.NtfType_LotteryCmd) <- lotteryCmd
			}()
		}
	}
	return err
}

func (g *GameConnection) QueryLotteryInfo(lotteryType Cmd.ELotteryType) (lotteryInfo *Cmd.QueryLotteryInfo) {
	cmd := &Cmd.QueryLotteryInfo{
		Type: &lotteryType,
	}
	g.AddNotifier(gameTypes.NtfType_LotteryQueryInfo)
	g.sendProtoCmd(cmd, Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"], Cmd.ItemParam_value["ITEMPARAM_QUERY_LOTTERYINFO"])
	res, err := g.waitForResponse("ITEMPARAM_QUERY_LOTTERYINFO")
	if err != nil {
		g.logger.Errorf("查询抽奖信息失败: %v", err)
		return g.QueryLotteryInfo(lotteryType)
	}
	if res != nil {
		lotteryInfo = res.(*Cmd.QueryLotteryInfo)
		return lotteryInfo
	}
	return nil
}

func (g *GameConnection) GetUserLotteryDailyCount(lotteryType Cmd.ELotteryType) (count uint32) {
	count = 0
	if t, ok := gameTypes.LotteryTypeToVarCountMap[lotteryType]; ok {
		g.Mutex.RLock()
		defer g.Mutex.RUnlock()
		if v, ok := g.Role.UserVars[t]; ok {
			count = v.GetValue()
		}
	} else {
		g.logger.Warnf("未知的抽奖类型: %v", lotteryType.String())
	}
	return count
}

func (g *GameConnection) GetUserMagicLotteryDailyCount() (count uint32) {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	if v, ok := g.Role.UserVars[Cmd.EVarType_EVARTYPE_DAY_LOTTERY_CNT_MAGIC]; ok {
		count = v.GetValue()
	}
	return count
}

func (g *GameConnection) LotteryDraw(lotteryType Cmd.ELotteryType, count, price, ticketId uint32, npcId uint64) *Cmd.LotteryCmd {
	skipA := true
	cmd := &Cmd.LotteryCmd{
		Npcid:    &npcId,
		SkipAnim: &skipA,
		Type:     &lotteryType,
		Count:    &count,
		Price:    &price,
	}
	if ticketId > 0 {
		cmd.Ticket = &ticketId
	}
	g.AddNotifier(gameTypes.NtfType_LotteryCmd)
	_ = g.sendProtoCmdIndex(cmd,
		Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"],
		Cmd.ItemParam_value["ITEMPARAM_LOTTERY"],
		1,
	)
	res, err := g.waitForResponse(gameTypes.NtfType_LotteryCmd)
	if err != nil {
		g.logger.Errorf("等待抽奖失败: %v", err)
	}
	if res != nil {
		lotteryCmd := res.(*Cmd.LotteryCmd)
		return lotteryCmd
	}
	return nil
}

// LotteryRecover 抽奖兑换
func (g *GameConnection) LotteryRecover(npcId uint64, lotteryType Cmd.ELotteryType, recoverItemGuidId []string) {
	cmd := Cmd.LotteryRecoveryCmd{
		Npcid: &npcId,
		Type:  &lotteryType,
		Guids: recoverItemGuidId,
	}
	_ = g.sendProtoCmd(&cmd, Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"], Cmd.ItemParam_value["ITEMPARAM_LOTTERY_RECOVERY"])
}
