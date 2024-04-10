package model

// Patrol 定义巡检任务的结构体
type Patrol struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	RoadID      uint   `json:"road_id"`
	InspectorID uint   `json:"inspector_id"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}
