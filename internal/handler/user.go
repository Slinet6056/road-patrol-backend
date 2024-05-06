package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/Slinet6056/road-patrol-backend/internal/model"
	"github.com/Slinet6056/road-patrol-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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

	tenantID := c.Query("tenant_id")

	userChan := make(chan model.User)
	errChan := make(chan error)

	go func() {
		var user model.User
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ? AND username = ? AND password = ?", tenantID, loginParams.Username, loginParams.Password).First(&user)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		userChan <- user
	}()

	select {
	case user := <-userChan:
		exp := time.Now().Add(time.Hour * 2)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"role":     user.Role,
			"exp":      exp.Unix(),
		})

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
	case err := <-errChan:
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户名或密码错误"})
		logger.Error(err.Error())
	}
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

	usersChan := make(chan []model.User)
	errChan := make(chan error)

	go func() {
		var users []model.User
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Find(&users)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		usersChan <- users
	}()

	select {
	case users := <-usersChan:
		c.JSON(200, users)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
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

	userChan := make(chan model.User)
	errChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Create(&user)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		var createdUser model.User
		config.DbMutex.Lock()
		config.DB.Where("id = ?", user.ID).First(&createdUser)
		config.DbMutex.Unlock()
		userChan <- createdUser
	}()

	select {
	case createdUser := <-userChan:
		c.JSON(201, createdUser)
	case err := <-errChan:
		c.JSON(500, gin.H{"error": err.Error()})
	}
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

	userChan := make(chan model.User)
	errChan := make(chan error)

	go func() {
		var existingUser model.User
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ? AND id = ?", tenantID, id).First(&existingUser)
		config.DbMutex.Unlock()
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errChan <- errors.New("no user found with given ID")
			} else {
				errChan <- result.Error
			}
			return
		}
		config.DbMutex.Lock()
		result = config.DB.Where("tenant_id = ?", tenantID).Model(&model.User{}).Where("id = ?", id).Updates(user)
		config.DbMutex.Unlock()
		if result.Error != nil {
			errChan <- result.Error
			return
		}
		if result.RowsAffected == 0 {
			userChan <- existingUser
		} else {
			var updatedUser model.User
			config.DbMutex.Lock()
			config.DB.Where("id = ?", id).First(&updatedUser)
			config.DbMutex.Unlock()
			userChan <- updatedUser
		}
	}()

	select {
	case user := <-userChan:
		if user.ID == 0 {
			c.JSON(200, gin.H{"message": "No changes made"})
		} else {
			c.JSON(200, user)
		}
	case err := <-errChan:
		if err.Error() == "no user found with given ID" {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	id := c.Param("id")

	resultChan := make(chan error)

	go func() {
		config.DbMutex.Lock()
		result := config.DB.Where("tenant_id = ?", tenantID).Delete(&model.User{}, id)
		config.DbMutex.Unlock()
		resultChan <- result.Error
	}()

	if err := <-resultChan; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted"})
}
