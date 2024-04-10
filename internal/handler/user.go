package handler

import (
	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
)

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
	var users []model.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, users)
}

// AddUser 添加新的用户
func AddUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(201, user)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	var user model.User
	id := c.Param("id")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result := config.DB.Model(&model.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No user found with given ID"})
		return
	}
	c.JSON(200, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&model.User{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "No user found with given ID"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}
