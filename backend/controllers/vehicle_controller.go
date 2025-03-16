package controllers

import (
	"net/http"
	"time"
	"vehicle_management/initializers"
	"vehicle_management/models"

	"github.com/gin-gonic/gin"
)

type SearchPlateRequest struct {
	CarPlate string `json:"car_plate"`
}

func SearchByPlate(c *gin.Context) {
	var req SearchPlateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("national_number")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var car models.Car
	cacheKey := "car:" + req.CarPlate

	if _, err := initializers.RedisClient.Get(initializers.RedisCtx, cacheKey).Result(); err == nil {
		if err := initializers.DB.Where("car_plate = ?", req.CarPlate).First(&car).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}
	} else {
		if err := initializers.DB.Where("car_plate = ?", req.CarPlate).First(&car).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}
	}

	var user models.User
	if err := initializers.DB.Where("national_number = ?", car.NationalNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch owner info"})
		return
	}

	if user.PhoneNumber == "" && user.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Owner contact information not found (no phone number or email)"})
		return
	}

	// Prepare response
	response := gin.H{
		"vehicle": gin.H{
			"model":     car.Model,
			"color":     car.Color,
			"car_plate": car.CarPlate,
			"type":      car.Type,
		},
		"owner": gin.H{
			"name":         user.Name,
			"phone_number": user.PhoneNumber,
			"email":        user.Email,
		},
	}

	err := initializers.RedisClient.Set(initializers.RedisCtx, cacheKey, response, time.Hour).Err()
	if err != nil {
	}

	saveSearchHistory(userID.(string), car.NationalNumber)

	c.JSON(http.StatusOK, response)
}

func saveSearchHistory(userID, searchedUser string) {
	var history models.History

	result := initializers.DB.Where("user_id = ? AND searched_user = ?", userID, searchedUser).First(&history)
	if result.Error == nil {
		history.SearchCount++
		history.SearchTime = time.Now()
		initializers.DB.Save(&history)
	} else {
		history = models.History{
			UserID:       userID,
			SearchedUser: searchedUser,
			SearchTime:   time.Now(),
			SearchCount:  1,
		}
		initializers.DB.Create(&history)
	}
}
