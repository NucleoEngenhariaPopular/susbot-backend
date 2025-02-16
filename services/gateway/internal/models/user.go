package models

import "time"

type User struct {
	Name         string    `json:"name"`
	CPF          string    `json:"cpf"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PhoneNumber  string    `json:"phone_number"`
	StreetName   string    `json:"street_name"`
	StreetNumber string    `json:"street_number"`
	Complement   string    `json:"complement"`
	Neighborhood string    `json:"neighborhood"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	CEP          string    `json:"cep"`
}

type UserData struct {
	User
	// MSG nao vai no pacote json
	Msg *TwilioMessage `json:"-"`
}

