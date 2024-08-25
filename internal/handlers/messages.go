package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/janschill/track-me/internal/repository"
	garmin "github.com/janschill/track-me/pkg/garmin"
)

type MessageHandler struct {
	repo   *repository.Repository
	client *garmin.Client
}

func NewMessageHandler(repo *repository.Repository, client *garmin.Client) *MessageHandler {
	return &MessageHandler{
		repo:   repo,
		client: client,
	}
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

	message := r.FormValue("message")
	name := r.FormValue("name")
	sender := r.FormValue("email")

	log.Printf("Message received: %s\n", message)
	log.Printf("Name: %s\n", name)

	if message == "" || name == "" {
		http.Error(w, "Name or message cannot be blank", http.StatusBadRequest)
		return
	}

	sentToGarmin, err := strconv.ParseBool(r.FormValue("sentToGarmin"))
	if err != nil {
		sentToGarmin = false
	}
	log.Printf("Sent to Garmin: %v\n", sentToGarmin)

	if sentToGarmin && r.FormValue("email") == "" {
		http.Error(w, "Email cannot be blank when sending to Garmin", http.StatusBadRequest)
		return
	}

	time := time.Now().Unix()

	m := repository.Message{
		TripID:       1,
		Message:      message,
		Name:         name,
		TimeStamp:    time,
		SentToGarmin: sentToGarmin,
		FromGarmin:   false,
	}

	if m.SentToGarmin {
		err = h.client.SendMessage(sender, message)
		if err != nil {
			log.Printf("Sent to Garmin failed %v", err)
		}
	}

	err = h.repo.Messages.Create(m)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"message":      m.Message,
		"name":         m.Name,
		"timeStamp":    strconv.FormatInt(m.TimeStamp, 10),
		"sentToGarmin": strconv.FormatBool(m.SentToGarmin),
		"fromGamin":    strconv.FormatBool(m.FromGarmin),
	})
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
