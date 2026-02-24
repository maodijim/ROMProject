package config

import (
	"io"
	"os"
	"reflect"

	"ROMProject/data"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	clientVersion = "1.3.1"
)

type EsConfig struct {
	Urls []string `yaml:"urls"`
}

type EnchantConfig struct {
	AutoSave        bool               `yaml:"autoSave" json:"autoSave" label:"自动保存"`
	EnchantType     string             `yaml:"enchantType" json:"enchantType" label:"附魔类型(高级/中级/低级)"`
	EnchantEquipPos string             `yaml:"enchantEquipPos" json:"enchantEquipPos" label:"附魔部位(武器/副手/盔甲/鞋子/披风/头饰/饰品1/饰品2/背部/尾部/脸部/嘴部)"`
	Condition       []EnchantCondition `yaml:"condition" json:"condition" label:"附魔目标条件"`
	AutoBuyCoin     AutoBuyCoinConfig  `yaml:"autoBuyCoin" json:"autoBuyCoin" label:"自动购买附魔材料"`
	EnchantCount    uint32             `yaml:"enchantCount" json:"enchantCount" label:"单次附魔次数"`
	BothCondition   bool               `yaml:"bothCondition" json:"bothCondition" label:"属性和词条同时满足才停止附魔"`
	AllAttrMatch    bool               `yaml:"allAttrMatch" json:"allAttrMatch" label:"所有属性条件全部满足才停止附魔，否则满足一个属性就停止附魔"`
}

func (c *EnchantConfig) UnmarshalYAML(value *yaml.Node) error {
	type NewC EnchantConfig
	var newStruct struct {
		NewC `yaml:",inline"`
	}
	if err := value.Decode(&newStruct); err == nil {
		*c = EnchantConfig(newStruct.NewC)
		return nil
	}
	*c = EnchantConfig(newStruct.NewC)
	var alias struct {
		NewC
		Condition EnchantCondition `yaml:"condition" json:"condition" label:"附魔目标条件"`
	}
	if err := value.Decode(&alias); err != nil {
		return err
	}
	*c = EnchantConfig(newStruct.NewC)
	c.Condition = []EnchantCondition{alias.Condition}
	return nil
}

func (c *EnchantConfig) ParseFromInterface(config map[string]any) EnchantConfig {
	utils.ParseConfigFromInterface(config, c)
	return *c
}

func (c *EnchantConfig) GetDefault() EnchantConfig {
	return EnchantConfig{
		AutoSave:        true,
		EnchantType:     "高级",
		EnchantEquipPos: "",
		Condition: []EnchantCondition{
			{
				Attributes: []string{"暴伤% > 80"},
				Extras:     []string{"尖锐4"},
			},
		},
	}
}

type AutoBuyCoinConfig struct {
	Enable        bool  `yaml:"enable" json:"enable" label:"开启/关闭"`
	MinZenyToKeep int64 `yaml:"minZenyToKeep" json:"minZenyToKeep" label:"金钱下限停止补充附魔币"`
	NumCoinsToBuy int   `yaml:"numCoinsToBuy" json:"numCoinsToBuy" default:"1000" label:"单次购买附魔币数量"`
}

type EnchantCondition struct {
	Attributes []string `yaml:"attributes" json:"attributes" label:"属性(必须跟游戏里面的属性描述一样要不然可能无法识别 比如 '暴伤% > 80')"`
	Extras     []string `yaml:"extras" json:"extras" label:"词条(比如 '尖锐4')"`
}

type HuntBossConfig struct {
	CarryTeam    bool     `yaml:"CarryTeam" json:"CarryTeam" label:"组队一起飞"`
	GameDuration int      `yaml:"GameDuration" json:"GameDuration" label:"狩猎时长(小时) 0为无限制"`
	Mini         []string `yaml:"Mini" json:"Mini" label:"Mini狩猎清单"`
	MVP          []string `yaml:"MVP" json:"MVP" label:"MVP狩猎清单"`
	HMVP         []string `yaml:"HMVP" json:"HMVP" label:"HMVP狩猎清单"`
}
type HuntMonsterConfig struct {
	UseDoubleEXP   bool     `yaml:"UseDoubleEXP" json:"UseDoubleEXP" label:"使用洋洋"`
	TimerFly       int      `yaml:"TimerFly" json:"TimerFly" label:"固定时间使用翅膀(0为不使用)"`
	TargetMonsters []string `yaml:"TargetMonsters" json:"TargetMonsters" label:"狩猎魔物清单，设置all自动全部狩猎"`
	TargetItems    []string `yaml:"TargetItems" json:"TargetItems" label:"狩猎物品清单"`
	Map            string   `yaml:"Map" json:"Map" label:"狩猎地图"`
	NatureType     string   `yaml:"NatureType" json:"NatureType" label:"使用属性类型"`
}
type HuntConfig struct {
	PrepEliteCD       int                `yaml:"PrepEliteCD" json:"PrepEliteCD" label:"备战精英技能冷却时间(秒)"`
	HuntMonsterConfig *HuntMonsterConfig `yaml:"HuntMonsterConfig" json:"HuntMonsterConfig" label:"自动狩猎设置"`
	HuntBossConfig    *HuntBossConfig    `yaml:"HuntBossConfig" json:"HuntBossConfig" label:"自动狩猎Boss设置"`
}

func (c *HuntConfig) ParseFromInterface(config map[string]any) HuntConfig {
	utils.ParseConfigFromInterface(config, c)
	return *c
}

func (c *HuntConfig) Merge(newCfg HuntConfig) {
	cv := reflect.ValueOf(c).Elem()
	nv := reflect.ValueOf(newCfg)

	for i := 0; i < cv.NumField(); i++ {
		oldField := cv.Field(i)
		newField := nv.Field(i)

		switch newField.Kind() {

		// --- Pointer / Slice / Map 需要处理 nil ---
		case reflect.Ptr, reflect.Slice, reflect.Map:
			if newField.IsNil() {
				// ⭐ 新值是 nil → 跳过（不覆盖）
				continue
			}
			oldField.Set(newField)

		// --- Struct → 递归深入 Merge ---
		case reflect.Struct:
			// 若旧值也是 struct，则递归
			mergeMethod := oldField.Addr().MethodByName("Merge")
			if mergeMethod.IsValid() {
				// ⭐ 调用子 struct 的 Merge
				mergeMethod.Call([]reflect.Value{newField})
			} else {
				// 没有 Merge 方法就直接覆盖
				oldField.Set(newField)
			}

		// --- 其他类型（int/string/bool）直接覆盖 ---
		default:
			oldField.Set(newField)
		}
	}
}

func (c *HuntBossConfig) GetDefault() HuntBossConfig {
	return HuntBossConfig{
		CarryTeam:    false,
		GameDuration: 0,
		Mini: []string{
			"狸猫", "蓝疯兔", "波利之王", "摇滚蝗虫", "幽灵波利", "蛙王", "直升机哥布灵", "龙蝇", "流浪之狼",
			"枯树精", "狮鹫兽", "安毕斯", "妖君", "兽人婴儿", "南瓜先生", "半龙人", "草精",
			"鹗枭首领", "爱丽丝女仆", "艾斯恩魔女", "弑神者", "迷幻之王",
		},
		MVP: []string{
			"天使波利", "黄金虫", "恶魔波利", "海盗之王", "海神", "哥布灵首领", "蜂后", "蚁后",
			"皮里恩", "虎王", "俄塞里斯", "月夜猫", "兽人英雄", "犬妖首领", "死灵", "阿特罗斯",
			"兽人酋长", "鹗枭男爵", "血腥骑士", "巴风特", "黑暗之王",
		},
		HMVP: []string{
			"卡仑", "狼外婆",
		},
	}
}

func (c *HuntMonsterConfig) GetDefault() HuntMonsterConfig {
	return HuntMonsterConfig{
		UseDoubleEXP:   true,
		TimerFly:       0,
		TargetMonsters: []string{"all"},
		TargetItems:    []string{},
		Map:            "",
		NatureType:     "",
	}
}

type LotteryConfig struct {
	DrawCount          uint32   `yaml:"drawCount" json:"drawCount" label:"抽奖次数"`
	LotteryType        []string `yaml:"lotteryType" json:"lotteryType" label:"抽奖类型"`
	UseTickets         bool     `yaml:"useTickets" json:"useTickets" label:"使用抽奖券"`
	SellTrash          bool     `yaml:"sellTrash" json:"sellTrash" label:"自动出售垃圾物品(只对宴机抽奖)"`
	SellPoringKingCard bool     `yaml:"sellPoriKingCard" json:"sellPoriKingCard" label:"自动分解国王波利卡片"`
	UseStones          bool     `yaml:"useStones" json:"useStones" label:"使用抽奖金币石"`
	MinStoneToKeep     uint32   `yaml:"minStoneToKeep" json:"minStoneToKeep" label:"保留最少抽奖金币石数量"`
}

func (l *LotteryConfig) UnmarshalYAML(value *yaml.Node) error {
	type NewL LotteryConfig
	var newStruct struct {
		NewL `yaml:",inline"`
	}
	if err := value.Decode(&newStruct); err == nil {
		*l = LotteryConfig(newStruct.NewL)
		return nil
	}
	*l = LotteryConfig(newStruct.NewL)
	var alias struct {
		NewL
		LotteryType string `yaml:"lotteryType" json:"lotteryType" label:"抽奖类型"`
	}
	if err := value.Decode(&alias); err != nil {
		return err
	}
	*l = LotteryConfig(newStruct.NewL)
	l.LotteryType = []string{alias.LotteryType}
	return nil
}

func (l *LotteryConfig) ParseFromInterface(config map[string]interface{}) any {
	utils.ParseConfigFromInterface(config, l)
	return *l
}

func (l *LotteryConfig) GetDefault() any {
	return LotteryConfig{
		DrawCount:          10,
		LotteryType:        []string{"幻想创造器·宴"},
		UseTickets:         false,
		SellTrash:          true,
		SellPoringKingCard: false,
		UseStones:          false,
		MinStoneToKeep:     0,
	}
}

type DailyTaskConfig struct {
	EnableItemCombine bool `yaml:"enableItemCombine" json:"enableItemCombine" label:"自动物品合成(精装卡册的残页, 卡册残页)"`
	EnableKanBan      bool `yaml:"enableKanBan" json:"enableKanBan" label:"完成看板任务"`
	// EnableWasteLandWeed       bool   `yaml:"enableWasteLandWeed" json:"enableWasteLandWeed" label:"完成清理荒地杂草"`
	EnableCrack               bool   `yaml:"enableCrack" json:"enableCrack" label:"完成裂缝任务"`
	EnableYuno                bool   `yaml:"enableZhuno" json:"enableZhuno" label:"完成朱诺任务"`
	EnableGuildEmperiumDonate bool   `yaml:"enableGuildEmperiumDonate" json:"enableGuildEmperiumDonate" label:"完成公会华丽金属捐献任务"`
	YunoTeamLeader            string `yaml:"YunoTeamLeader" json:"YunoTeamLeader" label:"朱诺任务队长名称(自动组队使用)"`
	YunForceContinue          bool   `yaml:"YunForceContinue" json:"YunForceContinue" label:"朱诺任务强制继续(即使已经完成了每日次数)"`
	PurchaseDailyBag          bool   `yaml:"purchaseDailyBag" json:"purchaseDailyBag" label:"自动购买每日福袋"`
	PurchaseDailyZeny         bool   `yaml:"purchaseDailyZeny" json:"purchaseDailyZeny" label:"自动购买每日Zeny"`
	PurchaseResurrection      bool   `yaml:"purchaseResurrection" json:"purchaseResurrection" label:"购买不死之证"`
	PurchaseDiamondZeny       bool   `yaml:"purchaseDiamondZeny" json:"purchaseDiamondZeny" label:"购买初心币Zeny"`
	PurchaseWeedPackage       bool   `yaml:"purchaseWeedPackage" json:"purchaseWeedPackage" label:"购买荒境除草卡片礼包"`
}

func (d *DailyTaskConfig) ParseFromInterface(config map[string]interface{}) any {
	utils.ParseConfigFromInterface(config, d)
	return *d
}

func (d *DailyTaskConfig) GetDefault() any {
	return DailyTaskConfig{
		EnableItemCombine: true,
		EnableKanBan:      true,
		// EnableWasteLandWeed:       true,
		EnableCrack:               true,
		EnableYuno:                true,
		EnableGuildEmperiumDonate: true,
		YunForceContinue:          false,
	}
}

type ServerConfigs struct {
	AuthServer     string `yaml:"authServer"`
	AuthPass       string `yaml:"authPass"`
	ZoneId         uint32 `yaml:"zoneId"`
	Char           uint   `yaml:"char"`
	AutoCreateChar bool   `yaml:"autoCreateChar"`
	// if not set, use random string during auto create char
	CharacterName string        `yaml:"characterName"`
	EnchantConfig EnchantConfig `yaml:"enchantConfig"`
	HuntConfig    HuntConfig    `yaml:"HuntConfig"`
	// Username and password for login if accId is not set
	Username           string             `yaml:"username"`
	Password           string             `yaml:"password"`
	Region             int                `yaml:"region"`
	Version            string             `yaml:"version"`
	ServerId           uint32             `yaml:"serverId"`
	AccId              uint64             `yaml:"accId"`
	Ip                 string             `yaml:"ip"`
	Domain             string             `yaml:"domain"`
	LineGrp            string             `yaml:"lineGrp"`
	Device             string             `yaml:"device"`
	DeviceId           string             `yaml:"deviceId"`
	ClientVer          string             `yaml:"clientVer"`
	LangZone           uint32             `yaml:"langZone"`
	IpPort             string             `yaml:"ipPort"`
	GameServer         string             `yaml:"gameServer"`
	Phone              string             `yaml:"phone"`
	SafeDevice         string             `yaml:"safeDevice"`
	Sha1Str            string             `yaml:"sha1Str"`
	AccessToken        string             `yaml:"accessToken"`
	ResVer             string             `yaml:"resourceVer"`
	PlatformVer        string             `yaml:"platV"`
	AppPreVer          uint32             `yaml:"appPreVer"`
	Lang               uint32             `yaml:"lang"`
	Model              string             `yaml:"model"`
	PhoneVer           string             `yaml:"phoneVer"`
	Authoriz           string             `yaml:"authoriz"`
	AuthParams         map[string]string  `yaml:"authParams"`
	EsConfig           EsConfig           `yaml:"elasticsearch"`
	TeamConfig         TeamConfig         `yaml:"team"`
	Loglines           int                `yaml:"-"`
	ChatMaxSize        int                `yaml:"-"`
	TradeMonitorConfig TradeMonitorConfig `yaml:"tradeMonitorConfig"`
	LotteryConfig      LotteryConfig      `yaml:"lotteryConfig"`
	DailyTaskConfig    DailyTaskConfig    `yaml:"dailyTaskConfig"`
}

func (s *ServerConfigs) UnmarshalYAML(value *yaml.Node) error {
	type NewS ServerConfigs
	var newStruct struct {
		NewS `yaml:",inline"`
	}
	if err := value.Decode(&newStruct); err != nil {
		return err
	}
	*s = ServerConfigs(newStruct.NewS)
	if s.Version == "" {
		s.Version = clientVersion
	} else if s.Version < clientVersion {
		log.Warnf("Client version in config (%s) is older than the default client version (%s). Consider updating it to avoid potential issues.", s.Version, clientVersion)
		s.Version = clientVersion
	}
	return nil
}

func (s *ServerConfigs) GetChatMaxSize() int {
	if s.ChatMaxSize <= 0 {
		return 500
	}
	return s.ChatMaxSize
}

func (s *ServerConfigs) SetTeamLeader(name string) {
	s.TeamConfig.LeaderName = name
}

func (s *ServerConfigs) SetFollowTeamLeader(yes bool) {
	s.TeamConfig.FollowTeamLeader = yes
}

func parseConfigYaml(r []byte, sc *ServerConfigs) error {
	var newSc ServerConfigs
	if err := yaml.Unmarshal(r, &newSc); err != nil {
		return err
	}
	// merge newSc into sc but keeping the original value if the new value is zero value (like 0 for int, "" for string, false for bool, nil for pointer/slice/map)
	scValue := reflect.ValueOf(sc).Elem()
	newScValue := reflect.ValueOf(&newSc).Elem()
	for i := 0; i < scValue.NumField(); i++ {
		oldField := scValue.Field(i)
		newField := newScValue.Field(i)

		switch newField.Kind() {

		// --- Pointer / Slice / Map 需要处理 nil ---
		case reflect.Ptr, reflect.Slice, reflect.Map:
			if newField.IsNil() {
				// ⭐ 新值是 nil → 跳过（不覆盖）
				continue
			}
			oldField.Set(newField)

		// --- Struct → 递归深入 Merge ---
		case reflect.Struct:
			// 若旧值也是 struct，则递归
			mergeMethod := oldField.Addr().MethodByName("Merge")
			if mergeMethod.IsValid() {
				// ⭐ 调用子 struct 的 Merge
				mergeMethod.Call([]reflect.Value{newField})
			} else {
				// 没有 Merge 方法就直接覆盖
				oldField.Set(newField)
			}

		case reflect.String:
			if newField.String() == "" {
				// ⭐ 新值是空字符串 → 跳过（不覆盖）
				continue
			}
			oldField.Set(newField)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if newField.Int() == 0 {
				continue
			}
			oldField.Set(newField)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if newField.Uint() == 0 {
				continue
			}
			oldField.Set(newField)

		case reflect.Float32, reflect.Float64:
			if newField.Float() == 0 {
				continue
			}
			oldField.Set(newField)
			// ⭐ 新值是0 → 跳过（不覆盖）
			continue

		// --- 其他类型（int/string/bool）直接覆盖 ---
		default:
			oldField.Set(newField)
		}
	}

	if sc.Region < 1 {
		sc.Region = 1
	} else {
		sc.Region -= 1
	}

	if sc.EnchantConfig.AutoBuyCoin.Enable {
		if sc.EnchantConfig.AutoBuyCoin.NumCoinsToBuy <= 0 {
			sc.EnchantConfig.AutoBuyCoin.NumCoinsToBuy = 1000
		}
		if sc.EnchantConfig.AutoBuyCoin.MinZenyToKeep <= 0 {
			sc.EnchantConfig.AutoBuyCoin.MinZenyToKeep = 100000000
		}
	}

	return nil
}

func NewServerConfigs(configYaml string) *ServerConfigs {
	configPath := configYaml
	configs := &ServerConfigs{}
	err := parseConfigYaml(data.ConfigYml, configs)
	if err != nil {
		log.Fatalf("failed to parse default config yaml: %v", err)
	}
	if configYaml == "" {
		configPath = "config.yml"
	}
	f, err := os.Open(configYaml)
	if err != nil {
		log.Errorf("failed to load %s: %s", configPath, err)
		return configs
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		log.Errorf("failed to read %s: %s", configPath, err)
		return configs
	}
	err = parseConfigYaml(content, configs)
	if err != nil {
		log.Errorf("parse config yaml failed: %s", err)
		return configs
	}

	return configs
}
