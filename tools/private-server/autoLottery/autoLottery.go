package autoLottery

import (
	"context"
	"io"
	"os"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"

	log "github.com/sirupsen/logrus"
)

type LotteryTask struct {
	GC             *gameConnection.GameConnection
	ctx            context.Context
	cancel         context.CancelFunc
	logWriter      io.Writer
	logger         *log.Logger
	completeStatus map[string]bool
}

func (l *LotteryTask) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(l.GC.LogWriter(), writer)
	l.logWriter = mw
	l.logger.SetOutput(mw)
}

func (l *LotteryTask) Start() {
	l.GC.GameServerLogin()

	if l.GC.Role.GetMapId() != gameTypes.MapId_Protera.Uint32() {
		l.logger.Warnf("当前地图不是普隆德拉，飞去普隆德拉中...")
		time.Sleep(time.Second * 5)
		l.GC.GoToMap(gameTypes.MapId_Protera.Uint32())
	}

	time.Sleep(time.Second * 5)

	// 移动到抽奖NPC位置
	l.logger.Infof("移动到抽奖NPC位置...")
	time.Sleep(time.Second * 2)
	_ = l.GC.MoveChartWait(l.GC.ParsePos(-25099, -513, -46410))
	time.Sleep(time.Second * 2)
	lotteryName := l.GC.Configs.LotteryConfig.LotteryType
	l.GC.MoveToNpcWait(lotteryName)
	time.Sleep(time.Second)
	npc, err := l.GC.VisitObjectByName(lotteryName)
	if err != nil {
		l.logger.Errorf("访问NPC失败: %v", err)
		return
	}
	time.Sleep(time.Second * 2)

	// 开始抽奖任务
	l.lotteryTask(npc)

	if l.GC.Configs.LotteryConfig.SellPoringKingCard {
		l.logger.Infof("开始分解波利国王卡片...")
		l.DecomposePoringKingCard()
	}
}

func (l *LotteryTask) Stop() {
	l.cancel()
}

func (l *LotteryTask) GetContext() context.Context {
	return l.ctx
}

func (l *LotteryTask) lotteryTask(npc Cmd.MapNpc) {
	l.logger.Infof("开始自动抽奖任务...")
	l.GC.SetQueryTimeout(time.Millisecond * 1000)
	lotteryName := l.GC.Configs.LotteryConfig.LotteryType
	lotteryType, ok := gameTypes.LotteryNameZh[l.GC.Configs.LotteryConfig.LotteryType]
	if !ok {
		l.logger.Errorf("未知的抽奖类型: %s", l.GC.Configs.LotteryConfig.LotteryType)
		return
	}
	maxCount := uint32(0)
	dailyCount := uint32(0)
	lotteryPrice := gameTypes.LotteryTypePriceMap[lotteryType]
	var lotteryInfo *Cmd.QueryLotteryInfo

	lotteryInfo = l.GC.QueryLotteryInfo(lotteryType)
	dailyCount = lotteryInfo.GetTodayCnt()
	maxCount = lotteryInfo.GetMaxCnt()

	if l.GC.Configs.LotteryConfig.SellTrash {
		l.SellTrash(*lotteryInfo, npc.GetId(), lotteryType)
	}

lotteryLoop:
	for {
		select {
		case <-l.ctx.Done():
			l.logger.Infof("自动抽奖任务已停止.")
			return
		default:
			// 检查抽奖机会
			if dailyCount == 0 && lotteryInfo == nil {
				l.logger.Infof("检查抽奖机会...")
				lotteryInfo = l.GC.QueryLotteryInfo(lotteryType)
				dailyCount = lotteryInfo.GetTodayCnt()
				maxCount = lotteryInfo.GetMaxCnt()
			}

			if dailyCount >= maxCount && !l.GC.Configs.LotteryConfig.UseTickets {
				l.logger.Infof("今日%s抽奖机会%d次已用完", lotteryName, maxCount)
				break lotteryLoop
			} else {
				l.logger.Infof("今日%s抽奖机会剩余: %d次", lotteryName, maxCount-dailyCount)
			}
			// 在这里添加具体的抽奖逻辑
			drawCount := min(l.GC.Configs.LotteryConfig.DrawCount, maxCount-dailyCount)
			l.logger.Infof("开始抽奖: 抽取 %d 次 %s", drawCount, lotteryName)
			l.logger.Infof("剩余票券: %d", l.GetTicketCount(*lotteryInfo))
			ticketId := uint32(0)
			if l.GC.Configs.LotteryConfig.UseTickets {
				ticketId = lotteryInfo.GetInfos()[0].GetSubInfo()[0].GetRecoverItemid()
				drawCount = min(drawCount, l.GetTicketCount(*lotteryInfo)/30)
				if drawCount == 0 {
					l.logger.Infof("票券不足，无法继续抽奖。")
					l.Stop()
					break lotteryLoop
				}
			} else {
				drawCount = min(drawCount, uint32(l.GC.Role.GetLottery()/lotteryPrice))
				if drawCount == 0 {
					l.logger.Infof("猫币不足，无法继续抽奖。")
					break lotteryLoop
				}
			}

			lotteryCmd := l.GC.LotteryDraw(lotteryType, drawCount, uint32(lotteryPrice)*drawCount, ticketId, npc.GetId())
			l.logger.Infof("抽奖结果: 获得 %d 个物品", len(lotteryCmd.GetItems()))
			if lotteryCmd != nil && lotteryCmd.GetTodayCnt() > 0 {
				dailyCount = lotteryCmd.GetTodayCnt()
			}
		}
	}

	// 使用银币宝石
	if !l.GC.Configs.LotteryConfig.UseStones {
		l.logger.Infof("跳过使用银币宝石。")
		return
	}

	l.logger.Infof("开始使用银币宝石...")
	l.completeStatus = map[string]bool{
		"红色玛瑙": false,
		"黑珍珠":  false,
		"金之星":  false,
	}

	for {
		select {
		case <-l.ctx.Done():
			l.logger.Infof("使用银币宝石已停止.")
			return
		default:
			if l.completeStatus["红色玛瑙"] && l.completeStatus["黑珍珠"] && l.completeStatus["金之星"] {
				l.logger.Infof("所有银币宝石使用完成。")
				return
			}
			l.UseSilverGem()
		}
	}
}

func (l *LotteryTask) GetTicketCount(lotteryInfo Cmd.QueryLotteryInfo) (totalCount uint32) {
	if lotteryInfo.GetInfos() == nil {
		return 0
	}
	if lotteryInfo.GetInfos()[0].GetSubInfo() == nil {
		return 0
	}
	ticketId := lotteryInfo.GetInfos()[0].GetSubInfo()[0].GetRecoverItemid()
	item := l.GC.FindPackItemByIdAll(ticketId, Cmd.EPackType_EPACKTYPE_MAIN)
	for _, it := range item {
		totalCount += it.GetBase().GetCount()
	}
	return totalCount
}

func (l *LotteryTask) SellTrash(lotteryInfo Cmd.QueryLotteryInfo, npcId uint64, lotteryType Cmd.ELotteryType) {
	l.logger.Infof("开始出售垃圾物品...")
	infos := lotteryInfo.GetInfos()
	if len(infos) == 0 {
		l.logger.Infof("没有可出售的垃圾物品信息。")
		return
	}
	subInfos := infos[0].GetSubInfo()
	if len(subInfos) == 0 {
		l.logger.Infof("没有可出售的垃圾物品子信息。")
		return
	}
	for _, subInfo := range subInfos {
		recoverId := subInfo.GetRecoverItemid()
		if recoverId == 0 {
			continue
		}
		trashItemId := subInfo.GetItemid()
		items := l.GC.FindPackItemByIdAll(trashItemId, Cmd.EPackType_EPACKTYPE_MAIN)
		var itemGuids []string
		var totalCount int
		for _, item := range items {
			totalCount += int(item.GetBase().GetCount())
			itemGuids = append(itemGuids, item.GetBase().GetGuid())
		}
		itemName, ok := l.GC.Items[trashItemId]
		var itemNameStr string
		if !ok {
			itemNameStr = "未知物品名"
		} else {
			itemNameStr = itemName.NameZh
		}
		if totalCount == 0 {
			continue
		}
		l.logger.Infof("找到 %d 个垃圾物品 (ID: %d, %s)，开始出售...", totalCount, trashItemId, itemNameStr)
		l.GC.LotteryRecover(npcId, lotteryType, itemGuids)
		time.Sleep(350 * time.Millisecond)
		ticketCount := l.GetTicketCount(lotteryInfo)
		l.logger.Infof("当前票券数量: %d", ticketCount)
	}
}

func (l *LotteryTask) DecomposePoringKingCard() {
	cardItem := l.GC.FindPackItemByName("国王波利的恩惠", Cmd.EPackType_EPACKTYPE_MAIN)
	cardCount := cardItem.GetBase().GetCount()
	if cardCount == 0 {
		l.logger.Infof("背包中没有国王波利的恩惠卡片，跳过分解。")
		return
	}
	l.logger.Infof("移动到波利国王...")
	// 移动到抽奖NPC位置
	l.GC.MoveToNpcWait("恶魔波利")
	time.Sleep(time.Second)
	npc, err := l.GC.VisitObjectByName("恶魔波利")
	if err != nil {
		l.logger.Errorf("访问NPC失败: %v", err)
		return
	}
	time.Sleep(time.Second * 2)

	l.logger.Infof("开始分解国王波利的恩惠卡片...")
	for {
		select {
		case <-l.ctx.Done():
			l.logger.Infof("分解国王波利的恩惠卡片任务已停止.")
			return
		default:
			if cardCount == 0 {
				l.logger.Infof("背包中没有国王波利的恩惠卡片，分解任务完成。")
				return
			}
			cardGuid := cardItem.GetBase().GetGuid()
			cardCount = cardItem.GetBase().GetCount()
			if l.GC.Role.GetSilver() < uint64(cardCount*10000) {
				l.logger.Infof("银币不足，无法继续分解国王波利的恩惠卡片。")
				return
			}
			cardList := make([]string, 0)
			for i := uint32(0); i < min(cardCount, 50); i++ {
				cardList = append(cardList, cardGuid)
			}
			res, err := l.GC.ExchangeCardDecompose(npc.GetId(), cardList...)
			if err != nil {
				l.logger.Errorf("分解国王波利的恩惠卡片失败: %v", err)
				return
			}
			items := res.GetItems()
			for _, item := range items {
				itemName := l.GC.FindItemNameById(item.GetId())
				l.logger.Infof("分解国王波利的恩惠成功，获得以下物品: %s %d个", itemName, item.GetCount())
			}
			time.Sleep(time.Second)
			cardItem = l.GC.FindPackItemByName("国王波利的恩惠", Cmd.EPackType_EPACKTYPE_MAIN)
		}
	}
}

func (l *LotteryTask) UseSilverGem() {
	gem1 := l.GC.FindPackItemByName("红色玛瑙", Cmd.EPackType_EPACKTYPE_MAIN)
	gem2 := l.GC.FindPackItemByName("黑珍珠", Cmd.EPackType_EPACKTYPE_MAIN)
	gem3 := l.GC.FindPackItemByName("金之星", Cmd.EPackType_EPACKTYPE_MAIN)

	l.useGem(gem1, "红色玛瑙")
	l.useGem(gem2, "黑珍珠")
	l.useGem(gem3, "金之星")
}

func (l *LotteryTask) useGem(item *Cmd.ItemData, gemName string) {
	if item.GetBase().GetCount() > 1 {
		useCount := min(99, item.GetBase().GetCount())
		if item.GetBase().GetCount() < l.GC.Configs.LotteryConfig.MinStoneToKeep {
			l.logger.Infof("%s保留数量不足，跳过使用。", gemName)
			l.completeStatus[gemName] = true
			return
		}
		l.logger.Infof("%s剩下%d个使用%d个...", gemName, item.GetBase().GetCount(), useCount)
		l.GC.UseItem(item.GetBase().GetGuid(), useCount)
		time.Sleep(time.Millisecond * 1500)
	} else {
		l.logger.Infof("%s数量不足，跳过使用。", gemName)
		l.completeStatus[gemName] = true
	}
}

func NewLotteryTask(ctx context.Context, gc *gameConnection.GameConnection) *LotteryTask {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	return &LotteryTask{
		GC:        gc,
		ctx:       newCtx,
		cancel:    cancel,
		logWriter: mw,
		logger:    logger,
	}
}
