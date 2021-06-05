package gameConnection

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	QuestProtoCmdId     = Cmd.Command_value["SCENE_USER_QUEST_PROTOCMD"]
	WantedQuestMaxCount = uint32(3)
	WantedQuestType     = "wanted"
)

func (g *GameConnection) GetWantedQuestCompleteCount() (count uint32) {
	if g.Role.UserVars[Cmd.EVarType_EVARTYPE_QUEST_WANTED] != nil {
		count = g.Role.UserVars[Cmd.EVarType_EVARTYPE_QUEST_WANTED].GetValue()
	}
	return count
}

func (g *GameConnection) GetWantedQuestList(questType Cmd.EQuestList) (wantedQuests []*Cmd.QuestData) {
	_, err := g.GetQuestList(questType, 101)
	if err != nil {
		log.Errorf("failed to get quest list: %v", err)
		return wantedQuests
	}
	for _, quest := range g.Role.QuestList.GetList() {
		for _, step := range quest.GetSteps() {
			if step.GetConfig() != nil && step.GetConfig().GetType() == WantedQuestType {
				wantedQuests = append(wantedQuests, quest)
				break
			}
		}
	}
	return wantedQuests
}

func (g *GameConnection) GetQuestList(questType Cmd.EQuestList, id uint32) (questList *Cmd.QuestList, err error) {
	cmd := &Cmd.QuestList{
		Type: &questType,
	}
	if id != 0 {
		cmd.Id = &id
	}
	g.addNotifier("QUESTPARAM_QUESTLIST")
	g.sendProtoCmd(cmd,
		QuestProtoCmdId,
		Cmd.QuestParam_value["QUESTPARAM_QUESTLIST"],
	)
	res, err := g.waitForResponse("QUESTPARAM_QUESTLIST")
	if res != nil {
		ql := res.(*Cmd.QuestList)
		g.Role.QuestList = ql
		questList = ql
	}
	return questList, err
}

func (g *GameConnection) GetWantedCanAcceptQuestList() []*Cmd.QuestData {
	return g.GetWantedQuestList(Cmd.EQuestList_EQUESTLIST_CANACCEPT)
}

func (g *GameConnection) GetSubmitQuestList() []*Cmd.QuestData {
	ql, _ := g.GetQuestList(Cmd.EQuestList_EQUESTLIST_SUBMIT, 0)
	return ql.GetList()
}

func (g *GameConnection) QuestAction(actionType Cmd.EQuestAction, questId uint32) {
	cmd := &Cmd.QuestAction{
		Action:  &actionType,
		Questid: &questId,
	}
	g.sendProtoCmd(cmd,
		QuestProtoCmdId,
		Cmd.QuestParam_value["QUESTPARAM_QUESTACTION"])
}

func (g *GameConnection) QuickSubmitWantedQuest(questId uint32) {
	g.QuestAction(Cmd.EQuestAction_EQUESTACTION_QUICK_SUBMIT_BOARD, questId)
}

func (g *GameConnection) AutoSubmitWantedQuest() {
	go func() {
		for {
			time.Sleep(15 * time.Second)
			completeCount := g.GetWantedQuestCompleteCount()
			if completeCount < WantedQuestMaxCount {
				qla := g.GetWantedCanAcceptQuestList()
				log.Debugf("%v  complete count: %d", qla, completeCount)
				for _, quest := range qla {
					log.Infof("提交委托任务: %s", quest.GetSteps()[0].GetConfig().GetName())
					g.QuickSubmitWantedQuest(quest.GetId())
					time.Sleep(5 * time.Second)
				}
			} else {
				log.Infof("每日委托任务次数已达到%d次 跳过", completeCount)
				return
			}
		}
	}()
}
