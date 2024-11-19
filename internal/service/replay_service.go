package service

import (
	"audit-system/internal/recorder"
	"fmt"
	"time"
)

type replayService struct{}

var ReplayService = new(replayService)

func (s *replayService) GetSessionRecording(sessionID string) ([]recorder.RecordEntry, error) {
	// 从存储中获取会话记录
	rec, err := recorder.NewSessionRecorder(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session recorder: %v", err)
	}

	return rec.GetRecords(time.Time{}, time.Time{}), nil
}
