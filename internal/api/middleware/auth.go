package middleware

import (
	"audit-system/pkg/utils/jwt"
	"strings"

	"github.com/kataras/iris/v12"
)

func JWT() iris.Handler {
	return func(ctx iris.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{
				"code":    401,
				"message": "未提供认证信息",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{
				"code":    401,
				"message": "认证格式错误",
			})
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{
				"code":    401,
				"message": "无效的token",
			})
			return
		}

		// 将用户信息存储在上下文中
		ctx.Values().Set("userID", claims.UserID)
		ctx.Values().Set("username", claims.Username)
		ctx.Values().Set("role", claims.Role)

		ctx.Next()
	}
}
