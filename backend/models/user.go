package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NationalNumber     string `gorm:"uniqueIndex;not null" json:"national_number"`
	Name               string `json:"name"`
	RegistrationNumber string `json:"registration_number"`
	PhoneNumber        string `json:"phone_number"`
	Email              string `json:"email"`
}
