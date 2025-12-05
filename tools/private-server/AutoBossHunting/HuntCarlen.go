package AutoBossHunting

import (
	"time"

	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

func i32(v int32) *int32 { return &v }

var CrazyRabbitPos = []Cmd.ScenePos{
	{X: i32(10165), Y: i32(27), Z: i32(-57086)},
	{X: i32(20845), Y: i32(27), Z: i32(-85483)},
	{X: i32(-29854), Y: i32(27), Z: i32(-84827)},
	{X: i32(-41499), Y: i32(127), Z: i32(-51326)},
}

var CarlenPos = []Cmd.ScenePos{
	{X: i32(-29854), Y: i32(62), Z: i32(-87337)},
}
var ScreamingDemonPos = []Cmd.ScenePos{
	{X: i32(-2277), Y: i32(11032), Z: i32(-48103)},
	{X: i32(-30273), Y: i32(11032), Z: i32(-63487)},
	{X: i32(-54962), Y: i32(11032), Z: i32(-62435)},
}

var BigBadWolfPos = []Cmd.ScenePos{
	{X: i32(5907), Y: i32(11032), Z: i32(-82878)},
	{X: i32(65253), Y: i32(8432), Z: i32(-24051)},
	{X: i32(89144), Y: i32(8432), Z: i32(7776)},
	{X: i32(-9598), Y: i32(8499), Z: i32(10583)},
}

func (b *BossHuntTask) huntCarlen() {
	TargetMVP := b.targetHiddenMVP.Info

	PrerequisiteMonsters := b.targetHiddenMVP.PrerequisiteMonsters

	PosList := b.targetHiddenMVP.PrerequisitePosList

	BossPosList := b.targetHiddenMVP.BossPosList

	for {
		b.checkApear(TargetMVP.NameZh) // 确认卡伦是否复活
		switch b.workState {
		// 初始化
		case Init:
			log.Infof("开始狩猎%s", TargetMVP.NameZh)
			b.findBoss = false
			b.haveBoss = false
			b.posCount = 0
			b.BossposCount = 0
			b.transition(TeleportMap)
			break
		// 传送到目标地图
		case TeleportMap:
			b.logger.Infof("传送到%s", gameTypes.MapIdToZh[b.targetHiddenMVP.Map])
			b.GC.InMap(b.targetHiddenMVP.Map.Uint32())
			b.transition(MOVE_PrerequisiteMonstersPOS)
			break
		// 移动到前置怪物地点
		case MOVE_PrerequisiteMonstersPOS:
			if int(b.posCount) < len(PosList)-1 {
				if b.GC.MoveChartWait(PosList[b.posCount]) {
					log.Infof("抵达%s出生点%d", PrerequisiteMonsters.NameZh, b.posCount+1)
					b.transition(CHECK_PrerequisiteMonsters)
				}
			} else {
				log.Infof("%s狩猎完成，开始狩猎%s", PrerequisiteMonsters.NameZh, TargetMVP.NameZh)
				b.transition(MOVE_BOSSPOS)
			}
			break
		// 确认有前置怪物
		case CHECK_PrerequisiteMonsters:
			if b.GC.IsMonsterInRange(PrerequisiteMonsters.NameZh) {
				b.logger.Infof("出生点%d发现%s，开始狩猎%s", b.posCount+1, PrerequisiteMonsters.NameZh, PrerequisiteMonsters.NameZh)
				b.haveBoss = true
				b.fightMonstStar(PrerequisiteMonsters.NameZh, 0)
				b.transition(HUNT_PrerequisiteMonsters)
			} else {
				b.logger.Infof("出生点%d未发现%s，到%s点%d", b.posCount+1, PrerequisiteMonsters.NameZh, PrerequisiteMonsters.NameZh, b.posCount+2)
				b.fightCancel()
				b.posCount++
				b.transition(MOVE_PrerequisiteMonstersPOS)
			}
			break
		// 狩猎疯兔
		case HUNT_PrerequisiteMonsters:
			if !b.GC.IsMonsterInRange(PrerequisiteMonsters.NameZh) {
				b.logger.Infof("%s点%d狩猎完成，到%s点%d", PrerequisiteMonsters.NameZh, b.posCount+1, PrerequisiteMonsters.NameZh, b.posCount+2)
				b.posCount++
				b.fightCancel()
				b.transition(MOVE_PrerequisiteMonstersPOS)
			}
			break
		// 移动到卡伦出生位置
		case MOVE_BOSSPOS:
			if int(b.BossposCount) < len(BossPosList) {
				b.GC.MoveChartWait(BossPosList[b.BossposCount])
				log.Infof("抵达%s出生点%d", TargetMVP.NameZh, b.BossposCount+1)
				b.transition(CHECK_BOSS)
			} else {
				b.transition(End)
			}
			break
		// 确认卡伦
		case CHECK_BOSS:
			if b.GC.IsMonsterInRange(TargetMVP.NameZh) {
				b.fightMonstStar(TargetMVP.NameZh, 0)
				b.transition(HUNT_BOSS)
			} else {
				b.fightCancel()
				b.fightStar = false
				b.BossposCount++
				b.transition(MOVE_BOSSPOS)
			}
			break
		// 狩猎卡伦
		case HUNT_BOSS:
			if b.GC.IsMonsterInRange(TargetMVP.NameZh) {
				b.haveBoss = true
				TargetID := b.GC.AtkStat.GetCurrentTargetId()
				MapNPC := b.GC.GetMapNpcs()
				if TargetID != 0 && MapNPC[TargetID].Attrs != nil {
					MonsterHP := utils.GetNpcAttrValByType(MapNPC[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
					MHP := utils.GetNpcAttrValByType(b.GC.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_HP)
					if MonsterHP != b.tempMHP || MHP != b.tempUHP {
						b.tempMHP = MonsterHP
						b.tempUHP = MHP
						b.logger.Infof("%s 未死亡，剩余血量:%d", TargetMVP.NameZh, MonsterHP)
						b.logger.Infof("我的血量:%d", MHP)
					} else if MonsterHP == 0 {
						b.transition(End)
					}
				}
			} else {
				b.transition(End)
			}
			break
		// 结束
		case End:
			if b.haveBoss {
				b.targetHiddenMVP.RespawnTime = time.Now().Add(30 * time.Minute)
			} else {
				b.targetHiddenMVP.RespawnTime = time.Now().Add(10 * time.Minute)
			}
			b.fightCancel()
			b.fightStar = false
			b.logger.Infof("%s已死亡，复活时间:%s", TargetMVP.NameZh, b.targetHiddenMVP.RespawnTime.Format("2006-01-02 15:04:05"))
			time.Sleep(time.Millisecond * 3000)
			b.transition(Init)
			return
		}

		time.Sleep(time.Millisecond * 100)
	}
}
func (b *BossHuntTask) checkApear(BossName string) bool {
	if !b.findBoss && b.GC.IsMonsterInRange(BossName) {
		b.findBoss = true
		b.logger.Infof("发现%s，开始狩猎", BossName)
		b.fightCancel()
		b.fightMonstStar(BossName, 0)
		b.transition(HUNT_BOSS)
		return true
	} else {
		return false
	}
}
