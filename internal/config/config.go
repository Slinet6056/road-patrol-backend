package config

import (
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "rpUser:RoadPatrolUser@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database server")
	}

	// 创建数据库
	DB.Exec("CREATE DATABASE IF NOT EXISTS road_patrol")
	DB.Exec("USE road_patrol")

	// 连接到具体的数据库
	dsn = "rpUser:RoadPatrolUser@tcp(127.0.0.1:3306)/road_patrol?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to road_patrol database")
	}

	// 自动迁移模式
	DB.AutoMigrate(&model.Road{}, &model.User{}, &model.Patrol{}, &model.Report{})
}
