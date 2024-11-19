package handler

import (
	"audit-system/internal/service"
	"audit-system/pkg/utils/jwt"

	"github.com/kataras/iris/v12"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx iris.Context) {
	var req LoginRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	user, err := service.AuthService.ValidateUser(req.Username, req.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role.Name)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "生成token失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": iris.Map{
			"token": token,
			"user":  user,
		},
	})
}
