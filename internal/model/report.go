package model

import "time"

// Report 定义巡检报告的结构体
type Report struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TenantID  uint      `json:"tenant_id"`
	PatrolID  uint      `json:"patrol_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Patrol Patrol `gorm:"foreignKey:PatrolID"`
}
