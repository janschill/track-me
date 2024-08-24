package garmin

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type ProcessPayloadFunc func(payload OutboundPayload) error

type OutboundHandler struct {
	processPayload ProcessPayloadFunc
}

func NewOutboundHandler(processPayload ProcessPayloadFunc) *OutboundHandler {
	return &OutboundHandler{processPayload: processPayload}
}

type OutboundPayload struct {
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
			Altitude  float64 `json:"altitude"`
			GpsFix    int     `json:"gpsFix"`
			Course    float64 `json:"course"`
			Speed     float64 `json:"speed"`
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

func (h *OutboundHandler) CreateOutboundEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var payload OutboundPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusInternalServerError)
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.CaptureMessage(err.Error())
			log.Printf("Error. %v\n", err)
		}
		return
	}

	log.Printf("Outbound payload received. %v event(s)\n", len(payload.Events))
	// h.prepareAndSave(payload)
	if err := h.processPayload(payload); err != nil {
		http.Error(w, "Error processing payload", http.StatusInternalServerError)
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
				hub.CaptureMessage(err.Error())
				log.Printf("Error. %v\n", err)
		}
		return
}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payload received successfully."))
}
