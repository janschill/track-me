package server

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/janschill/track-me/internal/db"
)

type httpServer struct {
	EventStore *EventStore
}

type EventStore struct {
	mu     sync.Mutex
	db     *sql.DB
	events []db.Event
	cache  map[int]db.Event
}

func (c *EventStore) prepareAndSave(payload GarminOutboundPayload) error {
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
			Altitude:    pEvent.Point.Altitude,
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

func reduce(events []db.Event) []db.Event {
	return events
}

func fillCache() {

}

func HttpServer(addr string) *http.Server {
	db, err := db.InitializeDB("./data/trips.db")
	if err != nil {
		log.Fatal(err)
	}
	server := &httpServer{
		EventStore: &EventStore{db: db},
	}
	fs := http.FileServer(http.Dir("assets/"))
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/", server.handleIndex).Methods("GET")
	router.HandleFunc("/events", server.handleEvents).Methods("GET")
	router.HandleFunc("/garmin-outbound", server.handleGarminOutbound).Methods("POST")

	return &http.Server{
		Addr:    ":" + addr,
		Handler: router,
	}
}
