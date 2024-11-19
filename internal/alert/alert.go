package alert

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"sync"
	"time"
)

type AlertLevel int

const (
	LevelInfo AlertLevel = iota
	LevelWarning
	LevelDanger
)

type Alert struct {
	ID        string
	Level     AlertLevel
	Type      string
	Message   string
	SessionID string
	UserID    uint
	Time      time.Time
	Handled   bool
}

type AlertRule struct {
	Type    string
	Pattern string
	Level   AlertLevel
	Message string
}

type AlertManager struct {
	rules    []AlertRule
	alerts   sync.Map
	handlers []AlertHandler
}

type AlertHandler interface {
	HandleAlert(alert Alert) error
}

func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules: []AlertRule{
			{
				Type:    "command",
				Pattern: "rm -rf /*",
				Level:   LevelDanger,
				Message: "危险命令执行告警",
			},
			{
				Type:    "login",
				Pattern: "failed_attempts > 3",
				Level:   LevelWarning,
				Message: "登录失败次数过多",
			},
		},
	}
}

func (m *AlertManager) AddHandler(handler AlertHandler) {
	m.handlers = append(m.handlers, handler)
}

func (m *AlertManager) CheckCommand(sessionID string, userID uint, command string) error {
	for _, rule := range m.rules {
		if rule.Type == "command" && matchPattern(command, rule.Pattern) {
			alert := Alert{
				ID:        generateAlertID(),
				Level:     rule.Level,
				Type:      rule.Type,
				Message:   rule.Message,
				SessionID: sessionID,
				UserID:    userID,
				Time:      time.Now(),
			}

			m.alerts.Store(alert.ID, alert)
			return m.notifyHandlers(alert)
		}
	}
	return nil
}

func generateAlertID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func matchPattern(input, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, input)
	return matched
}

func (m *AlertManager) notifyHandlers(alert Alert) error {
	for _, handler := range m.handlers {
		if err := handler.HandleAlert(alert); err != nil {
			return err
		}
	}
	return nil
}

func (m *AlertManager) GetRules() []AlertRule {
	return m.rules
}

func (m *AlertManager) AddRule(rule AlertRule) error {
	if rule.Type == "" || rule.Pattern == "" {
		return fmt.Errorf("规则类型和匹配模式不能为空")
	}

	for _, existingRule := range m.rules {
		if existingRule.Type == rule.Type && existingRule.Pattern == rule.Pattern {
			return fmt.Errorf("规则已存在")
		}
	}

	m.rules = append(m.rules, rule)
	return nil
}

func (m *AlertManager) RemoveRule(ruleType, pattern string) error {
	for i, rule := range m.rules {
		if rule.Type == ruleType && rule.Pattern == pattern {
			m.rules = append(m.rules[:i], m.rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("规则不存在")
}

func (m *AlertManager) UpdateRule(oldType, oldPattern string, newRule AlertRule) error {
	if err := m.RemoveRule(oldType, oldPattern); err != nil {
		return err
	}
	return m.AddRule(newRule)
}
