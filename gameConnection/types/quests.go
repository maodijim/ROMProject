package gameConnection

type SealQuestType uint32

const (
	SealQuestType_WestGate           SealQuestType = 101 // 西门裂缝
	SealQuestType_GlastHeimOutskirts SealQuestType = 112 // 古城郊外裂缝
	SealQuestType_MjolnirMountains   SealQuestType = 108 // 妙勒尼山脉裂缝
)

func (sealType SealQuestType) String() string {
	switch sealType {
	case SealQuestType_WestGate:
		return "西门裂缝"
	case SealQuestType_GlastHeimOutskirts:
		return "古城郊外裂缝"
	default:
		return "未知裂缝"
	}
}
