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

func HandleUBS(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUBS(w, r)
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/ubs/")
		if path == "" {
			listUBS(w, r)
		} else {
			getUBS(w, r)
		}
	case http.MethodPut:
		updateUBS(w, r)
	case http.MethodDelete:
		deleteUBS(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func createUBS(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUBSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	ubs := models.UBS{
		Name:    req.Name,
		Address: req.Address,
		City:    strings.ToUpper(req.City),
		State:   strings.ToUpper(req.State),
		CEP:     utils.NormalizeCEP(req.CEP),
	}

	if err := database.GetDB().Create(&ubs).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create UBS")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    ubs,
	})
}

func listUBS(w http.ResponseWriter, _ *http.Request) {
	var ubsList []models.UBS
	if err := database.GetDB().Find(&ubsList).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch UBS list")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    ubsList,
	})
}

func getUBS(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/ubs/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UBS ID")
		return
	}

	var ubs models.UBS
	if err := database.GetDB().Preload("Teams").First(&ubs, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "UBS not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    ubs,
	})
}

func updateUBS(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/ubs/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UBS ID")
		return
	}

	var req models.CreateUBSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var ubs models.UBS
	if err := database.GetDB().First(&ubs, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "UBS not found")
		return
	}

	ubs.Name = req.Name
	ubs.Address = req.Address
	ubs.City = strings.ToUpper(req.City)
	ubs.State = strings.ToUpper(req.State)
	ubs.CEP = utils.NormalizeCEP(req.CEP)

	if err := database.GetDB().Save(&ubs).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update UBS")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    ubs,
	})
}

func deleteUBS(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/ubs/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UBS ID")
		return
	}

	var ubs models.UBS
	if err := database.GetDB().First(&ubs, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "UBS not found")
		return
	}

	if err := database.GetDB().Delete(&ubs).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete UBS")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "UBS successfully deleted",
	})
}
