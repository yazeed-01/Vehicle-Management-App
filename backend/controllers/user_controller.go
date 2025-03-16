package controllers

import (
	"net/http"
	"time"
	"vehicle_management/initializers"
	"vehicle_management/models"

	"github.com/gin-gonic/gin"
)

func GetUserVehicles(c *gin.Context) {
	nationalNumber := c.GetString("national_number") 

	var cars []models.Car
	cacheKey := "vehicles:" + nationalNumber

	if val, err := initializers.RedisClient.Get(initializers.RedisCtx, cacheKey).Result(); err == nil {
		c.JSON(http.StatusOK, gin.H{"vehicles": val})
		return
	}

	if err := initializers.DB.Where("national_number = ?", nationalNumber).Find(&cars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}

	initializers.RedisClient.Set(initializers.RedisCtx, cacheKey, cars, time.Hour)

	c.JSON(http.StatusOK, gin.H{"vehicles": cars})
}

type UpdateUserRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func UpdateUserInfo(c *gin.Context) {
	nationalNumber := c.GetString("national_number") 

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("national_number = ?", nationalNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Email = req.Email
	user.PhoneNumber = req.PhoneNumber

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		return
	}

	cacheKey := "user_info:" + nationalNumber
	initializers.RedisClient.Del(initializers.RedisCtx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "User info updated successfully"})
}
func GetUserInfo(c *gin.Context) {
	nationalNumber := c.GetString("national_number") 

	var user models.User
	cacheKey := "user_info:" + nationalNumber

	if val, err := initializers.RedisClient.Get(initializers.RedisCtx, cacheKey).Result(); err == nil {
		c.JSON(http.StatusOK, gin.H{"user": val})
		return
	}

	if err := initializers.DB.Where("national_number = ?", nationalNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	initializers.RedisClient.Set(initializers.RedisCtx, cacheKey, user, time.Hour)

	c.JSON(http.StatusOK, gin.H{"user": user})
}
