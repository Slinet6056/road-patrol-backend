package model

// User 定义用户的结构体
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	TenantID uint   `json:"tenant_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
