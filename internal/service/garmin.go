package service

import (
	"log"

	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/utils"
	"github.com/janschill/track-me/pkg/garmin"
)

type GarminService struct {
	repo *repository.Repository
}

func NewGarminService(repo *repository.Repository) *GarminService {
	return &GarminService{repo: repo}
}


func (s *GarminService) ProcessPayload(payload garmin.OutboundPayload) error {
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

			if err := s.repo.Messages.Create(message); err != nil {
				log.Printf("Failed to save message from event %v", event.ID)
			}
		}

		if err := s.repo.Events.Create(event); err != nil {
			return err
		}
	}
	return nil
}
