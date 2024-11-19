package recorder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type CommandRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Command   string    `json:"command"`
	Output    string    `json:"output"`
	ExitCode  int       `json:"exit_code"`
}

type CommandRecorder struct {
	sessionID string
	file      *os.File
}

func NewCommandRecorder(sessionID string) (*CommandRecorder, error) {
	recordPath := filepath.Join("records", "commands", sessionID)
	if err := os.MkdirAll(recordPath, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(filepath.Join(recordPath, "commands.log"))
	if err != nil {
		return nil, err
	}

	return &CommandRecorder{
		sessionID: sessionID,
		file:      file,
	}, nil
}

func (r *CommandRecorder) RecordCommand(command string, output string, exitCode int) error {
	record := CommandRecord{
		Timestamp: time.Now(),
		Command:   command,
		Output:    output,
		ExitCode:  exitCode,
	}

	encoded, err := json.Marshal(record)
	if err != nil {
		return err
	}

	if _, err := r.file.Write(append(encoded, '\n')); err != nil {
		return err
	}

	return nil
}

func (r *CommandRecorder) Close() error {
	return r.file.Close()
}
