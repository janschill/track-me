package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/janschill/track-me/internal/repository"
)

type MessageHandler struct {
	repo *repository.Repository
}

func NewMessageHandler(repo *repository.Repository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	log.Printf("Message received: %s\n", r.FormValue("message"))
	log.Printf("Name: %s\n", r.FormValue("name"))
	log.Printf("Sent to Garmin: %v\n", r.FormValue("sentToGarmin"))

	if r.FormValue("message") == "" || r.FormValue("name") == "" {
		http.Error(w, "Name or message cannot be blank", http.StatusBadRequest)
		return
	}

	sentToGarmin, err := strconv.ParseBool(r.FormValue("sentToGarmin"))
	if err != nil {
		sentToGarmin = false
	}

	message := repository.Message{
			TripID:       1,
			Message:      r.FormValue("message"),
			Name:         r.FormValue("name"),
			TimeStamp:    time.Now().Unix(),
			SentToGarmin: sentToGarmin,
			FromGarmin: false,
	}

	if message.SentToGarmin {
		log.Printf("Sending message to Garmin: %s\n", message.Message)
	}

	h.repo.Messages.Create(message)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":      message.Message,
		"name":         message.Name,
		"timeStamp":    strconv.FormatInt(message.TimeStamp, 10),
		"sentToGarmin": strconv.FormatBool(message.SentToGarmin),
		"fromGamin":    strconv.FormatBool(message.FromGarmin),
	})
}
