package handlers

import (
	"address-api/internal/database"
	"address-api/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func HandleTeams(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createTeam(w, r)
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/teams/")
		if path == "" {
			listTeams(w, r)
		} else {
			getTeam(w, r)
		}
	case http.MethodPut:
		updateTeam(w, r)
	case http.MethodDelete:
		deleteTeam(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Verify if UBS exists
	var ubs models.UBS
	if err := database.GetDB().First(&ubs, req.UBSID).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UBS ID")
		return
	}

	team := models.Team{
		Name:  req.Name,
		UBSID: req.UBSID,
	}

	if err := database.GetDB().Create(&team).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create team")
		return
	}

	// Fetch the complete team data with UBS information
	if err := database.GetDB().Preload("UBS").First(&team, team.ID).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch team data")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    team,
	})
}

func listTeams(w http.ResponseWriter, _ *http.Request) {
	var teams []models.Team
	if err := database.GetDB().Preload("UBS").Find(&teams).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch teams")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    teams,
	})
}

func getTeam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var team models.Team
	if err := database.GetDB().Preload("UBS").First(&team, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Team not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    team,
	})
}

func updateTeam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	var team models.Team
	if err := database.GetDB().First(&team, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Team not found")
		return
	}

	// Verify if new UBS exists
	if err := database.GetDB().First(&models.UBS{}, req.UBSID).Error; err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UBS ID")
		return
	}

	team.Name = req.Name
	team.UBSID = req.UBSID

	if err := database.GetDB().Save(&team).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update team")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    team,
	})
}

func deleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/teams/"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var team models.Team
	if err := database.GetDB().First(&team, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Team not found")
		return
	}

	// Confere se um time possui segmentos de rua associado
	var count int64
	if err := database.GetDB().Model(&models.StreetSegment{}).Where("team_id = ?", id).Count(&count).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check team dependencies")
		return
	}

	if count > 0 {
		respondWithError(w, http.StatusBadRequest, "Cannot delete team with assigned street segments")
		return
	}

	if err := database.GetDB().Delete(&team).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete team")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    "Team successfully deleted",
	})
}
