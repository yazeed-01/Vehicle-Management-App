package initializers

import (
	"log"
	"vehicle_management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	dsn := "host=localhost user=postgres password=PASSWORD dbname=car port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Car{}, &models.History{})

	DB.Exec("CREATE INDEX idx_users_national_number ON users(national_number)")
	DB.Exec("CREATE INDEX idx_cars_national_number ON cars(national_number)")
	DB.Exec("CREATE INDEX idx_cars_car_plate ON cars(car_plate)")
}
