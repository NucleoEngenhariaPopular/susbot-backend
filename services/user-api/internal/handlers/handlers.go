package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"user-api/internal/clients"
	"user-api/internal/database"
	"user-api/internal/models"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/users/")

		if path == "" {
			getAllUsers(w, r)
			return
		}

		if strings.HasPrefix(path, "cpf/") {
			getUserByCPF(w, r)
		} else {
			getUser(w, r)
		}
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if req.Name == "" || req.CPF == "" {
		respondWithError(w, http.StatusBadRequest, "Name and CPF are required")
		return
	}

	// Create user instance
	user := models.User{
		Name:         req.Name,
		CPF:          req.CPF,
		DateOfBirth:  req.DateOfBirth,
		PhoneNumber:  req.PhoneNumber,
		StreetName:   req.StreetName,
		StreetNumber: req.StreetNumber,
		Complement:   req.Complement,
		Neighborhood: req.Neighborhood,
		City:         req.City,
		State:        req.State,
		CEP:          req.CEP,
	}

	// Save to database
	if err := database.GetDB().Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusConflict, "CPF already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Get team information based on address
	teamInfo, err := lookupTeamInfo(user)
	if err != nil {
		log.Printf("Error looking up team info: %v", err)
		// Still return user info even if team lookup fails
		respondWithJSON(w, http.StatusOK, models.APIResponse{
			Success: true,
			Data:    user,
		})
		return
	}

	response := models.UserWithTeam{
		User: user,
		Team: teamInfo,
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func getUserByCPF(w http.ResponseWriter, r *http.Request) {
	cpf := strings.TrimPrefix(r.URL.Path, "/users/cpf/")
	cpf = strings.TrimSpace(cpf)

	var user models.User
	if err := database.GetDB().Where("cpf = ?", cpf).First(&user).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Get team information based on address
	teamInfo, err := lookupTeamInfo(user)
	if err != nil {
		log.Printf("Error looking up team info: %v", err)
		// Still return user info even if team lookup fails
		respondWithJSON(w, http.StatusOK, models.APIResponse{
			Success: true,
			Data:    user,
		})
		return
	}

	response := models.UserWithTeam{
		User: user,
		Team: teamInfo,
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Update only provided fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.StreetName != "" {
		user.StreetName = req.StreetName
	}
	if req.StreetNumber != "" {
		user.StreetNumber = req.StreetNumber
	}
	if req.Complement != "" {
		user.Complement = req.Complement
	}
	if req.Neighborhood != "" {
		user.Neighborhood = req.Neighborhood
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.State != "" {
		user.State = req.State
	}
	if req.CEP != "" {
		user.CEP = req.CEP
	}

	if err := database.GetDB().Save(&user).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	if err := database.GetDB().Delete(&user).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "User successfully deleted",
	})
}

// addressClient is a package-level variable that will be initialized in main
var addressClient *clients.AddressClient

// SetAddressClient initializes the address client
func SetAddressClient(client *clients.AddressClient) {
	addressClient = client
}

// Helper function to look up team information from the address-api
func lookupTeamInfo(user models.User) (models.TeamInfo, error) {
	if addressClient == nil {
		log.Printf("Address client not initialized")
		return models.TeamInfo{}, fmt.Errorf("address client not initialized")
	}

	log.Printf("Looking up team info for address: %s, %s, %s, %s",
		user.StreetName, user.StreetNumber, user.City, user.State)

	teamInfo, err := addressClient.GetTeamInfo(
		user.StreetName,
		user.StreetNumber,
		user.City,
		user.State,
	)
	if err != nil {
		log.Printf("Error getting team info: %v", err)
		return models.TeamInfo{}, fmt.Errorf("error getting team info: %v", err)
	}

	if teamInfo == nil {
		log.Printf("No team found for address")
		return models.TeamInfo{}, nil
	}

	log.Printf("Found team: %+v", teamInfo)
	return *teamInfo, nil
}

func getAllUsers(w http.ResponseWriter, _ *http.Request) {
	var users []models.User

	result := database.GetDB().Find(&users)
	if result.Error != nil {
		log.Printf("Error fetching users: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    users,
	})
}

// Helper functions for response handling
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.APIResponse{
		Success: false,
		Error:   message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
