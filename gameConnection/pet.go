package gameConnection

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	PetCookWorkSpaceId = uint32(1)
)

var (
	PetProtoCmdId = Cmd.Command_value["SCENE_USER_PET_PROTOCMD"]
)

// QueryBattlePet 获取出战宠物
func (g *GameConnection) QueryBattlePet() (battlePet *Cmd.QueryBattlePetCmd) {
	cmd := &Cmd.QueryBattlePetCmd{}
	g.addNotifier("PETPARAM_ADVENTURE_QUERYBATTLEPET")
	g.sendProtoCmd(cmd, PetProtoCmdId, Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYBATTLEPET"])
	res, err := g.waitForResponse("PETPARAM_ADVENTURE_QUERYBATTLEPET")
	if err != nil {
		log.Errorf("failed to get battle pet: %v", err)
		return nil
	}
	if res != nil {
		battlePet = res.(*Cmd.QueryBattlePetCmd)
	}
	return battlePet
}

// QueryPetWorkData 获取宠物打工数据
func (g *GameConnection) QueryPetWorkData() (workData *Cmd.QueryPetWorkDataPetCmd, err error) {
	cmd := &Cmd.QueryPetWorkDataPetCmd{}
	g.addNotifier("PETPARAM_WORK_QUERYWORKDATA")
	g.sendProtoCmd(cmd,
		PetProtoCmdId,
		Cmd.PetParam_value["PETPARAM_WORK_QUERYWORKDATA"],
	)
	res, err := g.waitForResponse("PETPARAM_WORK_QUERYWORKDATA")
	if res != nil {
		workData = res.(*Cmd.QueryPetWorkDataPetCmd)
	}
	return workData, err
}

func (g *GameConnection) GetPetCookWorkSpaceInfo() (cookWorkSpace *Cmd.WorkSpace) {
	workData, err := g.QueryPetWorkData()
	if err != nil {
		log.Errorf("failed to get pet work data: %v", err)
	}
	if workData != nil {
		for _, wp := range workData.GetDatas() {
			if wp.GetId() == PetCookWorkSpaceId {
				cookWorkSpace = wp
				break
			}
		}
	}
	return cookWorkSpace
}

func (g *GameConnection) GetPetWorkReward(workSpaceId uint32) {
	cmd := &Cmd.GetPetWorkRewardPetCmd{
		Id: &workSpaceId,
	}
	g.sendProtoCmd(cmd,
		PetProtoCmdId,
		Cmd.PetParam_value["PETPARAM_WORK_GETREWARD"],
	)
}

func (g *GameConnection) GetCookWorkSpaceReward() {
	if g.IsCookWorkSpaceInUse() {
		workSpace := g.GetPetCookWorkSpaceInfo()
		startTime := time.Unix(int64(workSpace.GetStarttime()), 0)
		lastRewardTime := time.Unix(int64(workSpace.GetLastrewardtime()), 0)
		minWorkTime := 30 * time.Minute
		if time.Since(startTime) > minWorkTime && time.Since(lastRewardTime) > minWorkTime {
			log.Infof("%s: get pet cook workspace reward", g.Role.GetRoleName())
			g.GetPetWorkReward(PetCookWorkSpaceId)
		} else {
			log.Infof("%s: workspace start time %v last reward time %v has not been %v yet skip get work space reward",
				g.Role.GetRoleName(), startTime, lastRewardTime, minWorkTime,
			)
		}
	}
}

func (g *GameConnection) PetRestoreEgg(petId uint32) {
	cmd := &Cmd.EggRestorePetCmd{
		Petid: &petId,
	}
	g.sendProtoCmd(cmd,
		PetProtoCmdId,
		Cmd.PetParam_value["PETPARAM_RESTORE_EGG"],
	)
}

func (g *GameConnection) PetStartWork(workSpaceId uint32, petGuid []string) {
	cmd := &Cmd.StartWorkPetCmd{
		Id:   &workSpaceId,
		Pets: petGuid,
	}
	g.sendProtoCmd(cmd,
		PetProtoCmdId,
		Cmd.PetParam_value["PETPARAM_WORK_STARTWORK"],
	)
}

func (g *GameConnection) IsCookWorkSpaceInUse() bool {
	workSpace := g.GetPetCookWorkSpaceInfo()
	if workSpace.GetState() == Cmd.EWorkState_EWORKSTATE_WORKING {
		return true
	}
	return false
}

func (g *GameConnection) BattlePetToWork() {
	if !g.IsCookWorkSpaceInUse() {
		battlePet := g.QueryBattlePet()
		if len(battlePet.GetPets()) == 0 {
			log.Warnf("%s: 没有找到战斗宠物", g.Role.GetRoleName())
			return
		}
		egg := battlePet.GetPets()[0].GetEgg()
		petEggGuid := egg.GetGuid()
		petEggId := egg.GetId()
		petLv := egg.GetLv()
		if petLv < 70 {
			log.Warnf("%s: pet lv %d is lower than require lv 70 to work", g.Role.GetRoleName(), petLv)
			return
		}
		g.PetRestoreEgg(petEggId)
		time.Sleep(5 * time.Second)
		g.PetStartWork(PetCookWorkSpaceId, []string{petEggGuid})
	}
}
