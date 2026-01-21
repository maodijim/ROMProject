package main

import (
	"encoding/json"
	"net/http"
	"sort"

	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/tools/private-server/autoEnchant"
	"ROMProject/utils"
)

func GetMiniList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"狸猫",
		"蓝疯兔",
		"波利之王",
		"摇滚蝗虫",
		"幽灵波利",
		"蛙王",
		"直升机哥布灵",
		"龙蝇",
		"流浪之狼",
		"枯树精",
		"狮鹫兽",
		"安毕斯",
		"妖君",
		"兽人婴儿",
		"南瓜先生",
		"半龙人",
		"草精",
		"鹗枭首领",
		"爱丽丝女仆",
		"艾斯恩魔女",
		"弑神者",
		"迷幻之王",
		"大笨钟",
		"钟塔守护者",
		"魔灵娃娃",
		"炎之小魔女",
		"吹笛人",
		"银月魔女",
		"暗赛尼亚",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetMVPList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"天使波利",
		"黄金虫",
		"恶魔波利",
		"海盗之王",
		"海神",
		"哥布灵首领",
		"蜂后",
		"蚁后",
		"皮里恩",
		"虎王",
		"俄塞里斯",
		"月夜猫",
		"兽人英雄",
		"犬妖首领",
		"死灵",
		"阿特罗斯",
		"兽人酋长",
		"迪塔勒泰晤勒",
		"鹗枭男爵",
		"血腥骑士",
		"巴风特",
		"黑暗之王",
		"时间管理人",
		"斯佩夏尔",
		"冰暴骑士",
		"炎之领主卡浩",
		"圣天使波利",
		"暗·神射手迪文",
		"暗·超魔导师凯特莉娜",
		"暗·十字刺客艾勒梅斯",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetHMVPList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"卡仑",
		"狼外婆",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetMAPList(w http.ResponseWriter, r *http.Request) {
	ids := make([]int, 0, len(gameTypes.MapIdToZh))
	for id := range gameTypes.MapIdToZh {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	list := make([]string, 0, len(ids))
	for _, id := range ids {
		list = append(list, gameTypes.MapIdToZh[gameTypes.MapId(id)])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetNatureList(w http.ResponseWriter, r *http.Request) {
	list := make([]string, 0, len(gameTypes.NatureTypeZhMap))

	// 1. 先取得所有 key
	keys := make([]gameTypes.NatureType, 0, len(gameTypes.NatureTypeZhMap))
	for k := range gameTypes.NatureTypeZhMap {
		keys = append(keys, k)
	}

	// 2. 排序 key（NatureType 底层应该是 int）
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// 3. 按照排序后的 key 顺序取值
	for _, k := range keys {
		list = append(list, gameTypes.NatureTypeZhMap[k])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetEnchantEquipPosList(w http.ResponseWriter, r *http.Request) {

	// 1. 先取得所有 key
	list := make([]string, 0, len(autoEnchant.EnchantEquipPosMap))
	for k := range autoEnchant.EnchantEquipPosMap {
		list = append(list, k)
	}

	sort.Strings(list)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetEnchantTypeList(w http.ResponseWriter, r *http.Request) {

	// 1. 先取得所有 key
	list := make([]string, 0, len(autoEnchant.EnchantTypeMap))
	for k := range autoEnchant.EnchantTypeMap {
		list = append(list, k)
	}

	sort.Strings(list)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetExtrasList(w http.ResponseWriter, r *http.Request) {

	// 1. 先取得所有 key
	list := make([]string, 0, len(autoEnchant.ExtraZhMap))
	for _, k := range autoEnchant.ExtraZhMap {
		list = append(list, k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetTradeActionList(w http.ResponseWriter, r *http.Request) {
	list := []string{"买", "卖"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetTradeZhCategoriesList(w http.ResponseWriter, r *http.Request) {
	list := utils.GetMapKeys(gameTypes.TradeZhCategories)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetLotteryTypeList(w http.ResponseWriter, r *http.Request) {
	l := []string{
		"幻想创造器·宴",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    l,
	})
}

func GetDailyTaskLieFengTypeList(w http.ResponseWriter, r *http.Request) {
	l := []string{
		"裂缝",
		"朱诺",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    l,
	})
}
