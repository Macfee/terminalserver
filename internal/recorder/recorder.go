package recorder

import (
	"time"
)

// RecordEntry 表示一条会话记录
type RecordEntry struct {
	Timestamp time.Time
	Type      string // 记录类型：input/output
	Data      []byte // 记录内容
	SessionID string // 会话ID
}

// SessionRecorder 负责记录会话内容
type SessionRecorder struct {
	sessionID string
	records   []RecordEntry
}

// NewSessionRecorder 创建新的会话记录器
func NewSessionRecorder(sessionID string) (*SessionRecorder, error) {
	return &SessionRecorder{
		sessionID: sessionID,
		records:   make([]RecordEntry, 0),
	}, nil
}

// Record 记录一条会话数据
func (r *SessionRecorder) Record(recordType string, data []byte) error {
	entry := RecordEntry{
		Timestamp: time.Now(),
		Type:      recordType,
		Data:      data,
		SessionID: r.sessionID,
	}
	r.records = append(r.records, entry)
	return nil
}

// GetRecords 获取指定时间范围内的记录
func (r *SessionRecorder) GetRecords(start, end time.Time) []RecordEntry {
	var filtered []RecordEntry
	for _, record := range r.records {
		if (start.IsZero() || record.Timestamp.After(start)) &&
			(end.IsZero() || record.Timestamp.Before(end)) {
			filtered = append(filtered, record)
		}
	}
	return filtered
}
