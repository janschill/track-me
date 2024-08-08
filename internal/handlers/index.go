package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/service"
	"github.com/janschill/track-me/internal/utils"
)

type IndexHandler struct {
	repo *repository.Repository
}

type IndexPageData struct {
	Messages       []repository.Message
	LastEvent      repository.Event
	Ride           service.Ride
	Days           []repository.Day
	DaysEventsJSON template.JS
}

func NewIndexHandler(repo *repository.Repository) *IndexHandler {
	return &IndexHandler{repo: repo}
}

func (h *IndexHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"wroteOnTime": utils.WroteOnTime,
		"onDay":       utils.OnDay,
		"time":        utils.FormatTime,
		"oneDecimal":  utils.OneDecimal,
		"inKm":        utils.InKm,
		"addOne":      func(i int) int { return i + 1 },
	}
	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("web/templates/layout.html", "web/templates/index.html"))

	messages, err := h.repo.Messages.All()
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving messages: %v", err)
		return
	}
	log.Printf("Retrieved %d messages", len(messages))

	events, err := h.repo.Events.All()
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving events: %v", err)
		return
	}
	log.Printf("Retrieved %d events", len(events))

	days, err := h.repo.Days.All()
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving days: %v", err)
		return
	}

	var combinedPoints []string

	for _, day := range days {
		trimmedPoints := strings.Trim(day.Points, "[]")
		combinedPoints = append(combinedPoints, trimmedPoints)
	}

	daysEventsJSON := "[" + strings.Join(combinedPoints, ",") + "]"
	var lastEvent repository.Event
	if len(events) > 0 {
		lastEvent = events[len(events)-1]
	}

	eventsJSON, _ := json.Marshal(events)

	data := IndexPageData{
		Messages:       messages,
		LastEvent:      lastEvent,
		Ride:           ride,
		Days:           days,
		EventsJSON: template.JS(eventsJSON),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
