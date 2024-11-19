package handler

import (
	"audit-system/internal/service"

	"github.com/kataras/iris/v12"
)

func Logout(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code":    200,
		"message": "登出成功",
	})
}

func GetUserProfile(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	user, err := service.UserService.GetUserProfile(userID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": user,
	})
}
