package gameConnection

import (
	notifier "ROMProject/gameConnection/types"
)

func (g *GameConnection) AddNotifier(notifierType notifier.NotifierType) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	g.notifier[notifierType] = make(chan interface{})
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
