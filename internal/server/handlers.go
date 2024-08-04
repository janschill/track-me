package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"html/template"

	"github.com/getsentry/sentry-go"
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
		Imei              string `json:"imei"`
		MessageCode       int    `json:"messageCode"`
		FreeText          string `json:"freeText"`
		TimeStamp         int64  `json:"timeStamp"`
		PingbackReceived  int64  `json:"pingbackReceived"`
		PingbackResponded int64  `json:"pingbackResponded"`
		Addresses         []struct {
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
		Payload string `json:"payload"`
	} `json:"Events"`
}

func (c *Env) prepareAndSave(payload GarminOutboundPayload) error {
	for _, pEvent := range payload.Events {
		event := db.Event{
			TripID:      1,
			Imei:        pEvent.Imei,
			MessageCode: pEvent.MessageCode,
			FreeText:    pEvent.FreeText,
			TimeStamp:   pEvent.TimeStamp,
			Addresses:   make([]db.Address, len(pEvent.Addresses)),
			Latitude:    pEvent.Point.Latitude,
			Longitude:   pEvent.Point.Longitude,
			Altitude:    int64(pEvent.Point.Altitude),
			GpsFix:      pEvent.Point.GpsFix,
			Course:      pEvent.Point.Course,
			Speed:       pEvent.Point.Speed,
			Status: db.Status{
				Autonomous:     pEvent.Status.Autonomous,
				LowBattery:     pEvent.Status.LowBattery,
				IntervalChange: pEvent.Status.IntervalChange,
				ResetDetected:  pEvent.Status.ResetDetected,
			},
		}

		for i, addr := range pEvent.Addresses {
			event.Addresses[i] = db.Address{Address: addr.Address}
		}

		if err := event.Save(c.db); err != nil {
			return err
		}
	}
	return nil
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
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureMessage(err.Error())
			log.Printf("Error. %v\n", err)
		}
		return
	}

	log.Printf("GarminOutbound payload received. %v event(s)\n", len(payload.Events))
	s.Env.prepareAndSave(payload)

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

	if r.FormValue("message") == "" || r.FormValue("name") == "" {
		http.Error(w, "Name or message cannot be blank", http.StatusBadRequest)
		return
	}

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

	message.Save(s.Env.db)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":      message.Message,
		"name":         message.Name,
		"timeStamp":    strconv.FormatInt(message.TimeStamp, 10),
		"sentToGarmin": strconv.FormatBool(message.SentToGarmin),
	})
}

func (s *httpServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	s.Env.mu.Lock()
	defer s.Env.mu.Unlock()

	if err := json.NewEncoder(w).Encode(s.Env.events); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

type IndexPageData struct {
	Messages       []db.Message
	LastEvent      db.Event
	Ride           Ride
	Days           []db.Day
	DaysEventsJSON template.JS
}

func wroteOnTime(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("on 02 January at 15:04")
}

func (s *httpServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"wroteOnTime": wroteOnTime,
		"onDay":       onDay,
		"time":        formatTime,
		"oneDecimal":  oneDecimal,
		"inKm":        inKm,
		"addOne":      func(i int) int { return i + 1 },
	}
	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/index.html"))

	messages, err := db.GetAllMessages(s.Env.db)
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving messages: %v", err)
		return
	}
	log.Printf("Retrieved %d messages", len(messages))

	events, err := db.GetAllEvents(s.Env.db)
	if err != nil {
		http.Error(w, "An unexpected error happened.", http.StatusBadGateway)
		log.Printf("Error retrieving events: %v", err)
		return
	}
	log.Printf("Retrieved %d events", len(events))

	days, err := db.GetAllDays(s.Env.db)
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

	var (
		lastEvent            db.Event
		currentSpeed         float64
		isMoving             bool
		dist                 int64
		gain                 int64
		loss                 int64
		movingTime           int64
		movingTimeFormatted  string
		restingTimeFormatted string
	)

	if len(events) > 1 {
		lastEvent = events[len(events)-1]
		currentSpeed, isMoving = movement(events)
		dist = distance(days, events)
		gain, loss = elevation(days, events)
		movingTime = timeMoving(days, events)
		movingTimeFormatted = formatTime(movingTime)
		restingTimeFormatted = formatTime(restingTime(len(days), movingTime))
	}

	data := IndexPageData{
		Messages:  messages,
		LastEvent: lastEvent,
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
		Days:           days,
		DaysEventsJSON: template.JS(daysEventsJSON),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}
