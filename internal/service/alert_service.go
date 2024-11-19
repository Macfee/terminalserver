package service

import (
	"audit-system/internal/alert"
	"fmt"
)

type alertService struct {
	manager *alert.AlertManager
}

var AlertService = new(alertService)

func init() {
	AlertService.manager = alert.NewAlertManager()
}

func (s *alertService) GetAlertRules() []alert.AlertRule {
	return s.manager.GetRules()
}

func (s *alertService) AddAlertRule(rule alert.AlertRule) error {
	return s.manager.AddRule(rule)
}

func (s *alertService) UpdateRule(ruleID string, rule alert.AlertRule) error {
	// 从现有规则中查找要更新的规则
	rules := s.manager.GetRules()
	var oldRule alert.AlertRule
	found := false

	for _, r := range rules {
		if r.Type == ruleID { // 使用Type作为规则的唯一标识
			oldRule = r
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("规则不存在: %s", ruleID)
	}

	// 使用 AlertManager 的 UpdateRule 方法更新规则
	return s.manager.UpdateRule(oldRule.Type, oldRule.Pattern, rule)
}

func (s *alertService) DeleteRule(ruleID string) error {
	// 从现有规则中查找要删除的规则
	rules := s.manager.GetRules()
	var targetRule alert.AlertRule
	found := false

	for _, r := range rules {
		if r.Type == ruleID { // 使用Type作为规则的唯一标识
			targetRule = r
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("规则不存在: %s", ruleID)
	}

	// 使用 AlertManager 的 RemoveRule 方法删除规则
	return s.manager.RemoveRule(targetRule.Type, targetRule.Pattern)
}
