package server

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/janschill/track-me/internal/db"
)

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

	log.Printf("Payload received: %+v\n", payload)
	s.Events.save(payload)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payload received successfully."))
}

func (s *httpServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	s.Events.mu.Lock()
	defer s.Events.mu.Unlock()

	if err := json.NewEncoder(w).Encode(s.Events.events); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

type IndexPageData struct {
	Events []db.Event
}

func (s *httpServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.Events.mu.Lock()
	defer s.Events.mu.Unlock()

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/index.html"))

	data := IndexPageData{
		Events: s.Events.events,
	}

	tmpl.Execute(w, data)
}
