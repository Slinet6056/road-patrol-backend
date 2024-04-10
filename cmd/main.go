package main

import (
	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB() // 初始化数据库连接
	router := gin.Default()

	// 注册道路的路由
	router.GET("/roads", handler.GetRoads)
	router.POST("/road", handler.AddRoad)
	router.PUT("/road/:id", handler.UpdateRoad)
	router.DELETE("/road/:id", handler.DeleteRoad)

	// 注册巡检任务的路由
	router.GET("/patrols", handler.GetPatrols)
	router.POST("/patrol", handler.AddPatrol)
	router.PUT("/patrol/:id", handler.UpdatePatrol)
	router.DELETE("/patrol/:id", handler.DeletePatrol)

	router.Run() // 默认在8080端口监听
}
