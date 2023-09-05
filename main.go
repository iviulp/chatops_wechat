package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"yuguosheng/int/mychatops/config"
	"yuguosheng/int/mychatops/dao"
	"yuguosheng/int/mychatops/middleware"
	"yuguosheng/int/mychatops/routers"
)

func loadGin() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	// 使用中间件
	r.Use(middleware.LoggerToFile())
	r.Use(cors.Default())
	// 注册路由
	routers.LoadRouters(r)
	dao.LoadDatabase()
	return r
}

func main() {
	r := loadGin()
	_ = r.Run(":" + config.GetSystemConf().Port)
}
