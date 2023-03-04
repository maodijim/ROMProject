package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
)

const (
	ver = "0.1.1"
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
	g            *gameConnection.GameConnection
	allowRoleIds = []uint64{
		// 100100000223,
		// 100100001314,
		// 100100000612,
		// 100100001397,
	}
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logFile := "autoEnchant.log"
	var mw io.Writer
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw = io.MultiWriter(os.Stdout, f)
	} else {
		log.Warnf("无法写入日志文件 %s，使用默认输出", logFile)
		mw = os.Stdout
	}
	log.SetOutput(mw)
}

func main() {
	log.Infof("自动附魔版本 %s", ver)
	configPath := flag.String("config", "config.yml", "配置文件路径")
	enableDebug := flag.Bool("debug", false, "是否开启调试模式")
	speed := flag.Uint("speed", 850, "附魔速度，单位毫秒")
	autoSave := flag.Bool("autoSave", false, "是否自动保存 (默认不保存)")
	flag.Parse()
	items := utils.NewItemsLoader("", "", "")
	conf := config.NewServerConfigs(*configPath)
	skills := utils.NewSkillParser("")
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster("")
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}
	g.GameServerLogin()

	if len(allowRoleIds) > 0 && !utils.Contains(allowRoleIds, g.Role.GetRoleId()) {
		log.Fatalf("当前角色不在允许列表中，退出")
	}

	_ = g.GetAllPackItems()

	if g.Role.GetMapId() != gameTypes.MapId_Geffen.Uint32() {
		log.Warnf("当前地图不是积芬，飞去积芬中...")
		time.Sleep(time.Second * 5)
		g.ExitMapWait(gameTypes.MapId_Yuno.Uint32())
		g.ExitMapWait(gameTypes.MapId_Geffen.Uint32())
	}
	time.Sleep(time.Second * 3)
	log.Infof("寻找猫小友中...")
	// 猫小友附近
	g.MoveChartWait(g.ParsePos(12824, 2912, 37016))
	err := g.MoveToNpcWait("猫小友")
	if err != nil {
		log.Errorf("没有找到猫小友%s", err)
		return
	}
	_, err = g.VisitNpcByName("猫小友")
	if err != nil {
		log.Errorf("无法对话猫小友 %s", err)
	}

	checkEnchantType()

	checkEnchantEquipPos()

	log.Infof("又来附魔送死了吗?, 来吧来吧, 让我看看是谁不知天高地厚. 附魔类型: %s", g.Configs.EnchantConfig.EnchantType)

	targetEquip := getTargetItem()
	targetEnchant := conditionToEnchantCompare()
	log.Infof("附魔装备位置: %s", g.Configs.EnchantConfig.EnchantEquipPos)
	log.Infof("目标装备: %s", g.Items[targetEquip.GetBase().GetId()].NameZh)
	log.Infof("附魔停止条件: %v", g.Configs.EnchantConfig.Condition)
	log.Infof("坐稳了要开始附魔了!")
	time.Sleep(5 * time.Second)
	count := 1
	for {
		if g.EnchantContains(targetEquip.GetBase().GetGuid(), &targetEnchant) && *autoSave {
			enchantMap := enchantToZh(targetEquip.GetEnchant())
			log.Infof("已经有附魔要求的属性 %s", fumoStr(enchantMap))
			break
		}
		log.Infof("还有附魔币 %d", getFuMoBi())
		log.Infof("还有神谕之尘 %d", getDust())
		log.Infof("还有神谕之晶 %d", getCrystal())
		log.Infof("第 %d 次 %s附魔 %s", count, g.Configs.EnchantConfig.EnchantType, g.Items[targetEquip.GetBase().GetId()].NameZh)
		curEnchant := enchantToZh(targetEquip.GetEnchant())
		g.EnchantEquip(
			EnchantTypeMap[g.Configs.EnchantConfig.EnchantType],
			targetEquip.GetBase().GetGuid(),
		)
		time.Sleep(time.Millisecond * time.Duration(math.Max(float64(*speed), 200)))
		log.Infof("当前附魔: %s", fumoStr(curEnchant))
		targetEquip = getTargetItem()
		previewEnchant := enchantToZh(targetEquip.GetPreviewenchant())
		log.Infof("附魔结果: %s", fumoStr(previewEnchant))
		shouldSave := g.EnchantPreviewContains(
			targetEquip.GetBase().GetGuid(),
			&targetEnchant,
		)
		if shouldSave && *autoSave {
			log.Infof("自動保存附魔属性")
			g.EnchantSave(targetEquip.GetBase().GetGuid())
			time.Sleep(time.Second * 2)
			return // 保存后退出
		} else if shouldSave {
			log.Infof("附魔属性已达到要求，但未保存")
			time.Sleep(time.Second * 2)
			return
		}
		count++
	}
}

func checkEnchantType() {
	if g.Configs.EnchantConfig.EnchantType != "" {
		return
	}
	enchantTypes := utils.GetMapKeys(EnchantTypeMap)
	prompt := promptui.Select{
		Label: "请选择要附魔的类型",
		Items: enchantTypes,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Errorf("选择附魔类型失败: %s", err)
		return
	}
	g.Configs.EnchantConfig.EnchantType = result
}

func checkEnchantEquipPos() {
	if g.Configs.EnchantConfig.EnchantEquipPos != "" {
		return
	}
	enchantEquipPos := utils.GetMapKeys(EnchantEquipPosMap)
	prompt := promptui.Select{
		Label: "请选择要附魔的部位",
		Items: enchantEquipPos,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Errorf("选择附魔部位失败: %s", err)
		return
	}
	g.Configs.EnchantConfig.EnchantEquipPos = result
}

func getTargetItem() *Cmd.ItemData {
	equipItems := g.Role.GetPackItemsByType(Cmd.EPackType_EPACKTYPE_EQUIP)
	var targetEquip *Cmd.ItemData
	for _, item := range equipItems {
		if utils.Contains(EnchantEquipPosMap[g.Configs.EnchantConfig.EnchantEquipPos], item.GetBase().GetEquipType()) {
			if g.Configs.EnchantConfig.EnchantEquipPos == "饰品1" && item.GetBase().GetIndex() != 5 {
				continue
			} else if g.Configs.EnchantConfig.EnchantEquipPos == "饰品2" && item.GetBase().GetIndex() != 6 {
				continue
			}
			targetEquip = item
			break
		}
	}
	if targetEquip == nil {
		log.Fatalf("没有找到要附魔的装备")
		return nil
	}
	return targetEquip
}

func getFuMoBi() uint32 {
	item := g.FindPackItemByName("莫拉硬币", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		log.Errorf("没有找到莫拉硬币")
		return 0
	}
	return item.GetBase().GetCount()
}

func getCrystal() uint32 {
	item := g.FindPackItemByName("神谕之晶", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		log.Errorf("没有找到神谕之晶")
		return 0
	}
	return item.GetBase().GetCount()
}

func getDust() uint32 {
	item := g.FindPackItemByName("神谕之尘", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		log.Errorf("没有找到神谕之尘")
		return 0
	}
	return item.GetBase().GetCount()
}

func enchantToZh(data *Cmd.EnchantData) map[string][]string {
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
		zhName, ok := g.BuffItems[extra.GetBuffid()]
		if !ok {
			log.Errorf("没有找到词条: %d", extra.GetBuffid())
		}
		result["词条"] = append(result["词条"], zhName.BuffName)
	}
	return result
}

func fumoStr(input map[string][]string) string {
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

func conditionToEnchantCompare() gameConnection.EnchantCompare {
	data := gameConnection.EnchantCompare{}
	for _, attr := range g.Configs.EnchantConfig.Condition.Attributes {
		attrType, value, condition := stringToAttr(attr)
		data.Attrs = append(data.Attrs, &gameConnection.EnchantAttrCompare{
			EnchantAttr: Cmd.EnchantAttr{
				Type:  &attrType,
				Value: &value,
			},
			Condition: condition,
		},
		)
	}
	for _, extra := range g.Configs.EnchantConfig.Condition.Extras {
		ids, ok := g.BuffItemsByName[extra]
		if !ok {
			log.Errorf("没有找到词条: %s", extra)
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

func stringToAttr(in string) (attrType Cmd.EAttrType, value uint32, condition string) {
	p := strings.Split(in, " ")
	attrType = AttrMap[p[0]]
	v, _ := strconv.ParseUint(p[2], 10, 32)
	value = uint32(v)
	condition = p[1]
	return attrType, value, condition
}
