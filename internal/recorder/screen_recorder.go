package recorder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ScreenRecorder struct {
	sessionID  string
	file       *os.File
	startTime  time.Time
	resolution struct {
		width  int
		height int
	}
}

type ScreenFrame struct {
	Timestamp  time.Time `json:"timestamp"`
	ImageData  []byte    `json:"image_data"`
	Resolution struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"resolution"`
}

func NewScreenRecorder(sessionID string, width, height int) (*ScreenRecorder, error) {
	recorder := &ScreenRecorder{
		sessionID:  sessionID,
		startTime:  time.Now(),
		resolution: struct{ width, height int }{width, height},
	}

	recordPath := filepath.Join("records", "screen", sessionID)
	if err := os.MkdirAll(recordPath, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(filepath.Join(recordPath, "screen.rec"))
	if err != nil {
		return nil, err
	}
	recorder.file = file

	return recorder, nil
}

func (r *ScreenRecorder) RecordFrame(imageData []byte) error {
	frame := ScreenFrame{
		Timestamp: time.Now(),
		ImageData: imageData,
		Resolution: struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		}{
			Width:  r.resolution.width,
			Height: r.resolution.height,
		},
	}

	encoded, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	if _, err := r.file.Write(append(encoded, '\n')); err != nil {
		return err
	}

	return nil
}

func (r *ScreenRecorder) Close() error {
	return r.file.Close()
}
