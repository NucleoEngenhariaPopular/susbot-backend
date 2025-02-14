package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"conversation-api/internal/database"
	"conversation-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandleConversations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		saveMessage(w, r)
	case http.MethodGet:
		getConversation(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func saveMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	collection := database.GetCollection()
	ctx := context.Background()

	// Find the latest conversation for this user
	var conversation models.Conversation
	opts := options.FindOne().SetSort(bson.D{{Key: "start_time", Value: -1}})
	err := collection.FindOne(ctx,
		bson.M{
			"user_id":  msg.UserID,
			"end_time": nil, // Only find conversations that haven't ended
		},
		opts,
	).Decode(&conversation)

	if err != nil {
		// Create new conversation if none exists or last one is closed
		conversation = models.Conversation{
			UserID:    msg.UserID,
			StartTime: time.Now(),
			Messages:  []models.Message{msg},
		}

		result, err := collection.InsertOne(ctx, conversation)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create conversation")
			return
		}

		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			conversation.ID = oid.Hex()
		}
	} else {
		// Add message to existing conversation
		objectID, err := primitive.ObjectIDFromHex(conversation.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Invalid conversation ID")
			return
		}

		// Update the existing conversation with the new message
		update := bson.M{
			"$push": bson.M{"messages": msg},
		}

		_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to add message")
			return
		}

		// Get the updated conversation
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&conversation)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to fetch updated conversation")
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    conversation,
	})
}

func getConversation(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		respondWithError(w, http.StatusBadRequest, "Invalid conversation ID")
		return
	}

	id := parts[2]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid conversation ID format")
		return
	}

	var conversation models.Conversation
	err = database.GetCollection().FindOne(context.Background(), bson.M{"_id": objID}).Decode(&conversation)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    conversation,
	})
}

// Helper functions
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
