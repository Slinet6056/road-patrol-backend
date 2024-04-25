package handler

import (
	"errors"
	"strconv"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetReports 获取所有巡检报告
func GetReports(c *gin.Context) {
	tenantID := c.Query("tenant_id")

	reportChan := make(chan []model.Report)
	errChan := make(chan error)

	go func() {
		var reports []model.Report
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Find(&reports)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		reportChan <- reports
	}()

	select {
	case reports := <-reportChan:
		c.JSON(200, reports)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

// AddReport 添加新的巡检报告
func AddReport(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var report model.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	report.TenantID = uint(parsedTenantID)

	reportChan := make(chan model.Report)
	errChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Create(&report)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		var createdReport model.Report
		config.DbMutex.Lock()
		config.DB.Where("id = ?", report.ID).First(&createdReport)
		config.DbMutex.Unlock()
		reportChan <- createdReport
	}()

	select {
	case createdReport := <-reportChan:
		c.JSON(201, createdReport)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

// UpdateReport 更新巡检报告信息
func UpdateReport(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var report model.Report
	id := c.Param("id")
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	report.TenantID = uint(parsedTenantID)

	reportChan := make(chan model.Report)
	errChan := make(chan error)

	go func() {
		var existingReport model.Report
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ? AND id = ?", tenantID, id).First(&existingReport)
		config.DbMutex.Unlock()
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errChan <- errors.New("no report found with given ID")
			} else {
				errChan <- result.Error
			}
			return
		}
		config.DbMutex.Lock()
		result = config.DB.Where("tenant_id = ?", tenantID).Model(&model.Report{}).Where("id = ?", id).Updates(report)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		if result.RowsAffected == 0 {
			reportChan <- existingReport
		} else {
			var updatedReport model.Report
			config.DbMutex.Lock()
			config.DB.Where("id = ?", id).First(&updatedReport)
			config.DbMutex.Unlock()
			reportChan <- updatedReport
		}
	}()

	select {
	case report := <-reportChan:
		if report.ID == 0 {
			c.JSON(200, gin.H{"message": "No fields updated", "report": report})
		} else {
			c.JSON(200, report)
		}
	case err := <-errChan:
		if err.Error() == "no report found with given ID" {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}

// DeleteReport 删除巡检报告
func DeleteReport(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")

	resultChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Report{}, id)
		config.DbMutex.Unlock()
		resultChan <- result.Error
	}()

	if err := <-resultChan; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Report deleted"})
}
