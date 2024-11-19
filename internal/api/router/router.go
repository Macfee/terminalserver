package router

import (
	"audit-system/internal/api/handler"
	"audit-system/internal/api/middleware"

	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app *iris.Application) {
	api := app.Party("/api")
	{
		// 会话管理接口
		sessions := api.Party("/sessions")
		{
			sessions.Use(middleware.JWT())
			sessions.Get("/", handler.ListSessions)
			sessions.Get("/{id}", handler.GetSession)
			sessions.Delete("/{id}", handler.TerminateSession)
			sessions.Get("/{id}/replay", handler.GetSessionReplay)

			// WebSocket 连接
			sessions.Get("/{id}/ws", handler.HandleWebSocket)
			sessions.Get("/{id}/monitor", handler.HandleMonitorWebSocket)
		}

		// 告警配置接口
		alerts := api.Party("/alerts")
		{
			alerts.Use(middleware.JWT())
			alerts.Get("/rules", handler.ListAlertRules)
			alerts.Post("/rules", handler.CreateAlertRule)
			alerts.Put("/rules/{id}", handler.UpdateAlertRule)
			alerts.Delete("/rules/{id}", handler.DeleteAlertRule)
		}

		// 资源管理接口
		resources := api.Party("/resources")
		{
			resources.Use(middleware.JWT())
			resources.Post("/", handler.CreateResource)
			resources.Post("/{id}/access", handler.GrantResourceAccess)
		}

		// 用户认证接口
		auth := api.Party("/auth")
		{
			auth.Post("/login", handler.Login)
			auth.Post("/logout", middleware.JWT(), handler.Logout)
			auth.Get("/profile", middleware.JWT(), handler.GetUserProfile)
		}

		// 主机管理接口
		hosts := api.Party("/hosts")
		{
			hosts.Use(middleware.JWT())
			hosts.Get("/", handler.ListHosts)
			hosts.Post("/", handler.CreateHost)
			hosts.Get("/{id}", handler.GetHost)
			hosts.Put("/{id}", handler.UpdateHost)
			hosts.Delete("/{id}", handler.DeleteHost)
		}
	}
}
