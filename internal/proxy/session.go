package proxy

import (
	"audit-system/internal/recorder"
	"errors"
	"sync"
)

type Session struct {
	ID       string
	Proxy    Proxy
	Recorder *recorder.SessionRecorder
}

type SessionManager struct {
	sessions sync.Map
	factory  ProxyFactory
}

func NewSessionManager(factory ProxyFactory) *SessionManager {
	return &SessionManager{
		factory: factory,
	}
}

func (m *SessionManager) CreateSession(sessionID, protocol, host string, port int, username, password string) error {
	proxy, err := m.factory.CreateProxy(protocol, host, port, username, password)
	if err != nil {
		return err
	}

	recorder, err := recorder.NewSessionRecorder(sessionID)
	if err != nil {
		return err
	}

	session := &Session{
		ID:       sessionID,
		Proxy:    proxy,
		Recorder: recorder,
	}

	m.sessions.Store(sessionID, session)
	return proxy.Connect()
}

func (m *SessionManager) GetSession(sessionID string) (*Session, error) {
	if session, ok := m.sessions.Load(sessionID); ok {
		return session.(*Session), nil
	}
	return nil, errors.New("session not found")
}

func (m *SessionManager) CloseSession(sessionID string) error {
	if session, ok := m.sessions.Load(sessionID); ok {
		s := session.(*Session)
		if err := s.Proxy.Close(); err != nil {
			return err
		}
		m.sessions.Delete(sessionID)
		return nil
	}
	return errors.New("session not found")
}
