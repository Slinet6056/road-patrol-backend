package handler

import (
	"errors"
	"strconv"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlanDetail struct {
	model.Plan
	RoadIDs []uint `json:"road_ids"`
}

// GetPlans 获取所有巡检任务及其关联的道路ID
func GetPlans(c *gin.Context) {
	tenantID := c.Query("tenant_id")

	planDetailChan := make(chan []PlanDetail)
	errChan := make(chan error)

	go func() {
		var plans []model.Plan
		var planDetails []PlanDetail

		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Find(&plans)
		config.DbMutex.Unlock()

		if result.Error != nil {
			errChan <- result.Error
			return
		}

		for _, plan := range plans {
			var roadIDs []uint
			config.DbMutex.Lock()
			config.DB.Model(&model.PlanRoad{}).Where("plan_id = ?", plan.ID).Pluck("road_id", &roadIDs)
			config.DbMutex.Unlock()

			planDetails = append(planDetails, PlanDetail{Plan: plan, RoadIDs: roadIDs})
		}

		planDetailChan <- planDetails
	}()

	select {
	case planDetails := <-planDetailChan:
		c.JSON(200, planDetails)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

// AddPlan 添加新的巡检任务及其关联的道路
func AddPlan(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var planDetail PlanDetail
	if err := c.ShouldBindJSON(&planDetail); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	planDetail.TenantID = uint(parsedTenantID)

	plans := make(chan model.Plan)
	errChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Create(&planDetail.Plan)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}

		for _, roadID := range planDetail.RoadIDs {
			config.DbMutex.Lock()
			config.DB.Create(&model.PlanRoad{PlanID: planDetail.ID, RoadID: roadID})
			config.DbMutex.Unlock()
		}

		plans <- planDetail.Plan
	}()

	select {
	case createdPlan := <-plans:
		c.JSON(201, createdPlan)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
}

// UpdatePlan 更新巡检任务及其关联的道路
func UpdatePlan(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var planDetail PlanDetail
	id := c.Param("id")
	if err := c.ShouldBindJSON(&planDetail); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	planDetail.TenantID = uint(parsedTenantID)

	planChan := make(chan model.Plan)
	errChan := make(chan error)

	go func() {
		var existingPlan model.Plan
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ? AND id = ?", tenantID, id).First(&existingPlan)
		config.DbMutex.Unlock()
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errChan <- errors.New("no plan found with given ID")
			} else {
				errChan <- result.Error
			}
			return
		}

		// 更新 Plan 表
		config.DbMutex.Lock()
		result = config.DB.Where("tenant_id = ?", tenantID).Model(&model.Plan{}).Where("id = ?", id).Updates(planDetail.Plan)
		config.DbMutex.Unlock()

		// 更新 PlanRoad 表
		config.DbMutex.Lock()
		config.DB.Where("plan_id = ?", id).Delete(&model.PlanRoad{})
		for _, roadID := range planDetail.RoadIDs {
			config.DB.Create(&model.PlanRoad{PlanID: planDetail.ID, RoadID: roadID})
		}
		config.DbMutex.Unlock()

		if result.Error != nil {
			errChan <- result.Error
			return
		}
		if result.RowsAffected == 0 {
			planChan <- existingPlan
		} else {
			var updatedPlan model.Plan
			config.DbMutex.Lock()
			config.DB.Where("id = ?", id).First(&updatedPlan)
			config.DbMutex.Unlock()
			planChan <- updatedPlan
		}
	}()

	select {
	case plan := <-planChan:
		if plan.ID == 0 {
			c.JSON(200, gin.H{"message": "No fields updated", "plan": plan})
		} else {
			c.JSON(200, plan)
		}
	case err := <-errChan:
		if err.Error() == "no plan found with given ID" {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}

// DeletePlan 删除巡检任务及其关联的道路
func DeletePlan(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")

	resultChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		// 先删除 PlanRoad 表中的关联数据
		config.DB.Where("plan_id = ?", id).Delete(&model.PlanRoad{})
		// 再删除 Plan 表中的数据
		result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.Plan{}, id)
		config.DbMutex.Unlock()
		resultChan <- result.Error
	}()

	if err := <-resultChan; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Plan deleted"})
}
