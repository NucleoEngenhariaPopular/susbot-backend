package models

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:200;not null"`
	CPF         string    `json:"cpf" gorm:"size:11;unique;not null"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PhoneNumber string    `json:"phone_number" gorm:"size:20"`

	// Address information
	StreetName   string `json:"street_name" gorm:"size:200;not null"`
	StreetNumber string `json:"street_number" gorm:"size:20;not null"`
	Complement   string `json:"complement" gorm:"size:100"`
	Neighborhood string `json:"neighborhood" gorm:"size:100;not null"`
	City         string `json:"city" gorm:"size:100;not null"`
	State        string `json:"state" gorm:"size:2;not null"`
	CEP          string `json:"cep" gorm:"size:8;not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request/Response structures
type CreateUserRequest struct {
	Name         string    `json:"name" binding:"required"`
	CPF          string    `json:"cpf" binding:"required"`
	DateOfBirth  time.Time `json:"date_of_birth" binding:"required"`
	PhoneNumber  string    `json:"phone_number"`
	StreetName   string    `json:"street_name" binding:"required"`
	StreetNumber string    `json:"street_number" binding:"required"`
	Complement   string    `json:"complement"`
	Neighborhood string    `json:"neighborhood" binding:"required"`
	City         string    `json:"city" binding:"required"`
	State        string    `json:"state" binding:"required"`
	CEP          string    `json:"cep" binding:"required"`
}

type UpdateUserRequest struct {
	Name         string `json:"name"`
	PhoneNumber  string `json:"phone_number"`
	StreetName   string `json:"street_name"`
	StreetNumber string `json:"street_number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	CEP          string `json:"cep"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// UserWithTeam represents a user with their associated healthcare team information
type UserWithTeam struct {
	User User     `json:"user"`
	Team TeamInfo `json:"team"`
}

// TeamInfo represents the minimal team information we need
type TeamInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	UBSName string `json:"ubs_name"`
}

