package handlers

import (
	"address-api/internal/database"
	"address-api/internal/models"
	"address-api/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func HandleStreetSegments(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/streets/search" {
		findTeamByAddress(w, r)
		return
	}

	switch r.Method {
	case http.MethodPost:
		createStreetSegment(w, r)
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/streets/")
		if path == "" {
			listStreetSegments(w, r)
		} else {
			getStreetSegment(w, r)
		}
	case http.MethodPut:
		updateStreetSegment(w, r)
	case http.MethodDelete:
		deleteStreetSegment(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func createStreetSegment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateStreetSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Verify if team exists
	var team models.Team
	if err := database.GetDB().First(&team, req.TeamID).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Normalize address components
	normalizedStreetName := utils.NormalizeStreetName(req.StreetName)
	normalizedStreetType := utils.NormalizeStreetType(req.StreetType)

	segment := models.StreetSegment{
		StreetName:         normalizedStreetName,
		OriginalStreetName: req.StreetName,
		StreetType:         normalizedStreetType,
		Neighborhood:       strings.ToUpper(req.Neighborhood),
		City:               strings.ToUpper(req.City),
		State:              strings.ToUpper(req.State),
		StartNumber:        req.StartNumber,
		EndNumber:          req.EndNumber,
		CEPPrefix:          utils.NormalizeCEP(req.CEPPrefix),
		EvenOdd:            strings.ToLower(req.EvenOdd),
		TeamID:             req.TeamID,
	}

	// Basic validation
	if segment.StartNumber > segment.EndNumber {
		respondWithError(w, http.StatusBadRequest, "Start number cannot be greater than end number")
		return
	}

	if !isValidEvenOdd(segment.EvenOdd) {
		respondWithError(w, http.StatusBadRequest, "Invalid even/odd value. Must be 'even', 'odd', or 'all'")
		return
	}

	if err := database.GetDB().Create(&segment).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create street segment")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    segment,
	})
}

func listStreetSegments(w http.ResponseWriter, _ *http.Request) {
	var segments []models.StreetSegment
	if err := database.GetDB().Preload("Team.UBS").Find(&segments).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch street segments")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    segments,
	})
}

func updateStreetSegment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/streets/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid street segment ID")
		return
	}

	var req models.CreateStreetSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var segment models.StreetSegment
	if err := database.GetDB().First(&segment, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Street segment not found")
		return
	}

	// Verify if new team exists
	if err := database.GetDB().First(&models.Team{}, req.TeamID).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	// Normalize address components
	normalizedStreetName := utils.NormalizeStreetName(req.StreetName)
	normalizedStreetType := utils.NormalizeStreetType(req.StreetType)

	// Update fields
	segment.StreetName = normalizedStreetName
	segment.OriginalStreetName = req.StreetName
	segment.StreetType = normalizedStreetType
	segment.Neighborhood = strings.ToUpper(req.Neighborhood)
	segment.City = strings.ToUpper(req.City)
	segment.State = strings.ToUpper(req.State)
	segment.StartNumber = req.StartNumber
	segment.EndNumber = req.EndNumber
	segment.CEPPrefix = utils.NormalizeCEP(req.CEPPrefix)
	segment.EvenOdd = strings.ToLower(req.EvenOdd)
	segment.TeamID = req.TeamID

	// Validate updated data
	if segment.StartNumber > segment.EndNumber {
		respondWithError(w, http.StatusBadRequest, "Start number cannot be greater than end number")
		return
	}

	if !isValidEvenOdd(segment.EvenOdd) {
		respondWithError(w, http.StatusBadRequest, "Invalid even/odd value. Must be 'even', 'odd', or 'all'")
		return
	}

	if err := database.GetDB().Save(&segment).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update street segment")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    segment,
	})
}

func deleteStreetSegment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/streets/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid street segment ID")
		return
	}

	var segment models.StreetSegment
	if err := database.GetDB().First(&segment, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Street segment not found")
		return
	}

	if err := database.GetDB().Delete(&segment).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete street segment")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "Street segment successfully deleted",
	})
}

func getStreetSegment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/streets/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid street segment ID")
		return
	}

	var segment models.StreetSegment
	if err := database.GetDB().Preload("Team.UBS").First(&segment, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Street segment not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    segment,
	})
}

func isValidEvenOdd(value string) bool {
	validValues := map[string]bool{
		"even": true,
		"odd":  true,
		"all":  true,
	}
	return validValues[value]
}

func findTeamByAddress(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get query parameters
	streetName := r.URL.Query().Get("street")
	numberStr := r.URL.Query().Get("number")
	city := r.URL.Query().Get("city")
	state := r.URL.Query().Get("state")

	if streetName == "" || numberStr == "" || city == "" || state == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required parameters: street, number, city, state")
		return
	}

	// Convert house number to int
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid house number")
		return
	}

	// Normalize the input
	normalizedStreet := utils.NormalizeStreetName(streetName)
	normalizedCity := strings.ToUpper(city)
	normalizedState := strings.ToUpper(state)

	var segments []models.StreetSegment

	// Using similarity threshold of 0.3 (can be adjusted)
	// First, find all potential matches in the same city and state
	db := database.GetDB().
		Preload("Team").
		Preload("Team.UBS").
		Where("city = ? AND state = ?", normalizedCity, normalizedState).
		Where("similarity(street_name, ?) > 0.3", normalizedStreet)

	// Check number range and even/odd rules
	db = db.Where(
		"(start_number <= ? AND end_number >= ?) AND "+
			"(even_odd = 'all' OR "+
			"(even_odd = 'even' AND ? % 2 = 0) OR "+
			"(even_odd = 'odd' AND ? % 2 = 1))",
		number, number, number, number)

	if err := db.Find(&segments).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to search for address")
		return
	}

	if len(segments) == 0 {
		respondWithError(w, http.StatusNotFound, "No team found for this address")
		return
	}

	// Find the best match based on street name similarity
	var bestMatch models.StreetSegment
	bestSimilarity := -1.0

	for _, segment := range segments {
		similarity := calculateSimilarity(segment.StreetName, normalizedStreet)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatch = segment
		}
	}

	response := models.AddressSearchResponse{
		StreetSegment: bestMatch,
		Team:          bestMatch.Team,
		UBS:           bestMatch.Team.UBS,
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// Helper function to calculate string similarity
func calculateSimilarity(s1, s2 string) float64 {
	// Count matching characters
	matches := 0
	s1Runes := []rune(s1)
	s2Runes := []rune(s2)

	for i := 0; i < len(s1Runes) && i < len(s2Runes); i++ {
		if s1Runes[i] == s2Runes[i] {
			matches++
		}
	}

	// Calculate similarity as a ratio
	maxLen := float64(max(len(s1Runes), len(s2Runes)))
	if maxLen == 0 {
		return 0
	}
	return float64(matches) / maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
