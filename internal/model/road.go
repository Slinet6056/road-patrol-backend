package model

// Road 定义道路信息的结构体
type Road struct {
	ID               uint    `json:"id" gorm:"primaryKey"`
	Name             string  `json:"name"`
	Latitude         float64 `json:"latitude"`  // 纬度
	Longitude        float64 `json:"longitude"` // 经度
	Length           float64 `json:"length"`
	Type             string  `json:"type"`
	SurfaceMaterial  string  `json:"surface_material"`
	ConstructionYear int     `json:"construction_year"`
}
