package handler

import (
	"github.com/Slinet6056/road-patrol-backend/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Login 用户登录
func Login(c *gin.Context) {
	var loginParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginParams); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "服务器错误"})
		logger.Error(err.Error())
		return
	}

	// 验证用户名和密码
	var user model.User
	result := config.DB.Where("username = ? AND password = ?", loginParams.Username, loginParams.Password).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	exp := time.Now().Add(time.Hour * 2)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      exp.Unix(),
	})

	// 生成刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	refreshTokenString, _ := refreshToken.SignedString([]byte(config.JWTSecret))

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		logger.Error("Could not generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"username":     user.Username,
			"roles":        []string{user.Role},
			"accessToken":  tokenString,
			"refreshToken": refreshTokenString,
			"expires":      exp.Format("2006/01/02 15:04:05"),
		},
	})
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	var tokenParams struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&tokenParams); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		logger.Error(err.Error())
		return
	}

	// 解析刷新令牌
	token, err := jwt.Parse(tokenParams.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired refresh token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := time.Now().Add(time.Hour * 2)
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": claims["username"],
			"role":     claims["role"],
			"exp":      exp.Unix(),
		})
		newTokenString, err := newToken.SignedString([]byte(config.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not generate token"})
			logger.Error("Could not generate token")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"accessToken":  newTokenString,
				"refreshToken": tokenParams.RefreshToken,
				"expires":      exp.Format("2006/01/02 15:04:05"),
			},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired refresh token"})
	}
}

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var users []model.User
	result := config.DB.Where("tenant_id = ?", tenantID).Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, users)
}

// AddUser 添加新的用户
func AddUser(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	user.TenantID = uint(parsedTenantID)
	result := config.DB.Where("tenant_id = ?", tenantID).Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(201, user)
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	var user model.User
	id := c.Param("id")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	parsedTenantID, _ := strconv.ParseUint(tenantID, 10, 64)
	user.TenantID = uint(parsedTenantID)
	result := config.DB.Where("tenant_id = ?", tenantID).Model(&model.User{}).Where("id = ?", id).Updates(user)
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
	tenantID := c.Query("tenant_id")
	id := c.Param("id")
	result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.User{}, id)
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
