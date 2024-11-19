package playback

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type PlaybackType string

const (
    PlaybackCommand PlaybackType = "command"
    PlaybackScreen  PlaybackType = "screen"
)

type SessionPlayer struct {
    sessionID string
    startTime time.Time
    endTime   time.Time
    speed     float64 // 回放速度倍率
}

func NewSessionPlayer(sessionID string) *SessionPlayer {
    return &SessionPlayer{
        sessionID: sessionID,
        speed:     1.0,
    }
}

func (p *SessionPlayer) SetTimeRange(start, end time.Time) {
    p.startTime = start
    p.endTime = end
}

func (p *SessionPlayer) SetSpeed(speed float64) {
    p.speed = speed
}

func (p *SessionPlayer) PlayCommands(ch chan<- CommandRecord) error {
    recordPath := filepath.Join("records", "commands", p.sessionID, "commands.log")
    file, err := os.Open(recordPath)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    for decoder.More() {
        var record CommandRecord
        if err := decoder.Decode(&record); err != nil {
            return err
        }

        // 检查时间范围
        if !p.startTime.IsZero() && record.Timestamp.Before(p.startTime) {
            continue
        }
        if !p.endTime.IsZero() && record.Timestamp.After(p.endTime) {
            break
        }

        ch <- record
        
        // 根据回放速度控制发送间隔
        if decoder.More() {
            var
        }
    }
    return nil
} 