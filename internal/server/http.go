package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/janschill/track-me/internal/db"
)

type httpServer struct {
	Events *EventStore
}

type EventStore struct {
	mu     sync.Mutex
	events []db.Event
}


func (c *EventStore) saveEventToStorage(event db.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	event.ID = int64(len(c.events)) + 1
	c.events = append(c.events, event)

	return nil
}

func (c *EventStore) save(payload GarminOutboundPayload) error {
	for _, pEvent := range payload.Events {
		event := db.Event{
			TripID:      1,
			Imei:        pEvent.Imei,
			MessageCode: pEvent.MessageCode,
			FreeText:    pEvent.FreeText,
			TimeStamp:   pEvent.TimeStamp,
			Addresses:   make([]db.Address, len(pEvent.Addresses)),
			Point: db.Point{
				Latitude:  pEvent.Point.Latitude,
				Longitude: pEvent.Point.Longitude,
				Altitude:  pEvent.Point.Altitude,
				GpsFix:    pEvent.Point.GpsFix,
				Course:    pEvent.Point.Course,
				Speed:     pEvent.Point.Speed,
			},
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

		if err := c.saveEventToStorage(event); err != nil {
			return err
		}
	}
	return nil

}

func HttpServer(addr string) *http.Server {
	server := &httpServer{
		Events: &EventStore{},
	}
	fs := http.FileServer(http.Dir("assets/"))
	router := mux.NewRouter()
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", server.handleIndex).Methods("GET")
	router.HandleFunc("/events", server.handleEvents).Methods("GET")
	router.HandleFunc("/garmin-outbound", server.handleGarminOutbound).Methods("POST")

	return &http.Server{
		Addr:    ":" + addr,
		Handler: router,
	}
}
