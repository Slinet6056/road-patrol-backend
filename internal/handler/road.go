package handler

import (
	"errors"
	"strconv"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRoads 获取所有道路信息
func GetRoads(c *gin.Context) {
	tenantID := c.Query("tenant_id")

	roadChan := make(chan []model.Road)
	errChan := make(chan error)

	go func() {
		var roads []model.Road
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Find(&roads)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		roadChan <- roads
	}()

	select {
	case roads := <-roadChan:
		c.JSON(200, roads)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

// AddRoad 添加新的道路信息
func AddRoad(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var road model.Road
	if err := c.ShouldBindJSON(&road); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	road.TenantID = uint(parsedTenantID)

	roadChan := make(chan model.Road)
	errChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Create(&road)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		var createdRoad model.Road
		config.DbMutex.Lock()
		config.DB.Where("id = ?", road.ID).First(&createdRoad)
		config.DbMutex.Unlock()
		roadChan <- createdRoad
	}()

	select {
	case createdRoad := <-roadChan:
		c.JSON(201, createdRoad)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
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
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	road.TenantID = uint(parsedTenantID)

	roadChan := make(chan model.Road)
	errChan := make(chan error)

	go func() {
		var existingRoad model.Road
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ? AND id = ?", tenantID, id).First(&existingRoad)
		config.DbMutex.Unlock()
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errChan <- errors.New("no road found with given ID")
			} else {
				errChan <- result.Error
			}
			return
		}
		config.DbMutex.Lock()
		result = config.DB.Where("tenant_id = ?", tenantID).Model(&model.Road{}).Where("id = ?", id).Updates(road)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		if result.RowsAffected == 0 {
			roadChan <- existingRoad
		} else {
			var updatedRoad model.Road
			config.DbMutex.Lock()
			config.DB.Where("id = ?", id).First(&updatedRoad)
			config.DbMutex.Unlock()
			roadChan <- updatedRoad
		}
	}()

	select {
	case road := <-roadChan:
		if road.ID == 0 {
			c.JSON(200, gin.H{"message": "No fields updated", "road": road})
		} else {
			c.JSON(200, road)
		}
	case err := <-errChan:
		if err.Error() == "no road found with given ID" {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}

// DeleteRoad 删除道路信息
func DeleteRoad(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")

	resultChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Road{}, id)
		config.DbMutex.Unlock()
		resultChan <- result.Error
	}()

	if err := <-resultChan; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Road deleted"})
}
