package handler

import (
	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// GetPatrols 获取所有巡检任务
func GetPatrols(c *gin.Context) {
	var patrols []model.Patrol
	result := config.DB.Find(&patrols)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, patrols)
}

// AddPatrol 添加新的巡检任务
func AddPatrol(c *gin.Context) {
	var patrol model.Patrol
	if err := c.ShouldBindJSON(&patrol); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Create(&patrol)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(201, patrol)
}

// UpdatePatrol 更新巡检任务
func UpdatePatrol(c *gin.Context) {
	var patrol model.Patrol
	id := c.Param("id")
	if err := c.ShouldBindJSON(&patrol); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Model(&model.Patrol{}).Where("id = ?", id).Updates(patrol)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No patrol found with given ID"})
		return
	}
	c.JSON(200, patrol)
}

// DeletePatrol 删除巡检任务
func DeletePatrol(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&model.Patrol{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No patrol found with given ID"})
		return
	}
	c.JSON(200, gin.H{"message": "Patrol deleted"})
}
