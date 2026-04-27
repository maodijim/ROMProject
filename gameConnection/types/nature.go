package gameConnection

type NatureType string

const (
	NatureType_Fire     NatureType = "Fire"
	NatureType_Water    NatureType = "Water"
	NatureType_Earth    NatureType = "Earth"
	NatureType_Wind     NatureType = "Wind"
	NatureType_Ghost    NatureType = "Ghost"
	NatureType_Neutral  NatureType = "Neutral"
	NatureType_Holy     NatureType = "Holy"
	NatureType_Undead   NatureType = "Undead"
	NatureType_Shadow   NatureType = "Shadow"
	NatureType_Poison   NatureType = "Poison"
	NatureType_Formless NatureType = "Formless"
)

var NatureTypeZhMap = map[NatureType]string{
	NatureType_Fire:     "火属性",
	NatureType_Water:    "水属性",
	NatureType_Earth:    "地属性",
	NatureType_Wind:     "风属性",
	NatureType_Ghost:    "念属性",
	NatureType_Neutral:  "无属性",
	NatureType_Holy:     "圣属性",
	NatureType_Undead:   "不死属性",
	NatureType_Shadow:   "暗属性",
	NatureType_Poison:   "毒属性",
	NatureType_Formless: "无形属性",
}

var NatureTypeFromZhMap = func() map[string]NatureType {
	m := make(map[string]NatureType)
	for k, v := range NatureTypeZhMap {
		m[v] = k
	}
	return m
}()

func GetNatureTypeFromZhFast(zh string) NatureType {
	t, ok := NatureTypeFromZhMap[zh]
	if ok {
		return t
	} else {
		return ""
	}
}
