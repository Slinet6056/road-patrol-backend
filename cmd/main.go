package main

import (
	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/handler"
	"github.com/Slinet6056/road-patrol-backend/pkg/logger"
	"github.com/Slinet6056/road-patrol-backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig() // 初始化配置
	config.InitDB()     // 初始化数据库连接
	err := logger.Init()
	if err != nil {
		return
	}

	gin.SetMode(config.GinMode)
	router := gin.Default()

	// 注册登录路由
	router.POST("/login", handler.Login)

	// 注册刷新令牌路由
	router.POST("/refresh-token", handler.RefreshToken)

	authorizedAdmin := router.Group("/")
	authorizedAdmin.Use(middleware.JWTAuth([]string{"admin"}))
	{
		authorizedAdmin.POST("/road", handler.AddRoad)
		authorizedAdmin.PUT("/road/:id", handler.UpdateRoad)
		authorizedAdmin.DELETE("/road/:id", handler.DeleteRoad)

		authorizedAdmin.GET("/users", handler.GetUsers)
		authorizedAdmin.POST("/user", handler.AddUser)
		authorizedAdmin.PUT("/user/:id", handler.UpdateUser)
		authorizedAdmin.DELETE("/user/:id", handler.DeleteUser)
	}

	authorizedInspector := router.Group("/")
	authorizedInspector.Use(middleware.JWTAuth([]string{"admin", "inspector"}))
	{
		authorizedInspector.GET("/roads", handler.GetRoads)

		authorizedInspector.GET("/plans", handler.GetPlans)
		authorizedInspector.POST("/plan", handler.AddPlan)
		authorizedInspector.PUT("/plan/:id", handler.UpdatePlan)
		authorizedInspector.DELETE("/plan/:id", handler.DeletePlan)

		authorizedInspector.GET("/reports", handler.GetReports)
		authorizedInspector.POST("/report", handler.AddReport)
		authorizedInspector.PUT("/report/:id", handler.UpdateReport)
		authorizedInspector.DELETE("/report/:id", handler.DeleteReport)
	}

	err = router.Run(":" + config.GinPort)
	if err != nil {
		return
	}
}
