package main

import (
	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	router.Run() // 默认在8080端口监听
}
