package model

// Patrol 定义巡检任务的结构体
type Patrol struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	TenantID    uint   `json:"tenant_id"`
	RoadID      uint   `json:"road_id"`
	InspectorID uint   `json:"inspector_id"`
	Date        string `json:"date"`
	Status      string `json:"status"`

	Road      Road `gorm:"foreignKey:RoadID"`
	Inspector User `gorm:"foreignKey:InspectorID"`
}
