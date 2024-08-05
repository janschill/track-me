package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/utils"
)

type GarminHandler struct {
	repo *repository.Repository
}

func NewGarminHandler(repo *repository.Repository) *GarminHandler {
	return &GarminHandler{repo: repo}
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

func (h *GarminHandler) prepareAndSave(payload GarminOutboundPayload) error {
	for _, pEvent := range payload.Events {
		event := repository.Event{
			TripID:      1,
			Imei:        pEvent.Imei,
			MessageCode: pEvent.MessageCode,
			FreeText:    pEvent.FreeText,
			TimeStamp:   pEvent.TimeStamp / 1000, // comes in millisecond format
			Addresses:   make([]repository.Address, len(pEvent.Addresses)),
			Latitude:    pEvent.Point.Latitude,
			Longitude:   pEvent.Point.Longitude,
			Altitude:    pEvent.Point.Altitude,
			GpsFix:      pEvent.Point.GpsFix,
			Course:      pEvent.Point.Course,
			Speed:       pEvent.Point.Speed,
			Status: repository.Status{
				Autonomous:     pEvent.Status.Autonomous,
				LowBattery:     pEvent.Status.LowBattery,
				IntervalChange: pEvent.Status.IntervalChange,
				ResetDetected:  pEvent.Status.ResetDetected,
			},
		}

		for i, addr := range pEvent.Addresses {
			event.Addresses[i] = repository.Address{Address: addr.Address}
		}

		if utils.HasMessage(event) {
			message := repository.Message{
				TripID:     1,
				Message:    event.FreeText,
				Name:       "Automated Message",
				TimeStamp:  event.TimeStamp,
				FromGarmin: true,
			}

			if err := h.repo.Messages.Create(message); err != nil {
				log.Printf("Failed to save message from event %v", event.ID)
			}
		}

		if err := h.repo.Events.Create(event); err != nil {
			return err
		}
	}
	return nil
}

func (h *GarminHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
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
	h.prepareAndSave(payload)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Payload received successfully."))
}
