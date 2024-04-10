package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 是全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	dsn := "rpUser:RoadPatrolUser@tcp(127.0.0.1:3306)/road_patrol?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
