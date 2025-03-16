package main

import (
	"vehicle_management/initializers"
	"vehicle_management/routes"
)

func main() {
	initializers.ConnectToDatabase()

	initializers.ConnectToRedis()

	r := routes.SetupRouter()

	r.Run(":8080")
}
