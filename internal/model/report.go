package model

import "time"

// Report 定义巡检报告的结构体
type Report struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TenantID  uint      `json:"tenant_id"`
	PlanID    uint      `json:"plan_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Plan Plan `gorm:"foreignKey:PlanID"`
}
