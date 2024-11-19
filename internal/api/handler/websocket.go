package handler

import (
	"audit-system/internal/proxy"
	"audit-system/internal/service"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
)

var sessionManager *proxy.SessionManager
var wsServer *neffos.Server

func init() {
	factory := proxy.NewProxyFactory()
	sessionManager = proxy.NewSessionManager(*factory)

	wsServer = websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: handleMessage,
		"disconnect":              handleDisconnect,
	})
}

func HandleWebSocket(ctx iris.Context) {
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

	// 创建代理会话
	err = sessionManager.CreateSession(
		sessionID,
		session.Protocol,
		session.TargetHost,
		session.TargetPort,
		session.Username,
		"", // 密码需要从安全存储中获取
	)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "创建代理会话失败",
		})
		return
	}

	// 升级到WebSocket连接
	websocket.Upgrade(ctx, websocket.DefaultIDGenerator, wsServer)
}

func handleMessage(ns *neffos.NSConn, msg neffos.Message) error {
	sessionID := ns.Conn.ID()
	proxySession, err := sessionManager.GetSession(sessionID)
	if err != nil {
		return err
	}

	// 记录审计日志
	service.AuditService.LogCommand(sessionID, string(msg.Body))

	// 写入代理
	if _, err := proxySession.Proxy.Write(msg.Body); err != nil {
		return err
	}

	// 从代理读取响应
	buffer := make([]byte, 1024)
	n, err := proxySession.Proxy.Read(buffer)
	if err != nil {
		return err
	}

	// 发送响应回客户端
	if !ns.Conn.Write(neffos.Message{
		Body:     buffer[:n],
		IsNative: true,
	}) {
		return fmt.Errorf("failed to write message to websocket")
	}
	return nil
}

func handleDisconnect(ns *neffos.NSConn, msg neffos.Message) error {
	sessionID := ns.Conn.ID()
	if err := sessionManager.CloseSession(sessionID); err != nil {
		return err
	}
	return nil
}
