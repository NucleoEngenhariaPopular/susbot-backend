package handlers

import (
	"gateway/internal/clients"
	"gateway/internal/config"
	"gateway/internal/models"
	"gateway/internal/services"
	"gateway/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	userClient         *clients.UserClient
	addressClient      *clients.AddressClient
	conversationClient *clients.ConversationClient
}

// Cria o handler para as rotas e inicializa os clients
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		userClient:         clients.NewUserClient(cfg),
		addressClient:      clients.NewAddressClient(cfg),
		conversationClient: clients.NewConversationClient(cfg),
	}
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received message from WhatsApp")

	// Parse na mensagem da Twilio
	twilioMessage, err := utils.ParseTwilioRequest(r)
	if err != nil {
		log.Printf("Error parsing Twilio request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Normalizando o numero (remover o prefix do WhatsApp)
	phoneNumber := twilioMessage.From
	phoneNumber = strings.TrimPrefix(phoneNumber, "whatsapp:")

	// Deixando nome de usuario padrao se nao tiver
	userName := twilioMessage.ProfileName
	if userName == "" {
		userName = phoneNumber
	}

	// Create user data with available information
	userData := models.UserData{
		User: models.User{
			Name:         userName,
			CPF:          "PENDENTE",
			DateOfBirth:  time.Now(),
			PhoneNumber:  phoneNumber,
			StreetName:   "Pendente",
			StreetNumber: "S/N",
			Complement:   "",
			Neighborhood: "Pendente",
			City:         "São Paulo",
			State:        "SP",
			CEP:          "00000000",
		},
		Msg: twilioMessage,
	}

	log.Printf("Attempting to save user data: %+v", userData)

	if err := h.userClient.SaveUser(userData); err != nil {
		log.Printf("Error saving user: %v", err)
	} else {
		log.Printf("Sucessfully saved user data")
	}

	// Criar uma mensagem para salvar na api-conversation
	userMessage := models.Message{
		UserID:    twilioMessage.AccountSid,
		Sender:    twilioMessage.AccountSid,
		Text:      twilioMessage.Body,
		Timestamp: time.Now(),
	}

	// Salvando a mensagem na conversa
	if err := h.conversationClient.SaveMessage(userMessage); err != nil {
		log.Printf("Error saving message: %v", err)
	}

	// Enviando a mensagem para o Botkit
	reply, err := services.SendToBotkit(*twilioMessage)
	if err != nil {
		log.Printf("Error sending message to BotKit: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Salvando as respostas do botkit
	for _, msg := range reply {
		botMessage := models.Message{
			UserID:    twilioMessage.AccountSid,
			Sender:    "BotKit",
			Text:      msg,
			Timestamp: time.Now(),
		}
		if err := h.conversationClient.SaveMessage(botMessage); err != nil {
			log.Printf("Error saving bot message: %v", err)
		}
	}

	// Enviando as reposta de volta ao usuário
	if err := services.RespondToUser(w, reply); err != nil {
		log.Printf("Error sending response to user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
