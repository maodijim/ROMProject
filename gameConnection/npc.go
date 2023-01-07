package gameConnection

import (
	"errors"
	"fmt"
	"time"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

var (
	sceneUserIntertCmdId = Cmd.Command_value["SCENE_USER_INTER_PROTOCMD"]
)

func (g *GameConnection) VisitNpc(npcId uint64) {
	cmdMap := Cmd.MapObjectData{
		Guid: &npcId,
	}
	_ = g.sendProtoCmd(
		&cmdMap,
		sceneUserQuestId,
		Cmd.CmdParam_value["MAP_OBJECT_DATA"],
	)
	time.Sleep(500 * time.Millisecond)
	cmd := Cmd.VisitNpcUserCmd{
		Npctempid: &npcId,
	}
	_ = g.sendProtoCmd(
		&cmd,
		sceneUserQuestId,
		Cmd.QuestParam_value["QUESTPARAM_VISIT_NPC"],
	)
}

func (g *GameConnection) WaitForInterQuestion(interId uint32) (inter *Cmd.NewInter, err error) {
	var iq *Cmd.NewInter
	for {
		select {
		case <-time.After(3 * time.Second):
			if iq != nil {
				go func() {
					g.notifier["INTER_QUESTION"] <- iq
				}()
			}
			return nil, errors.New(fmt.Sprintf("wait for inter question %d timeout", interId))
		case note := <-g.notifier["INTER_QUESTION"]:
			iq = note.(*Cmd.NewInter)
			if iq.GetInter().GetInterid() == interId {
				return iq, nil
			} else {
				log.Warnf("inter id not match: got %v want %d", iq.GetInter().GetInterid(), interId)
			}
		}
	}
}

func (g *GameConnection) Answer(npcId uint64, interId, Answer uint32) {
	iq, err := g.WaitForInterQuestion(interId)
	if err != nil {
		log.Errorf("failed to wait for inter question: %v", err)
		return
	}
	log.Infof("Answering inter question: %v", iq)
	guid := iq.GetInter().GetGuid()
	cmd := Cmd.Answer{
		Npcid:   &npcId,
		Interid: &interId,
		Answer:  &Answer,
		Guid:    &guid,
	}
	_ = g.sendProtoCmd(
		&cmd,
		sceneUserIntertCmdId,
		Cmd.InterParam_value["INTERPARAM_ANSWERINTER"],
	)
}
