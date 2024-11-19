package handler

import (
	"audit-system/internal/service"
	"strconv"

	"github.com/kataras/iris/v12"
)

type CreateSessionRequest struct {
	Protocol   string `json:"protocol"`
	TargetHost string `json:"target_host"`
	TargetPort int    `json:"target_port"`
	Username   string `json:"username"`
}

func CreateSession(ctx iris.Context) {
	var req CreateSessionRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	userID := ctx.Values().Get("userID").(uint)
	session, err := service.SessionService.CreateSession(
		userID,
		req.Protocol,
		req.TargetHost,
		req.TargetPort,
		req.Username,
	)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "创建会话失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": session,
	})
}

func ListSessions(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.URLParamDefault("page_size", "10"))

	sessions, total, err := service.SessionService.ListSessions(page, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "获取会话列表失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": iris.Map{
			"total":    total,
			"sessions": sessions,
		},
	})
}

func GetSession(ctx iris.Context) {
	sessionID := ctx.Params().Get("id")
	session, err := service.SessionService.GetSession(sessionID)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"code":    404,
			"message": "会话不存在",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": session,
	})
}

func TerminateSession(ctx iris.Context) {
	sessionID := ctx.Params().Get("id")
	if err := service.SessionService.TerminateSession(sessionID); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "终止会话失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "会话已终止",
	})
}
