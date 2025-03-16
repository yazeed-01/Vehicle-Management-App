package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Constants
const (
	NUM_USERS        = 1_000_000 // Number of users
	NUM_CARS         = 800_000   // Number of cars
	PHONE_EMAIL_PROB = 0.7       // Probability a user has phone/email
)

// User represents a user record
type User struct {
	NationalNumber     string
	Name               string
	RegistrationNumber string // Comma-separated list of registration numbers for this user
	PhoneNumber        string
	Email              string
}

// Car represents a car record
type Car struct {
	RegistrationNumber string
	NationalNumber     string
	CarPlate           string
	Model              string
	Color              string
	Type               string
}

// Helper functions
func generateNationalNumber(usedNationalNumbers map[string]bool) string {
	for {
		nationalNumber := fmt.Sprintf("%010d", rand.Intn(1_000_000_000_0))
		if !usedNationalNumbers[nationalNumber] {
			usedNationalNumbers[nationalNumber] = true
			return nationalNumber
		}
	}
}

func generateRegistrationNumber(usedRegistrationNumbers map[string]bool) string {
	for {
		// Generate a 9-digit registration number
		registrationNumber := fmt.Sprintf("%09d", rand.Intn(1_000_000_000))
		if !usedRegistrationNumbers[registrationNumber] {
			usedRegistrationNumbers[registrationNumber] = true
			return registrationNumber
		}
	}
}

func generateCarPlate(usedPlates map[string]bool) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	letterChoice := func() byte { return letters[rand.Intn(len(letters))] }
	digitChoice := func() byte { return digits[rand.Intn(len(digits))] }
	for {
		carPlate := fmt.Sprintf("%c%c%c%c%c%c%c",
			digitChoice(), digitChoice(),
			letterChoice(), letterChoice(), letterChoice(),
			digitChoice(), digitChoice())
		if !usedPlates[carPlate] {
			usedPlates[carPlate] = true
			return carPlate
		}
	}
}

func generateName() string {
	firstNames := []string{"John", "Jane", "Michael", "Emily", "David", "Sarah", "Robert", "Lisa"}
	lastNames := []string{"Smith", "Johnson", "Brown", "Davis", "Wilson", "Clark", "Lewis", "Taylor"}
	return fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))])
}

func generatePhone() string {
	digits := "0123456789"
	digitChoice := func() byte { return digits[rand.Intn(len(digits))] }
	return fmt.Sprintf("+1-%c%c%c-%c%c%c-%c%c%c%c",
		digitChoice(), digitChoice(), digitChoice(),
		digitChoice(), digitChoice(), digitChoice(),
		digitChoice(), digitChoice(), digitChoice(), digitChoice())
}

func generateEmail(name string) string {
	domains := []string{"gmail.com", "yahoo.com", "hotmail.com", "example.com"}
	emailName := strings.ReplaceAll(strings.ToLower(name), " ", ".")
	return fmt.Sprintf("%s@%s", emailName, domains[rand.Intn(len(domains))])
}

func generateUserData(usedRegistrationNumbers map[string]bool, numCars int) []User {
	users := make([]User, 0, NUM_USERS)
	usedNationalNumbers := make(map[string]bool)

	// First, ensure we have enough car owners to cover the number of cars
	// We'll distribute cars among users, allowing some users to own multiple cars
	// Let's assume a user can own between 1 and 5 cars, with a skewed distribution
	carOwners := make([]User, 0)
	carsAssigned := 0

	// Generate users who own cars first
	for carsAssigned < numCars {
		nationalNumber := generateNationalNumber(usedNationalNumbers)
		name := generateName()
		phone := ""
		email := ""
		if rand.Float64() < PHONE_EMAIL_PROB {
			phone = generatePhone()
			email = generateEmail(name)
		}

		// Decide how many cars this user will own (skewed towards fewer cars)
		numCarsForUser := 1 + rand.Intn(5) // 1 to 5 cars
		if numCarsForUser > numCars-carsAssigned {
			numCarsForUser = numCars - carsAssigned // Don't assign more cars than needed
		}

		registrationNumbers := make([]string, numCarsForUser)
		for j := 0; j < numCarsForUser; j++ {
			registrationNumbers[j] = generateRegistrationNumber(usedRegistrationNumbers)
		}

		carOwners = append(carOwners, User{
			NationalNumber:     nationalNumber,
			Name:               name,
			RegistrationNumber: strings.Join(registrationNumbers, ","), // Comma-separated list
			PhoneNumber:        phone,
			Email:              email,
		})
		carsAssigned += numCarsForUser
	}

	// Add car owners to the users slice
	users = append(users, carOwners...)

	// Generate remaining users without cars
	for len(users) < NUM_USERS {
		nationalNumber := generateNationalNumber(usedNationalNumbers)
		name := generateName()
		phone := ""
		email := ""
		if rand.Float64() < PHONE_EMAIL_PROB {
			phone = generatePhone()
			email = generateEmail(name)
		}

		users = append(users, User{
			NationalNumber:     nationalNumber,
			Name:               name,
			RegistrationNumber: "",
			PhoneNumber:        phone,
			Email:              email,
		})
	}

	return users
}

func generateCarData(users []User, usedRegistrationNumbers map[string]bool) []Car {
	cars := make([]Car, 0, NUM_CARS)
	usedPlates := make(map[string]bool)

	// Create a map of national_number to registration_numbers for easy lookup
	userCars := make(map[string][]string)
	for _, user := range users {
		if user.RegistrationNumber != "" {
			registrationNumbers := strings.Split(user.RegistrationNumber, ",")
			userCars[user.NationalNumber] = registrationNumbers
		}
	}

	// Flatten the list of registration numbers and their corresponding national numbers
	type carOwner struct {
		NationalNumber     string
		RegistrationNumber string
	}
	var carOwners []carOwner
	for nationalNumber, regNumbers := range userCars {
		for _, regNumber := range regNumbers {
			carOwners = append(carOwners, carOwner{
				NationalNumber:     nationalNumber,
				RegistrationNumber: regNumber,
			})
		}
	}

	// Ensure we have enough registration numbers to cover the number of cars
	if len(carOwners) < NUM_CARS {
		log.Fatalf("Not enough registration numbers (%d) to generate %d cars", len(carOwners), NUM_CARS)
	}

	// Shuffle car owners to ensure random assignment
	rand.Shuffle(len(carOwners), func(i, j int) {
		carOwners[i], carOwners[j] = carOwners[j], carOwners[i]
	})

	models := []string{"Toyota Camry", "Honda Civic", "Ford Mustang", "Tesla Model 3", "BMW X5"}
	colors := []string{"Red", "Blue", "Black", "White", "Silver"}
	carTypes := []string{"Sedan", "SUV", "Truck", "Coupe", "Hatchback"}

	for i := 0; i < NUM_CARS; i++ {
		owner := carOwners[i] // Use a unique registration number for each car
		carPlate := generateCarPlate(usedPlates)
		cars = append(cars, Car{
			RegistrationNumber: owner.RegistrationNumber,
			NationalNumber:     owner.NationalNumber,
			CarPlate:           carPlate,
			Model:              models[rand.Intn(len(models))],
			Color:              colors[rand.Intn(len(colors))],
			Type:               carTypes[rand.Intn(len(carTypes))],
		})
	}
	return cars
}

func saveToCSV(filename string, headers []string, data [][]string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		log.Fatalf("Failed to write headers to %s: %v", filename, err)
	}

	// Write data
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write row to %s: %v", filename, err)
		}
	}
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Generating user data...")
	usedRegistrationNumbers := make(map[string]bool)
	users := generateUserData(usedRegistrationNumbers, NUM_CARS)

	fmt.Println("Generating car data...")
	cars := generateCarData(users, usedRegistrationNumbers)

	// Prepare user data for CSV
	userHeaders := []string{"national_number", "name", "registration_number", "phone_number", "email"}
	userData := make([][]string, len(users))
	for i, user := range users {
		userData[i] = []string{user.NationalNumber, user.Name, user.RegistrationNumber, user.PhoneNumber, user.Email}
	}

	// Prepare car data for CSV
	carHeaders := []string{"registration_number", "national_number", "car_plate", "model", "color", "type"}
	carData := make([][]string, len(cars))
	for i, car := range cars {
		carData[i] = []string{car.RegistrationNumber, car.NationalNumber, car.CarPlate, car.Model, car.Color, car.Type}
	}

	// Save to CSV
	fmt.Println("Saving data to CSV files...")
	saveToCSV("users.csv", userHeaders, userData)
	saveToCSV("cars.csv", carHeaders, carData)
	fmt.Println("Data generation complete! Files saved as 'users.csv' and 'cars.csv'.")
}
