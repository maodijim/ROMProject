package gameConnection

import (
	notifier "ROMProject/gameConnection/types"
)

func (g *GameConnection) AddNotifier(notifierType notifier.NotifierType) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	if val, exists := g.notifier[notifierType]; exists && val != nil {
		return
	}
	g.notifier[notifierType] = make(chan interface{}, 3)
}

func (g *GameConnection) RemoveNotifier(notifierType notifier.NotifierType) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	g.notifier[notifierType] = nil
}

func (g *GameConnection) Notifier(notifierType notifier.NotifierType) chan interface{} {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	return g.notifier[notifierType]
}
