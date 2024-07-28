package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/janschill/track-me/internal/db"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://your-oauth-provider.com/token",
		},
	}
	staticToken = "your-static-token"
)

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

func reduce(events []db.Event) []db.Event {
	return events
}

func fillCache() {

}

func authorize(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("token: %v", token)
		if token == staticToken {
			log.Println("")
			next.ServeHTTP(w, r)
			return
		}

		tokenSource := oauthConfig.TokenSource(context.TODO(), &oauth2.Token{AccessToken: token})
		_, err := tokenSource.Token()
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func HttpServer(addr string) *http.Server {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}

	db, err := db.InitializeDB(dbPath)
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
	router.HandleFunc("/messages", server.handleMessages).Methods("POST")
	router.HandleFunc("/garmin-outbound", authorize(http.HandlerFunc(server.handleGarminOutbound))).Methods("POST")

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	loggedRouter := Tracing(nextRequestID)(Logging(logger)(router))

	return &http.Server{
		Addr:     ":" + addr,
		Handler:  loggedRouter,
		ErrorLog: logger,
	}
}
