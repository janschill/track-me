package server

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/janschill/track-me/internal/db"
	"github.com/joho/godotenv"
)

type httpServer struct {
	Env *Env
}

type Env struct {
	mu     sync.Mutex
	db     *sql.DB
	events []db.Event
	cache  map[int]db.Event
}

func authorize(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		expectedToken := os.Getenv("AUTHORIZATION_TOKEN")
		log.Printf("token: %v", token)

		if token == expectedToken {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		next.ServeHTTP(w, r)
	})
}

func newHTTPHandler(server *httpServer) http.Handler {
	mux := http.NewServeMux()
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs)) // Correctly add to mux

	mux.Handle("/", sentryHandler.Handle(http.HandlerFunc(server.handleIndex)))
	mux.Handle("/events", sentryHandler.Handle(http.HandlerFunc(server.handleEvents)))
	mux.Handle("/messages", sentryHandler.Handle(http.HandlerFunc(server.handleMessages)))
	mux.Handle("/garmin-outbound", sentryHandler.Handle(http.HandlerFunc(server.handleGarminOutbound)))

	return mux
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}
}

func HttpServer(addr string, ctx context.Context) *http.Server {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}
	db, err := db.InitializeDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	server := &httpServer{
		Env: &Env{db: db},
	}

	return &http.Server{
		Addr:         ":" + addr,
		Handler:      newHTTPHandler(server),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
