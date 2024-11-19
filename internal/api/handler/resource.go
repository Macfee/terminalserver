package handler

import (
	"audit-system/internal/auth"
	"audit-system/internal/service"

	"github.com/kataras/iris/v12"
)

type GrantAccessRequest struct {
	UserID      uint              `json:"user_id"`
	Permissions []auth.Permission `json:"permissions"`
}

func CreateResource(ctx iris.Context) {
	var resource auth.Resource
	if err := ctx.ReadJSON(&resource); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	// 验证当前用户权限
	userID := ctx.Values().Get("userID").(uint)
	if !service.AuthService.CanManageResource(userID, resource.Type) {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"code":    403,
			"message": "没有资源管理权限",
		})
		return
	}

	if err := service.ResourceService.CreateResource(resource); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "创建资源失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "资源创建成功",
		"data":    resource,
	})
}

func GrantResourceAccess(ctx iris.Context) {
	resourceID := ctx.Params().Get("id")
	var req GrantAccessRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	// 验证当前用户权限
	userID := ctx.Values().Get("userID").(uint)
	if !service.AuthService.CanGrantAccess(userID, resourceID) {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"code":    403,
			"message": "没有授权权限",
		})
		return
	}

	if err := service.ResourceService.GrantAccess(resourceID, req.UserID, req.Permissions); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "授权失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "授权成功",
	})
}
