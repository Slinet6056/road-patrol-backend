package handler

import (
	"errors"
	"strconv"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetPatrols 获取所有巡检任务
func GetPatrols(c *gin.Context) {
	tenantID := c.Query("tenant_id")

	patrolChan := make(chan []model.Patrol)
	errChan := make(chan error)

	go func() {
		var patrols []model.Patrol
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Find(&patrols)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		patrolChan <- patrols
	}()

	select {
	case patrols := <-patrolChan:
		c.JSON(200, patrols)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
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

	patrolChan := make(chan model.Patrol)
	errChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Create(&patrol)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		var createdPatrol model.Patrol
		config.DbMutex.Lock()
		config.DB.Where("id = ?", patrol.ID).First(&createdPatrol)
		config.DbMutex.Unlock()
		patrolChan <- createdPatrol
	}()

	select {
	case createdPatrol := <-patrolChan:
		c.JSON(201, createdPatrol)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
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

	patrolChan := make(chan model.Patrol)
	errChan := make(chan error)

	go func() {
		var existingPatrol model.Patrol
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ? AND id = ?", tenantID, id).First(&existingPatrol)
		config.DbMutex.Unlock()
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errChan <- errors.New("no patrol found with given ID")
			} else {
				errChan <- result.Error
			}
			return
		}
		config.DbMutex.Lock()
		result = config.DB.Where("tenant_id = ?", tenantID).Model(&model.Patrol{}).Where("id = ?", id).Updates(patrol)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		if result.RowsAffected == 0 {
			patrolChan <- existingPatrol
		} else {
			var updatedPatrol model.Patrol
			config.DbMutex.Lock()
			config.DB.Where("id = ?", id).First(&updatedPatrol)
			config.DbMutex.Unlock()
			patrolChan <- updatedPatrol
		}
	}()

	select {
	case patrol := <-patrolChan:
		if patrol.ID == 0 {
			c.JSON(200, gin.H{"message": "No fields updated", "patrol": patrol})
		} else {
			c.JSON(200, patrol)
		}
	case err := <-errChan:
		if err.Error() == "no patrol found with given ID" {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}

// DeletePatrol 删除巡检任务
func DeletePatrol(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")

	resultChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Patrol{}, id)
		config.DbMutex.Unlock()
		resultChan <- result.Error
	}()

	if err := <-resultChan; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Patrol deleted"})
}
