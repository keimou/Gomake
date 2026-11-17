package ninex

import (
	v1 "gmkdb/api/v1"

	"github.com/gin-gonic/gin"
)

type CollectRouter struct{}

func (*CollectRouter) InitRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("collect")
	api := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		router.POST("login", api.Login)
		router.POST("captcha", api.Captcha)
	}
	return router
}
