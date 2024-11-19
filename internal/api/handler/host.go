package handler

import (
	"audit-system/internal/service"
	"strconv"

	"github.com/kataras/iris/v12"
)

func ListHosts(ctx iris.Context) {
	page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.URLParamDefault("page_size", "10"))

	hosts, total, err := service.HostService.ListHosts(page, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "获取主机列表失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": iris.Map{
			"total": total,
			"hosts": hosts,
		},
	})
}

func CreateHost(ctx iris.Context) {
	var host service.Host
	if err := ctx.ReadJSON(&host); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	if err := service.HostService.CreateHost(&host); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "创建主机失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": host,
	})
}

func GetHost(ctx iris.Context) {
	hostID := ctx.Params().Get("id")
	host, err := service.HostService.GetHost(hostID)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"code":    404,
			"message": "主机不存在",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": 200,
		"data": host,
	})
}

func UpdateHost(ctx iris.Context) {
	hostID := ctx.Params().Get("id")
	var host service.Host
	if err := ctx.ReadJSON(&host); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"code":    400,
			"message": "无效的请求参数",
		})
		return
	}

	if err := service.HostService.UpdateHost(hostID, &host); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "更新主机失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "更新成功",
	})
}

func DeleteHost(ctx iris.Context) {
	hostID := ctx.Params().Get("id")
	if err := service.HostService.DeleteHost(hostID); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"code":    500,
			"message": "删除主机失败",
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":    200,
		"message": "删除成功",
	})
}
