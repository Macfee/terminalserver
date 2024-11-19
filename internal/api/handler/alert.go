package handler

import (
	"audit-system/internal/alert"
	"audit-system/internal/service"

	"github.com/kataras/iris/v12"
)

func UpdateAlertRule(ctx iris.Context) {
	ruleID := ctx.Params().Get("id")
	var rule alert.AlertRule
	if err := ctx.ReadJSON(&rule); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	if err := service.AlertService.UpdateRule(ruleID, rule); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "更新告警规则失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "更新成功",
	})
}

func DeleteAlertRule(ctx iris.Context) {
	ruleID := ctx.Params().Get("id")
	if err := service.AlertService.DeleteRule(ruleID); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "删除告警规则失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "删除成功",
	})
}
