package handler

import (
	"audit-system/internal/alert"
	"audit-system/internal/service"

	"github.com/kataras/iris/v12"
)

type CreateAlertRuleRequest struct {
	Type    string           `json:"type"`
	Pattern string           `json:"pattern"`
	Level   alert.AlertLevel `json:"level"`
	Message string           `json:"message"`
}

func ListAlertRules(ctx iris.Context) {
	// 验证权限
	userID := ctx.Values().Get("userID").(uint)
	if !service.AuthService.CanManageAlerts(userID) {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"code":    403,
			"message": "没有告警管理权限",
		})
		return
	}

	rules := service.AlertService.GetAlertRules()
	ctx.JSON(iris.Map{
		"code": 200,
		"data": rules,
	})
}

func CreateAlertRule(ctx iris.Context) {
	var req CreateAlertRuleRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	rule := alert.AlertRule{
		Type:    req.Type,
		Pattern: req.Pattern,
		Level:   req.Level,
		Message: req.Message,
	}

	if err := service.AlertService.AddAlertRule(rule); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "创建告警规则失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "告警规则创建成功",
		"data":    rule,
	})
}
