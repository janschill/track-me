package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/janschill/track-me/internal/repository"
)

type KudosHandler struct {
	repo *repository.Repository
}

func NewKudosHandler(repo *repository.Repository) *KudosHandler {
	return &KudosHandler{
		repo: repo,
	}
}

func (h *KudosHandler) CreateKudos(w http.ResponseWriter, r *http.Request) {
	log.Print("kudos")
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		log.Print("method")
		return
	}

	var requestData struct {
		Day string `json:"day"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Day == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Print("requestdata")
		return
	}

	err = h.repo.Kudos.Increment(requestData.Day)
	if err != nil {
		http.Error(w, "Failed to update kudos", http.StatusInternalServerError)
		log.Print("repo")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
