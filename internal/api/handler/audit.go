package handler

import (
	"audit-system/internal/service"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
)

func GetAuditLogs(ctx iris.Context) {
	sessionID := ctx.URLParam("session_id")
	logType := ctx.URLParam("type")
	startTimeStr := ctx.URLParam("start_time")
	endTimeStr := ctx.URLParam("end_time")
	page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.URLParamDefault("page_size", "10"))

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startTimeStr)
	}
	if endTimeStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endTimeStr)
	}

	logs, total, err := service.AuditService.GetAuditLogs(
		sessionID,
		logType,
		startTime,
		endTime,
		page,
		pageSize,
	)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "获取审计日志失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": iris.Map{
			"total": total,
			"logs":  logs,
		},
	})
}
