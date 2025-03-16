package models

type Car struct {
	RegistrationNumber string `gorm:"not null" json:"registration_number"`
	NationalNumber     string `gorm:"index;not null" json:"national_number"`
	CarPlate           string `gorm:"uniqueIndex;not null" json:"car_plate"`
	Model              string `json:"model"`
	Color              string `json:"color"`
	Type               string `json:"type"`
}
