package common

type Address struct {
	StreetName   string `json:"street_name" validate:"required,max=200"`
	StreetNumber string `json:"street_number" validate:"required,max=20"`
	Complement   string `json:"complement" validate:"max=100"`
	Neighborhood string `json:"neighborhood" validate:"required,max=100"`
	City         string `json:"city" validate:"required,max=100"`
	State        string `json:"state" validate:"required,len=2"`
	CEP          string `json:"cep" validate:"required,len=8"`
}
