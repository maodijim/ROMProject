package main

import "time"

type Work int

const (
	Init Work = iota
	TeleportMap
	MOVE_RABBIDSPO
	CHECK_RABBIDS
	MOVE_RABBIDSP1
	HUNT_RABBIDS

	End
)

var (
	WorkState Work = Init
)

func HuntCarlen() {
	for {
		switch WorkState {
		case Init:
			Transition(TeleportMap)
			break
		case TeleportMap:
			if g.Configs.HuntConfig.CarryTeam {
				g.TeamGoToMap(*TargetMonster.Mapid)
			} else {
				g.GoToMap(*TargetMonster.Mapid)
			}

			time.Sleep(time.Millisecond * 1000)
			useFlyWing()
			time.Sleep(time.Millisecond * 3200)
			Useskill()
			time.Sleep(time.Millisecond * 1000)
			Useskill()
			time.Sleep(time.Millisecond * 1000)
			Transition(MOVE_RABBIDSPO)
			break
		case MOVE_RABBIDSPO:
			g.MoveChartWait(g.ParsePos(15218, 27, -63822))
			Transition(CHECK_RABBIDS)
			break
		case CHECK_RABBIDS:
			if g.IsMonsterInRange("疯兔") {
				Transition(MOVE_RABBIDSPO)
			} else {
				Transition(MOVE_RABBIDSPO)
			}
			break
		case HUNT_RABBIDS:
			if g.IsMonsterInRange("疯兔") {
				Transition(MOVE_RABBIDSPO)
			} else {
				Transition(MOVE_RABBIDSPO)
			}
			break
		case MOVE_RABBIDSP1:
			g.MoveChartWait(g.ParsePos(15218, 27, -63822))
			Transition(CHECK_RABBIDS)
			break
		}
	}
}

func Transition(SwitchState Work) {
	if WorkState != SwitchState {
		WorkState = SwitchState
	}
}
