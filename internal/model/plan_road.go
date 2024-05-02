package model

// PlanRoad 定义计划和报告之间多对多关系的结构体
type PlanRoad struct {
	TenantID uint `json:"tenant_id" gorm:"primaryKey"`
	PlanID   uint `json:"plan_id" gorm:"primaryKey"`
	RoadID   uint `json:"road_id" gorm:"primaryKey"`

	Plan Plan `gorm:"foreignKey:PlanID"`
	Road Road `gorm:"foreignKey:RoadID"`
}
