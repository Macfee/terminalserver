package service

import (
	"audit-system/internal/model"
	"audit-system/pkg/database"
	"time"
)

type auditService struct{}

var AuditService = new(auditService)

func (s *auditService) LogCommand(sessionID string, command string) error {
	log := &model.AuditLog{
		SessionID: sessionID,
		Type:      "command",
		Content:   command,
		Timestamp: time.Now(),
	}
	return database.DB.Create(log).Error
}

func (s *auditService) LogFileTransfer(sessionID string, filename string, direction string) error {
	log := &model.AuditLog{
		SessionID: sessionID,
		Type:      "file_transfer",
		Content:   filename + ":" + direction,
		Timestamp: time.Now(),
	}
	return database.DB.Create(log).Error
}

func (s *auditService) GetAuditLogs(sessionID string, logType string, startTime, endTime time.Time, page, pageSize int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	db := database.DB.Model(&model.AuditLog{}).Where("session_id = ?", sessionID)

	if logType != "" {
		db = db.Where("type = ?", logType)
	}

	if !startTime.IsZero() {
		db = db.Where("timestamp >= ?", startTime)
	}

	if !endTime.IsZero() {
		db = db.Where("timestamp <= ?", endTime)
	}

	db.Count(&total)

	if err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("timestamp desc").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
