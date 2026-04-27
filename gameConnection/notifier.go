package gameConnection

import (
	"time"

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

func (g *GameConnection) SendToNotifier(notifier notifier.NotifierType, msg interface{}) {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	timeout := time.After(10 * time.Second)
	if ch := g.notifier[notifier]; ch != nil {
		go func() {
			select {
			case ch <- msg:
			case <-timeout:
				g.logger.Warnf("send to notifier %s timeout", notifier)
			}
		}()
	}
}
