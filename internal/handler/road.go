package handler

import (
	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// GetRoads 获取所有道路信息
func GetRoads(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var roads []model.Road
	result := config.DB.Where("tenant_id = ?", tenantID).Find(&roads)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, roads)
}

// AddRoad 添加新的道路信息
func AddRoad(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var road model.Road
	if err := c.ShouldBindJSON(&road); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Where("tenant_id = ?", tenantID).Create(&road)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(201, road)
}

// UpdateRoad 更新道路信息
func UpdateRoad(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var road model.Road
	id := c.Param("id")
	if err := c.ShouldBindJSON(&road); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Where("tenant_id = ?", tenantID).Model(&model.Road{}).Where("id = ?", id).Updates(road)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No road found with given ID"})
		return
	}
	c.JSON(200, road)
}

// DeleteRoad 删除道路信息
func DeleteRoad(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")
	result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Road{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No road found with given ID"})
		return
	}
	c.JSON(200, gin.H{"message": "Road deleted"})
}
