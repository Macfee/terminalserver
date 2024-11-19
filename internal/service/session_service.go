package service

import (
	"audit-system/internal/model"
	"audit-system/pkg/database"
	"time"

	"github.com/google/uuid"
)

type sessionService struct{}

var SessionService = new(sessionService)

func (s *sessionService) CreateSession(userID uint, protocol, targetHost string, targetPort int, username string) (*model.Session, error) {
	session := &model.Session{
		SessionID:  uuid.New().String(),
		Protocol:   protocol,
		TargetHost: targetHost,
		TargetPort: targetPort,
		Username:   username,
		Status:     "active",
		StartTime:  time.Now(),
		UserID:     userID,
	}

	if err := database.DB.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionService) GetSession(sessionID string) (*model.Session, error) {
	var session model.Session
	if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *sessionService) ListSessions(page, pageSize int) ([]model.Session, int64, error) {
	var sessions []model.Session
	var total int64

	db := database.DB.Model(&model.Session{})
	db.Count(&total)

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (s *sessionService) TerminateSession(sessionID string) error {
	return database.DB.Model(&model.Session{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":   "terminated",
			"end_time": time.Now(),
		}).Error
}
