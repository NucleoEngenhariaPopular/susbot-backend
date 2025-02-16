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
		if r.URL.Path == "/conversations/" {
			getAllConversations(w, r)
			return
		}
		getConversation(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func getAllConversations(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	collection := database.GetCollection()

	// Define the projection to only get ID and UserID
	projection := bson.D{
		{Key: "_id", Value: 1},
		{Key: "user_id", Value: 1},
		{Key: "start_time", Value: 1},
		{Key: "end_time", Value: 1},
	}

	// Configure options for sorting by start_time in descending order
	findOptions := options.Find().
		SetProjection(projection).
		SetSort(bson.D{{Key: "start_time", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch conversations")
		return
	}
	defer cursor.Close(ctx)

	type ConversationSummary struct {
		ID        primitive.ObjectID `json:"id" bson:"_id"`
		UserID    string             `json:"user_id" bson:"user_id"`
		StartTime time.Time          `json:"start_time" bson:"start_time"`
		EndTime   *time.Time         `json:"end_time,omitempty" bson:"end_time"`
		Status    string             `json:"status"`
	}

	var conversations []ConversationSummary
	for cursor.Next(ctx) {
		var conv ConversationSummary
		if err := cursor.Decode(&conv); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error decoding conversation")
			return
		}

		// Determine conversation status
		if conv.EndTime == nil {
			conv.Status = "active"
		} else {
			conv.Status = "closed"
		}

		conversations = append(conversations, conv)
	}

	if err := cursor.Err(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error iterating conversations")
		return
	}

	// Create response with string IDs
	type ConversationResponse struct {
		ID        string     `json:"id"`
		UserID    string     `json:"user_id"`
		StartTime time.Time  `json:"start_time"`
		EndTime   *time.Time `json:"end_time,omitempty"`
		Status    string     `json:"status"`
	}

	response := make([]ConversationResponse, len(conversations))
	for i, conv := range conversations {
		response[i] = ConversationResponse{
			ID:        conv.ID.Hex(),
			UserID:    conv.UserID,
			StartTime: conv.StartTime,
			EndTime:   conv.EndTime,
			Status:    conv.Status,
		}
	}

	respondWithJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data: struct {
			TotalConversations int                    `json:"total_conversations"`
			Conversations      []ConversationResponse `json:"conversations"`
		}{
			TotalConversations: len(conversations),
			Conversations:      response,
		},
	})
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
