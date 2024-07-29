package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/janschill/track-me/internal/db"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
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

func newHTTPHandler(server *httpServer) http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	fs := http.FileServer(http.Dir("assets/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs)) // Correctly add to mux
	// Register handlers.
	handleFunc("/", server.handleIndex)
	handleFunc("/events", server.handleEvents)
	handleFunc("/messages", server.handleMessages)
	handleFunc("/garmin-outbound", authorize(http.HandlerFunc(server.handleGarminOutbound)))

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func HttpServer(addr string) *http.Server {
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
		Addr:     ":" + addr,
		Handler:  newHTTPHandler(server),
	}
}
