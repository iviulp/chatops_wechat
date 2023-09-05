package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yuguosheng/int/mychatops/controller"
)

func RegistryWXRouter(r *gin.Engine) {
	gptApi := r.Group("/gpt")
	{
		gptApi.GET("", controller.VerifyCallBack)
		gptApi.POST("", controller.WxChatCommand)
	}
}

func TestRouter(r *gin.Engine) {
	testGroup := r.Group("/test")
	testGroup.GET("", func(context *gin.Context) {
		context.String(http.StatusOK, "Pong")
	})
}
