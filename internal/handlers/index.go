package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/service"
	"github.com/janschill/track-me/internal/utils"
)

type IndexHandler struct {
	repo       *repository.Repository
	dayService *service.DayService
}

type IndexPageData struct {
	Messages   []repository.Message
	LastEvent  repository.Event
	Ride       service.Ride
	Days       []service.Day
	EventsJSON template.JS
}

func NewIndexHandler(repo *repository.Repository, service *service.DayService) *IndexHandler {
	return &IndexHandler{
		repo:       repo,
		dayService: service,
	}
}

func (h *IndexHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"wroteOnTime":     utils.WroteOnTime,
		"onDay":           utils.OnDay,
		"onDayFromString": utils.OnDayFromString,
		"time":            utils.FormatTime,
		"oneDecimal":      utils.OneDecimal,
		"inKm":            utils.InKm,
		"addOne":          func(i int) int { return i + 1 },
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

	days, ride := h.dayService.GetDays(events)

	var lastEvent repository.Event
	if len(events) > 0 {
		lastEvent = events[len(events)-1]
	}

	eventsJSON, _ := json.Marshal(events)

	data := IndexPageData{
		Messages:   messages,
		LastEvent:  lastEvent,
		Ride:       ride,
		Days:       days,
		EventsJSON: template.JS(eventsJSON),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
