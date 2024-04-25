package config

import (
	"sync"

	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	JWTSecret string
	GinPort   string
	GinMode   string
	DbMutex   sync.Mutex
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read the config file")
	}
	JWTSecret = viper.GetString("jwt_secret")
	GinPort = viper.GetString("gin.port")
	GinMode = viper.GetString("gin.mode")
}

func InitDB() {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database server")
	}

	// 创建数据库
	DbMutex.Lock()
	DB.Exec("CREATE DATABASE IF NOT EXISTS road_patrol")
	DB.Exec("USE road_patrol")
	DbMutex.Unlock()

	// 连接到具体的数据库
	dsn = username + ":" + password + "@tcp(" + host + ":" + port + ")/road_patrol?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to road_patrol database")
	}

	DbMutex.Lock()
	err = DB.AutoMigrate(&model.Road{}, &model.User{}, &model.Patrol{}, &model.Report{})
	DbMutex.Unlock()
	if err != nil {
		panic("failed to auto migrate database")
	}
}
