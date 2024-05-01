package model

import "time"

// Plan 定义巡检任务的结构体
type Plan struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenant_id"`
	InspectorID uint      `json:"inspector_id"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`

	Inspector User `gorm:"foreignKey:InspectorID"`
}
