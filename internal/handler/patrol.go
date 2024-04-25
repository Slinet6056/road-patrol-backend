package handler

import (
	"errors"
	"gorm.io/gorm"
	"strconv"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// GetPatrols 获取所有巡检任务
func GetPatrols(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var patrols []model.Patrol
	result := config.DB.Where("tenant_id = ?", tenantID).Find(&patrols)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, patrols)
}

// AddPatrol 添加新的巡检任务
func AddPatrol(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var patrol model.Patrol
	if err := c.ShouldBindJSON(&patrol); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	patrol.TenantID = uint(parsedTenantID)
	result := config.DB.Where("tenant_id = ?", tenantID).Create(&patrol)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	var createdPatrol model.Patrol
	config.DB.Where("id = ?", patrol.ID).First(&createdPatrol)
	c.JSON(201, createdPatrol)
}

// UpdatePatrol 更新巡检任务
func UpdatePatrol(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var patrol model.Patrol
	id := c.Param("id")
	if err := c.ShouldBindJSON(&patrol); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	patrol.TenantID = uint(parsedTenantID)
	var existingPatrol model.Patrol
	result := config.DB.Where("tenant_id = ? AND id = ?", tenantID, id).First(&existingPatrol)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "No patrol found with given ID"})
		} else {
			c.JSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}
	result = config.DB.Where("tenant_id = ?", tenantID).Model(&model.Patrol{}).Where("id = ?", id).Updates(patrol)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(200, gin.H{"message": "No patrol updated"})
	} else {
		var updatedPatrol model.Patrol
		config.DB.Where("id = ?", id).First(&updatedPatrol)
		c.JSON(200, updatedPatrol)
	}
}

// DeletePatrol 删除巡检任务
func DeletePatrol(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")
	result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Patrol{}, id)
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
