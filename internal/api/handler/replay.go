package handler

import (
	"audit-system/internal/recorder"
	"audit-system/internal/service"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
)

func GetSessionReplay(ctx iris.Context) {
	sessionID := ctx.Params().Get("id")
	startTimeStr := ctx.URLParam("start_time")
	endTimeStr := ctx.URLParam("end_time")
	speed, _ := strconv.ParseFloat(ctx.URLParamDefault("speed", "1.0"), 64)

	// 验证权限
	userID := ctx.Values().Get("userID").(uint)
	if !service.AuthService.CanReplaySession(userID, sessionID) {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"code":    403,
			"message": "没有回放权限",
		})
		return
	}

	var startTime, endTime time.Time
	if startTimeStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startTimeStr)
	}
	if endTimeStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endTimeStr)
	}

	// 获取回放数据
	entries, err := service.ReplayService.GetSessionRecording(sessionID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "获取会话回放数据失败",
		})
		return
	}

	// 过滤指定时间范围的记录
	var filteredEntries []recorder.RecordEntry
	for _, entry := range entries {
		if (startTime.IsZero() || entry.Timestamp.After(startTime)) &&
			(endTime.IsZero() || entry.Timestamp.Before(endTime)) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": iris.Map{
			"session_id": sessionID,
			"speed":      speed,
			"entries":    filteredEntries,
		},
	})
}
