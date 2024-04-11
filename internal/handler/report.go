package handler

import (
	"strconv"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// GetReports 获取所有巡检报告
func GetReports(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var reports []model.Report
	result := config.DB.Where("tenant_id = ?", tenantID).Find(&reports)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, reports)
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
	result := config.DB.Where("tenant_id = ?", tenantID).Create(&report)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(201, report)
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
	result := config.DB.Where("tenant_id = ?", tenantID).Model(&model.Report{}).Where("id = ?", id).Updates(report)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No report found with given ID"})
		return
	}
	c.JSON(200, report)
}

// DeleteReport 删除巡检报告
func DeleteReport(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")
	result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Report{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No report found with given ID"})
		return
	}
	c.JSON(200, gin.H{"message": "Report deleted"})
}
