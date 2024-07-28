package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"html/template"

	"github.com/janschill/track-me/internal/db"
)

type Ride struct {
	IsMoving      bool
	LastPing      int64
	Distance      int64
	Progress      float64
	ElevationGain int64
	ElevationLoss int64
	MovingTime    string
	RestingTime   string
	ElapsedDays   int
	RemainingDays int
	CurrentSpeed  float64
}

type GarminOutboundPayload struct {
	Version string `json:"Version"`
	Events  []struct {
		Imei        string `json:"imei"`
		MessageCode int    `json:"messageCode"`
		FreeText    string `json:"freeText"`
		TimeStamp   int64  `json:"timeStamp"`
		Addresses   []struct {
			Address string `json:"address"`
		} `json:"addresses"`
		Point struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Altitude  int     `json:"altitude"`
			GpsFix    int     `json:"gpsFix"`
			Course    int     `json:"course"`
			Speed     int     `json:"speed"`
		} `json:"point"`
		Status struct {
			Autonomous     int `json:"autonomous"`
			LowBattery     int `json:"lowBattery"`
			IntervalChange int `json:"intervalChange"`
			ResetDetected  int `json:"resetDetected"`
		} `json:"status"`
	} `json:"Events"`
}

func (s *httpServer) handleGarminOutbound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var payload GarminOutboundPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		return
	}

	log.Printf("GarminOutbound payload received. %v event(s)\n", len(payload.Events))
	s.EventStore.prepareAndSave(payload)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payload received successfully."))
}

func (s *httpServer) handleMessages(w http.ResponseWriter, r *http.Request) {
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
	sentToGarmin, err := strconv.ParseBool(r.FormValue("sentToGarmin"))
	if err != nil {
		sentToGarmin = false
	}

	message := db.Message{
		TripID:       1,
		Message:      r.FormValue("message"),
		Name:         r.FormValue("name"),
		TimeStamp:    time.Now().Unix(),
		SentToGarmin: sentToGarmin,
	}

	if message.SentToGarmin {
		log.Printf("Sending message to Garmin: %s\n", message.Message)
	}

	message.Save(s.EventStore.db)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received successfully."))
}

func (s *httpServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	s.EventStore.mu.Lock()
	defer s.EventStore.mu.Unlock()

	if err := json.NewEncoder(w).Encode(s.EventStore.events); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

type IndexPageData struct {
	Messages   []db.Message
	Events     []db.Event
	LastEvent  db.Event
	EventsJSON template.JS
	Ride       Ride
	Days       []db.Day
}

func wroteOnTime(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("on 02 January at 15:04")
}

func onDay(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("02 January")
}

func (s *httpServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"wroteOnTime": wroteOnTime,
		"onDay": onDay,
		"time": formatTime,
		"oneDecimal": oneDecimal,
		"inKm": inKm,
		"addOne": func(i int) int { return i + 1 },
	}
	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/index.html"))

	messages, err := db.GetAllMessages(s.EventStore.db)
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving messages: %v", err)
		return
	}
	log.Printf("Retrieved %d messages", len(messages))

	events, err := db.GetAllEvents(s.EventStore.db)
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving events: %v", err)
		return
	}
	log.Printf("Retrieved %d events", len(events))

	days, err := db.GetAllDays(s.EventStore.db)
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving days: %v", err)
		return
	}

	lastEvent := events[len(events)-1]
	currentSpeed, isMoving := isMoving(events)
	dist := distance(days, events)
	gain, loss := elevation(days, events)
	movingTime := movingTime(days, events)
	movingTimeFormatted := formatTime(movingTime)
	restingTimeFormatted := formatTime(restingTime(len(days), movingTime))

	events = db.Rdp(events, 0.0002) // roughly 1500 -> 321
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error marshalling events: %v", err)
		return
	}

	data := IndexPageData{
		Messages:   messages,
		Events:     events,
		LastEvent:  lastEvent,
		EventsJSON: template.JS(eventsJSON),
		Ride: Ride{
			IsMoving:      isMoving,
			LastPing:      lastEvent.TimeStamp,
			Distance:      dist,
			Progress:      progress(dist),
			CurrentSpeed:  currentSpeed,
			ElevationGain: gain,
			ElevationLoss: loss,
			MovingTime:    movingTimeFormatted,
			RestingTime:   restingTimeFormatted,
			ElapsedDays:   len(days),
			RemainingDays: 30 - len(days),
		},
		Days: days,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
