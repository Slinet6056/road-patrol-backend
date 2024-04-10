package model

// User 定义用户的结构体
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
