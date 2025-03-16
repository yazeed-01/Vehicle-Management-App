package routes

import (
	"fmt"
	"time"

	"net/http"
	"vehicle_management/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		fmt.Printf("Request Method: %s, Path: %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Println("Headers:", c.Request.Header)
		c.Next()
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost",
			"http://localhost:8080",
			"http://localhost:57902",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/login", controllers.Login)

	auth := r.Group("/api")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/vehicles", controllers.GetUserVehicles)
		auth.GET("/user", controllers.GetUserInfo)
		auth.PUT("/user", controllers.UpdateUserInfo)
		auth.POST("/search-plate", controllers.SearchByPlate)
	}

	return r
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_jwt_secret"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("national_number", claims["national_number"])
		c.Next()
	}
}
