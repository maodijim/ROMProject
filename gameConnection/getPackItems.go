package gameConnection

import (
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

const (
	packRetry = "packRetry"
)

func (g *GameConnection) GetMainPackItems() (err error) {
	packTypes := []Cmd.EPackType{
		Cmd.EPackType_EPACKTYPE_MAIN,
	}
	for _, pType := range packTypes {
		err = g.GetPackItem(&pType)
		time.Sleep(time.Second)
	}
	return err
}

func (g *GameConnection) GetEquipPackItems() (err error) {
	packType := Cmd.EPackType_EPACKTYPE_EQUIP
	err = g.GetPackItem(&packType)
	return err
}

func (g *GameConnection) GetTempMainPackItems() (err error) {
	packTypes := []Cmd.EPackType{
		Cmd.EPackType_EPACKTYPE_TEMP_MAIN,
	}
	for _, pType := range packTypes {
		err = g.GetPackItem(&pType)
		time.Sleep(time.Second)
	}
	return err
}

func (g *GameConnection) GetAllPackItems() (err error) {
	packTypes := []Cmd.EPackType{
		Cmd.EPackType_EPACKTYPE_MAIN,
		Cmd.EPackType_EPACKTYPE_EQUIP,
		// Cmd.EPackType_EPACKTYPE_STORE,
		Cmd.EPackType_EPACKTYPE_PERSONAL_STORE,
		Cmd.EPackType_EPACKTYPE_FOOD,
		Cmd.EPackType_EPACKTYPE_PET,
		Cmd.EPackType_EPACKTYPE_QUEST,
	}
	wg := sync.WaitGroup{}
	for range packTypes {
		wg.Add(1)
	}
	for i, pType := range packTypes {
		go func(packType Cmd.EPackType, i int) {
			time.Sleep(time.Duration(i) * 250 * time.Millisecond)
			err = g.GetPackItem(&packType)
			wg.Done()
		}(pType, i)
	}
	wg.Wait()
	return err
}

func (g *GameConnection) GetPackItem(packType *Cmd.EPackType) (err error) {
	cmd := &Cmd.PackageItem{
		Type: packType,
	}
	_ = g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"],
		Cmd.ItemParam_value["ITEMPARAM_PACKAGEITEM"],
	)
	err = g.WaitForGetPackItems(packType)
	if err != nil {
		log.Errorf("get pack item error: %s", err)
		// Retries
		if err == ErrQueryTimeout {
			log.Warnf("get pack item %v timed out retrying %d", packType, g.retries[packRetry])
			if g.retries[packRetry] < maxRetry {
				g.retries[packRetry] += 1
				err = g.GetPackItem(packType)
				if err == nil {
					g.retries[packRetry] = 0
					return err
				}
			}
		} else {
			g.retries[tradeHis] = 0
		}
	}
	return err
}

func (g *GameConnection) WaitForGetPackItems(pType *Cmd.EPackType) (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("Checking for pack items response")
		g.Mutex.Lock()
		if g.Role.GetPackItems()[*pType] == nil {
			g.Mutex.Unlock()
			time.Sleep(2 * time.Second)
			continue
		} else {
			g.Mutex.Unlock()
			break
		}
	}
	if time.Since(startQueryTime) > queryTimeout {
		err = ErrQueryTimeout
	}
	return err
}
