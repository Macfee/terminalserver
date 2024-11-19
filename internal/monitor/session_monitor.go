package monitor

import (
	"sync"
	"time"
)

type SessionEvent struct {
	SessionID string
	Type      string
	Data      interface{}
	Time      time.Time
}

type SessionMonitor struct {
	observers sync.Map // map[string][]chan SessionEvent
	mu        sync.RWMutex
}

func NewSessionMonitor() *SessionMonitor {
	return &SessionMonitor{}
}

func (m *SessionMonitor) Subscribe(sessionID string) chan SessionEvent {
	ch := make(chan SessionEvent, 100)
	m.mu.Lock()
	defer m.mu.Unlock()

	observers, _ := m.observers.LoadOrStore(sessionID, make([]chan SessionEvent, 0))
	observerList := observers.([]chan SessionEvent)
	observerList = append(observerList, ch)
	m.observers.Store(sessionID, observerList)

	return ch
}

func (m *SessionMonitor) Unsubscribe(sessionID string, ch chan SessionEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if observers, ok := m.observers.Load(sessionID); ok {
		observerList := observers.([]chan SessionEvent)
		for i, observer := range observerList {
			if observer == ch {
				observerList = append(observerList[:i], observerList[i+1:]...)
				break
			}
		}
		if len(observerList) == 0 {
			m.observers.Delete(sessionID)
		} else {
			m.observers.Store(sessionID, observerList)
		}
	}
	close(ch)
}

func (m *SessionMonitor) BroadcastEvent(event SessionEvent) {
	if observers, ok := m.observers.Load(event.SessionID); ok {
		for _, ch := range observers.([]chan SessionEvent) {
			select {
			case ch <- event:
			default:
				// 通道已满,丢弃事件
			}
		}
	}
}
