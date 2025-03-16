package controllers

import (
	"net/http"
	"time"
	"vehicle_management/initializers"
	"vehicle_management/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	NationalNumber string `json:"national_number"`
	SpecificKey    string `json:"specific_key"`
}

var jwtSecret = []byte("your_jwt_secret")

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	cacheKey := "user:" + req.NationalNumber

	if val, err := initializers.RedisClient.Get(initializers.RedisCtx, cacheKey).Result(); err == nil {
		// If token is cached, fetch user data to include phone_number and email
		if err := initializers.DB.Where("national_number = ?", req.NationalNumber).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		/*
			if user.SpecificKey != req.SpecificKey {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
		*/
		c.JSON(http.StatusOK, gin.H{
			"token":        val,
			"phone_number": user.PhoneNumber,
			"email":        user.Email,
		})
		return
	}

	if err := initializers.DB.Where("national_number = ?", req.NationalNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	/*
		if user.SpecificKey != req.SpecificKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	*/

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"national_number": user.NationalNumber,
		"exp":             time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	err = initializers.RedisClient.Set(initializers.RedisCtx, cacheKey, tokenString, 24*time.Hour).Err()
	if err != nil {

	}

	c.JSON(http.StatusOK, gin.H{
		"token":        tokenString,
		"phone_number": user.PhoneNumber,
		"email":        user.Email,
	})
}
