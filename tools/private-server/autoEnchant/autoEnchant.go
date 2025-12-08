package autoEnchant

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
)

var (
	AttrZhMap = map[Cmd.EAttrType]string{
		// 属性
		Cmd.EAttrType_EATTRTYPE_STR: "力量Str",
		Cmd.EAttrType_EATTRTYPE_AGI: "敏捷Agi",
		Cmd.EAttrType_EATTRTYPE_INT: "智力Int",
		Cmd.EAttrType_EATTRTYPE_VIT: "体质Vit",
		Cmd.EAttrType_EATTRTYPE_DEX: "灵巧Dex",
		Cmd.EAttrType_EATTRTYPE_LUK: "幸运Luk",

		// 基础属性
		Cmd.EAttrType_EATTRTYPE_MAXHP:        "MaxHp",
		Cmd.EAttrType_EATTRTYPE_MAXSP:        "MaxSp",
		Cmd.EAttrType_EATTRTYPE_MAXHPPER:     "MaxHp%",
		Cmd.EAttrType_EATTRTYPE_MAXSPPER:     "MaxSp%",
		Cmd.EAttrType_EATTRTYPE_ATK:          "物理攻击",
		Cmd.EAttrType_EATTRTYPE_MATK:         "魔法攻击",
		Cmd.EAttrType_EATTRTYPE_DEF:          "物理防御",
		Cmd.EAttrType_EATTRTYPE_MDEF:         "魔法防御",
		Cmd.EAttrType_EATTRTYPE_HIT:          "命中",
		Cmd.EAttrType_EATTRTYPE_CRI:          "暴击",
		Cmd.EAttrType_EATTRTYPE_FLEE:         "闪避",
		Cmd.EAttrType_EATTRTYPE_CRIRES:       "暴击防护",
		Cmd.EAttrType_EATTRTYPE_CRIDAMPER:    "暴伤%",
		Cmd.EAttrType_EATTRTYPE_CRIDEFPER:    "暴伤减免%",
		Cmd.EAttrType_EATTRTYPE_HEALENCPER:   "治疗加成%",
		Cmd.EAttrType_EATTRTYPE_BEHEALENCPER: "受治疗加成%",
		Cmd.EAttrType_EATTRTYPE_DAMINCREASE:  "物伤加成",
		Cmd.EAttrType_EATTRTYPE_DAMREDUC:     "物伤减免",
		Cmd.EAttrType_EATTRTYPE_EQUIPASPD:    "装备攻速",

		// 防御属性
		Cmd.EAttrType_EATTRTYPE_SILENCEDEF: "沉默抵抗",
		Cmd.EAttrType_EATTRTYPE_FREEZEDEF:  "冰冻抵抗",
		Cmd.EAttrType_EATTRTYPE_STONEDEF:   "石化抵抗",
		Cmd.EAttrType_EATTRTYPE_STUNDEF:    "眩晕抵抗",
		Cmd.EAttrType_EATTRTYPE_POSIONDEF:  "中毒抵抗",
		Cmd.EAttrType_EATTRTYPE_SLEEPDEF:   "睡眠抵抗",
		Cmd.EAttrType_EATTRTYPE_CHAOSDEF:   "恐惧抵抗",
		Cmd.EAttrType_EATTRTYPE_CURSEDEF:   "诅咒抵抗",
		Cmd.EAttrType_EATTRTYPE_SLOWDEF:    "减速抵抗",
		Cmd.EAttrType_EATTRTYPE_BLINDDEF:   "致盲抵抗",
	}
	AttrMap        = utils.RevertMap(AttrZhMap)
	EnchantTypeMap = map[string]Cmd.EEnchantType{
		"高级": Cmd.EEnchantType_EENCHANTTYPE_SENIOR,
		"中级": Cmd.EEnchantType_EENCHANTTYPE_MEDIUM,
		"低级": Cmd.EEnchantType_EENCHANTTYPE_PRIMARY,
	}
	EnchantEquipPosMap = map[string][]Cmd.EEquipType{
		"武器": {Cmd.EEquipType_EEQUIPTYPE_WEAPON},
		"副手": {
			Cmd.EEquipType_EEQUIPTYPE_SHIELD,
			Cmd.EEquipType_EEQUIPTYPE_BRACELET,
			Cmd.EEquipType_EEQUIPTYPE_EIKON,
			Cmd.EEquipType_EEQUIPTYPE_HANDBRACELET,
			// Cmd.EEquipType_EEQUIPTYPE_PEARL,
		},
		"盔甲":  {Cmd.EEquipType_EEQUIPTYPE_ARMOUR},
		"鞋子":  {Cmd.EEquipType_EEQUIPTYPE_SHOES},
		"披风":  {Cmd.EEquipType_EEQUIPTYPE_ROBE},
		"饰品1": {Cmd.EEquipType_EEQUIPTYPE_ACCESSORY},
		"饰品2": {Cmd.EEquipType_EEQUIPTYPE_ACCESSORY},
		"头饰":  {Cmd.EEquipType_EEQUIPTYPE_HEAD},
		"背部":  {Cmd.EEquipType_EEQUIPTYPE_BACK},
		"尾部":  {Cmd.EEquipType_EEQUIPTYPE_TAIL},
		"脸部":  {Cmd.EEquipType_EEQUIPTYPE_FACE},
		"嘴部":  {Cmd.EEquipType_EEQUIPTYPE_MOUTH},
	}
	ExtraZhMap = []string{
		"名弓1", "名弓2", "名弓3", "名弓4",
		"尖锐1", "尖锐2", "尖锐3", "尖锐4",
		"坚韧1", "坚韧2", "坚韧3", "坚韧4",
		"耐心1", "耐心2", "耐心3", "耐心4",
		"破魔1", "破魔2", "破魔3", "破魔4",
		"利刃1", "利刃2", "利刃3", "利刃4",
		"奥法1", "奥法2", "奥法3", "奥法4",
		"神佑1", "神佑2", "神佑3", "神佑4",
		"铁甲1", "铁甲2", "铁甲3", "铁甲4",
		"狂热1", "狂热2", "狂热3", "狂热4",
		"洞察1", "洞察2", "洞察3", "洞察4",
		"亵渎1", "亵渎2", "亵渎3", "亵渎4",
		"破甲1", "破甲2", "破甲3", "破甲4",
	}
	AllowRoleIds = []uint64{
		// 100100000223,
		// 100100001314,
		// 100100000612,
		// 100100001397,
	}
)

type EnchantTask struct {
	GC        *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	fumoSpeed uint
	fumoCount int
	logWriter io.Writer
	logger    *log.Logger
}

func (e *EnchantTask) GetContext() context.Context {
	return e.ctx
}

func (e *EnchantTask) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(e.GC.LogWriter(), writer)
	e.logWriter = mw
	e.logger.SetOutput(mw)
}

func (e *EnchantTask) Start() {
	e.GC.ShouldChangeScene = true
	e.GC.GameServerLogin()

	if len(AllowRoleIds) > 0 && !utils.Contains(AllowRoleIds, e.GC.Role.GetRoleId()) {
		e.logger.Fatalf("当前角色不在允许列表中，退出")
	}

	_ = e.GC.GetAllPackItems()

	if e.GC.Role.GetMapId() != gameTypes.MapId_Geffen.Uint32() {
		e.logger.Warnf("当前地图不是积芬，飞去积芬中...")
		time.Sleep(time.Second * 5)
		// g.ExitMapWait(gameTypes.MapId_Yuno.Uint32())
		// g.ExitMapWait(gameTypes.MapId_Geffen.Uint32())
		_ = e.GC.GoToGear(gameTypes.MapId_Geffen.Uint32())
		e.GC.ChangeMap(gameTypes.MapId_Geffen.Uint32())
	}

	time.Sleep(time.Second * 3)
	e.logger.Infof("寻找猫小友中...")
	// 猫小友附近
	e.GC.MoveChartWait(e.GC.ParsePos(10739, 2970, 38585))
	err := e.GC.MoveToNpcWait("猫小友")
	if err != nil {
		e.logger.Errorf("没有找到猫小友%s", err)
		return
	}
	_, err = e.GC.VisitNpcByName("猫小友")
	if err != nil {
		e.logger.Errorf("无法对话猫小友 %s", err)
	}

	e.CheckEnchantType()

	e.CheckEnchantEquipPos()

	e.logger.Infof("又来附魔送死了吗?, 来吧来吧, 让我看看是谁不知天高地厚. 附魔类型: %s", e.GC.Configs.EnchantConfig.EnchantType)

	targetEquip := e.GetTargetItem()
	if targetEquip == nil {
		e.logger.Errorf("没有找到要附魔的装备，退出")
		// e.cancel()
	}
	targetEnchant := e.ConditionToEnchantCompare()
	enchantCount := uint(10)
	if e.GC.Configs.EnchantConfig.EnchantCount > 0 {
		enchantCount = uint(e.GC.Configs.EnchantConfig.EnchantCount)
	}
	e.logger.Infof("附魔装备位置: %s", e.GC.Configs.EnchantConfig.EnchantEquipPos)
	e.logger.Infof("目标装备: %s", e.GC.Items[targetEquip.GetBase().GetId()].NameZh)
	e.logger.Infof("附魔停止条件: %v", e.GC.Configs.EnchantConfig.Condition)
	e.logger.Infof("坐稳了要开始附魔了!")
	time.Sleep(5 * time.Second)

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				e.logger.Infof("附魔任务已取消")
				e.GC.Close()
				return
			default:
				e.fumoTask(targetEquip, &targetEnchant, enchantCount)
			}
		}
	}()
}

func (e *EnchantTask) Stop() {
	e.cancel()
}

func (e *EnchantTask) fumoTask(targetEquip *Cmd.ItemData, targetEnchant *gameConnection.EnchantCompare, enchantCount uint) {
	if e.GC.EnchantContains(targetEquip.GetBase().GetGuid(), targetEnchant) && e.GC.Configs.EnchantConfig.AutoSave {
		enchantMap := e.EnchantToZh(targetEquip.GetEnchant())
		e.logger.Infof("已经有附魔要求的属性 %s", FumoStr(enchantMap))

	}
	curCoins := e.GetFuMoBi()
	e.logger.Infof("还有附魔币 %d", curCoins)
	e.logger.Infof("还有神谕之尘 %d", e.GetDust())
	e.logger.Infof("还有神谕之晶 %d", e.GetCrystal())
	e.logger.Infof("第 %d 次 %s附魔 %s", e.fumoCount, e.GC.Configs.EnchantConfig.EnchantType, e.GC.Items[targetEquip.GetBase().GetId()].NameZh)

	// handle auto buy
	leastCoin := uint32(enchantCount * 4)
	if e.GC.Configs.EnchantConfig.AutoBuyCoin.Enable && curCoins <= leastCoin {
		if uint64(e.GC.Configs.EnchantConfig.AutoBuyCoin.MinZenyToKeep) >= e.GC.Role.GetSilver() {
			e.logger.Infof("附魔币不足，但银币低于保留阈值，无法自动购买附魔币，停止附魔")
			return
		}
		e.logger.Infof("附魔币不足，自动购买中...")
		numToBuy := math.Max(float64(e.GC.Configs.EnchantConfig.AutoBuyCoin.NumCoinsToBuy), float64(leastCoin))
		shopConfig, err := e.GC.QueryShopConfig(gameTypes.ShopType_Item, 10)
		if err != nil {
			e.logger.Errorf("购买附魔币查询商店配置失败 %s", err)
		}
		for _, item := range shopConfig.GetGoods() {
			if item.GetId() == 6000 {
				e.logger.Infof("购买%d附魔币", numToBuy)
				e.GC.BuyShopItem(item, uint32(numToBuy))
			}
		}
		time.Sleep(time.Second * 2)
	}

	curEnchant := e.EnchantToZh(targetEquip.GetEnchant())
	e.GC.EnchantEquip(
		EnchantTypeMap[e.GC.Configs.EnchantConfig.EnchantType],
		targetEquip.GetBase().GetGuid(),
		uint32(enchantCount),
	)
	time.Sleep(time.Millisecond * time.Duration(math.Max(float64(e.fumoSpeed), 200)))
	e.logger.Infof("当前附魔: %s", FumoStr(curEnchant))
	targetEquip = e.GetTargetItem()

	previewEnchants := targetEquip.GetPreviewenchant()
	for i, preview := range previewEnchants {
		enchantZh := e.EnchantToZh(preview)
		e.logger.Infof("附魔结果%d: %s", i, FumoStr(enchantZh))
	}
	shouldSave, targetNum := e.GC.EnchantPreviewContains(
		targetEquip.GetBase().GetGuid(),
		targetEnchant,
	)
	if shouldSave && e.GC.Configs.EnchantConfig.AutoSave {
		e.logger.Infof("自動保存附魔属性")
		e.GC.EnchantSave(targetEquip.GetBase().GetGuid(), targetNum)
		time.Sleep(time.Second * 2)
		return // 保存后退出
	} else if shouldSave {
		e.logger.Infof("附魔属性已达到要求，但未保存")
		time.Sleep(time.Second * 2)
		return
	}
	e.fumoCount++
}

func (e *EnchantTask) CheckEnchantType() {
	if e.GC.Configs.EnchantConfig.EnchantType != "" {
		return
	}
	enchantTypes := utils.GetMapKeys(EnchantTypeMap)
	prompt := promptui.Select{
		Label: "请选择要附魔的类型",
		Items: enchantTypes,
	}
	_, result, err := prompt.Run()
	if err != nil {
		e.logger.Errorf("选择附魔类型失败: %s", err)
		return
	}
	e.GC.Configs.EnchantConfig.EnchantType = result
}

func (e *EnchantTask) CheckEnchantEquipPos() {
	if e.GC.Configs.EnchantConfig.EnchantEquipPos != "" {
		return
	}
	enchantEquipPos := utils.GetMapKeys(EnchantEquipPosMap)
	prompt := promptui.Select{
		Label: "请选择要附魔的部位",
		Items: enchantEquipPos,
	}
	_, result, err := prompt.Run()
	if err != nil {
		e.logger.Errorf("选择附魔部位失败: %s", err)
		return
	}
	e.GC.Configs.EnchantConfig.EnchantEquipPos = result
}

func (e *EnchantTask) GetTargetItem() *Cmd.ItemData {
	equipItems := e.GC.Role.GetPackItemsByType(Cmd.EPackType_EPACKTYPE_EQUIP)
	var targetEquip *Cmd.ItemData
	for _, item := range equipItems {
		if utils.Contains(EnchantEquipPosMap[e.GC.Configs.EnchantConfig.EnchantEquipPos], item.GetBase().GetEquipType()) {
			if e.GC.Configs.EnchantConfig.EnchantEquipPos == "饰品1" && item.GetBase().GetIndex() != 5 {
				continue
			} else if e.GC.Configs.EnchantConfig.EnchantEquipPos == "饰品2" && item.GetBase().GetIndex() != 6 {
				continue
			}
			targetEquip = item
			break
		}
	}
	if targetEquip == nil {
		e.logger.Errorf("没有找到要附魔的装备")
		return nil
	}
	return targetEquip
}

func (e *EnchantTask) GetFuMoBi() uint32 {
	item := e.GC.FindPackItemByName("莫拉硬币", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		e.logger.Errorf("没有找到莫拉硬币")
		return 0
	}
	return item.GetBase().GetCount()
}

func (e *EnchantTask) GetCrystal() uint32 {
	item := e.GC.FindPackItemByName("神谕之晶", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		e.logger.Errorf("没有找到神谕之晶")
		return 0
	}
	return item.GetBase().GetCount()
}

func (e *EnchantTask) GetDust() uint32 {
	item := e.GC.FindPackItemByName("神谕之尘", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		e.logger.Errorf("没有找到神谕之尘")
		return 0
	}
	return item.GetBase().GetCount()
}

func (e *EnchantTask) EnchantToZh(data *Cmd.EnchantData) map[string][]string {
	result := make(map[string][]string)
	result["属性"] = []string{}
	result["词条"] = []string{}
	for _, attr := range data.GetAttrs() {
		zhName, ok := AttrZhMap[attr.GetType()]
		if !ok {
			zhName = attr.GetType().String()
		}
		result["属性"] = append(result["属性"],
			fmt.Sprintf("%s +%d",
				zhName,
				attr.GetValue(),
			),
		)
	}
	for _, extra := range data.GetExtras() {
		zhName, ok := e.GC.BuffItems[extra.GetBuffid()]
		if !ok {
			e.logger.Errorf("没有找到词条: %d", extra.GetBuffid())
		}
		result["词条"] = append(result["词条"], zhName.BuffName)
	}
	return result
}

func (e *EnchantTask) ConditionToEnchantCompare() gameConnection.EnchantCompare {
	data := gameConnection.EnchantCompare{}
	for _, attr := range e.GC.Configs.EnchantConfig.Condition.Attributes {
		attrType, value, condition := StringToAttr(attr)
		data.Attrs = append(data.Attrs, &gameConnection.EnchantAttrCompare{
			EnchantAttr: Cmd.EnchantAttr{
				Type:  &attrType,
				Value: &value,
			},
			Condition: condition,
		},
		)
	}
	for _, extra := range e.GC.Configs.EnchantConfig.Condition.Extras {
		ids, ok := e.GC.BuffItemsByName[extra]
		if !ok {
			e.logger.Errorf("没有找到词条: %s", extra)
			continue
		}
		var buffId uint32
		for _, id := range ids.Items {
			if id.BuffName == extra {
				i, _ := id.Id.Int64()
				buffId = uint32(i)
				break
			}
		}

		data.Extras = append(data.Extras, &Cmd.EnchantExtra{
			Buffid: &buffId,
		})
	}
	return data
}

func FumoStr(input map[string][]string) string {
	result := ""
	for k, v := range input {
		switch k {
		case "属性":
			result += "\n属性: "
			for _, attr := range v {
				result += attr + ", "
			}
		case "词条":
			result += "\n词条: "
			for _, extra := range v {
				result += extra + ", "
			}
		}
	}
	return result
}

func StringToAttr(in string) (attrType Cmd.EAttrType, value uint32, condition string) {
	p := strings.Split(in, " ")
	attrType = AttrMap[p[0]]
	v, _ := strconv.ParseUint(p[2], 10, 32)
	value = uint32(v)
	condition = p[1]
	return attrType, value, condition
}

func NewEnchantTask(ctx context.Context, gc *gameConnection.GameConnection, fumoSpeed uint) *EnchantTask {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	return &EnchantTask{
		GC:        gc,
		ctx:       newCtx,
		cancel:    cancel,
		fumoSpeed: fumoSpeed,
		logWriter: mw,
		logger:    logger,
	}
}
