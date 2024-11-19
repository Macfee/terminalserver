package handler

import (
	"audit-system/internal/monitor"
	"audit-system/internal/service"
	"encoding/json"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
)

var sessionMonitor *monitor.SessionMonitor

func init() {
	sessionMonitor = monitor.NewSessionMonitor()
}

func HandleMonitorWebSocket(ctx iris.Context) {
	sessionID := ctx.Params().Get("id")

	// 检查权限
	userID := ctx.Values().Get("userID").(uint)
	if !service.AuthService.CanMonitorSession(userID, sessionID) {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	// 订阅会话事件
	eventChan := sessionMonitor.Subscribe(sessionID)
	defer sessionMonitor.Unsubscribe(sessionID, eventChan)

	// 升级到WebSocket连接
	conn := websocket.Upgrade(ctx, websocket.DefaultIDGenerator, wsServer)
	if conn == nil {
		return
	}

	// 处理监控事件
	go func() {
		for event := range eventChan {
			eventJSON, _ := json.Marshal(event)
			conn.Write(neffos.Message{
				Body:     eventJSON,
				IsNative: true,
			})
		}
	}()
}

// 在会话处理中添加事件广播
func broadcastSessionEvent(sessionID string, eventType string, data interface{}) {
	event := monitor.SessionEvent{
		SessionID: sessionID,
		Type:      eventType,
		Data:      data,
		Time:      time.Now(),
	}
	sessionMonitor.BroadcastEvent(event)
}
