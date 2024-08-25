package garmin

import (
	"encoding/json"
	"log"
	"net/http"
)

type ProcessPayloadFunc func(payload OutboundPayload) error

type OutboundHandler struct {
	processPayload ProcessPayloadFunc
}

func NewOutboundHandler(processPayload ProcessPayloadFunc) *OutboundHandler {
	return &OutboundHandler{processPayload: processPayload}
}

type Address struct {
	Address string `json:"address"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	GpsFix    int     `json:"gpsFix"`
	Course    float64 `json:"course"`
	Speed     float64 `json:"speed"`
}

type Status struct {
	Autonomous     int `json:"autonomous"`
	LowBattery     int `json:"lowBattery"`
	IntervalChange int `json:"intervalChange"`
	ResetDetected  int `json:"resetDetected"`
}

type Event struct {
	Imei              string    `json:"imei"`
	MessageCode       int       `json:"messageCode"`
	FreeText          string    `json:"freeText"`
	TimeStamp         int64     `json:"timeStamp"`
	PingbackReceived  int64     `json:"pingbackReceived"`
	PingbackResponded int64     `json:"pingbackResponded"`
	Addresses         []Address `json:"addresses"`
	Point             Point     `json:"point"`
	Status            Status    `json:"status"`
	Payload           string    `json:"payload"`
}

type OutboundPayload struct {
	Version string  `json:"Version"`
	Events  []Event `json:"Events"`
}

func (h *OutboundHandler) CreateOutboundEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var payload OutboundPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)

		return
	}

	log.Printf("Outbound payload received. %v event(s)\n", len(payload.Events))

	if err := h.processPayload(payload); err != nil {
		http.Error(w, "Error processing payload", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Payload received successfully."))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
