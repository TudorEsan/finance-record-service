package routes

import (
	controller "App/controllers"
	"App/middleware"

	"github.com/gin-gonic/gin"
)

func NetWorthRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.Use(middleware.VerifyAuth())
	incomingRoutes.POST("", controller.AddRecord())
	incomingRoutes.DELETE(":id", controller.DeleteRecord())
	incomingRoutes.GET(":id", controller.GetRecord())
	incomingRoutes.PUT(":id", controller.UpdateRecord())
	incomingRoutes.GET("", controller.GetRecords())
	incomingRoutes.GET("count", controller.GetRecordCount())
}
