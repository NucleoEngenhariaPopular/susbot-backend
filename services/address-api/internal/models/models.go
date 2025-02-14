package models

import "time"

// UBS represents a Basic Health Unit
type UBS struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:200;not null"`
	Address   string    `json:"address" gorm:"size:200;not null"`
	City      string    `json:"city" gorm:"size:100;not null"`
	State     string    `json:"state" gorm:"size:2;not null"`
	CEP       string    `json:"cep" gorm:"size:8;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Teams     []Team    `json:"teams,omitempty" gorm:"foreignKey:UBSID"`
}

// Team represents a medical team within a UBS
type Team struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	UBSID     uint      `json:"ubs_id" gorm:"not null"`
	UBS       UBS       `json:"ubs,omitempty" gorm:"foreignKey:UBSID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StreetSegment represents a street segment assigned to a team
type StreetSegment struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	StreetName         string    `json:"street_name" gorm:"size:200;not null"`
	OriginalStreetName string    `json:"original_street_name" gorm:"size:200;not null"`
	StreetType         string    `json:"street_type" gorm:"size:50;not null"` // Rua, Avenida, etc.
	Neighborhood       string    `json:"neighborhood" gorm:"size:100;not null"`
	City               string    `json:"city" gorm:"size:100;not null"`
	State              string    `json:"state" gorm:"size:2;not null"`
	StartNumber        int       `json:"start_number"`
	EndNumber          int       `json:"end_number"`
	CEPPrefix          string    `json:"cep_prefix" gorm:"size:5"` // First 5 digits of CEP
	EvenOdd            string    `json:"even_odd" gorm:"size:4"`   // 'even', 'odd', or 'all'
	TeamID             uint      `json:"team_id" gorm:"not null"`
	Team               Team      `json:"team,omitempty" gorm:"foreignKey:TeamID"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Request/Response structures
type CreateUBSRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	CEP     string `json:"cep" binding:"required"`
}

type CreateTeamRequest struct {
	Name  string `json:"name" binding:"required"`
	UBSID uint   `json:"ubs_id" binding:"required"`
}

type CreateStreetSegmentRequest struct {
	StreetName   string `json:"street_name" binding:"required"`
	StreetType   string `json:"street_type" binding:"required"`
	Neighborhood string `json:"neighborhood" binding:"required"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	StartNumber  int    `json:"start_number"`
	EndNumber    int    `json:"end_number"`
	CEPPrefix    string `json:"cep_prefix"`
	EvenOdd      string `json:"even_odd" binding:"required"`
	TeamID       uint   `json:"team_id" binding:"required"`
}

type AddressSearchRequest struct {
	StreetName string `json:"street_name" binding:"required"`
	Number     int    `json:"number" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
}

type AddressSearchResponse struct {
	StreetSegment StreetSegment `json:"street_segment"`
	Team          Team          `json:"team"`
	UBS           UBS           `json:"ubs"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

