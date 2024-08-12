package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/janschill/track-me/internal/clients"
	"github.com/janschill/track-me/internal/config"
	"github.com/janschill/track-me/internal/db"
	"github.com/janschill/track-me/internal/handlers"
	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/service"
)

var conf *config.Config

func newHTTPHandler(repo *repository.Repository, dayService *service.DayService, garminClient *clients.GarminClient) http.Handler {
	mux := http.NewServeMux()
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	fs := http.FileServer(http.Dir("web/assets/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/", sentryHandler.Handle(http.HandlerFunc(handlers.NewIndexHandler(repo, dayService).GetIndex)))
	mux.Handle("/messages", sentryHandler.Handle(http.HandlerFunc(handlers.NewMessageHandler(repo, garminClient).CreateMessage)))
	mux.Handle("/kudos", sentryHandler.Handle(http.HandlerFunc(handlers.NewKudosHandler(repo).CreateKudos)))
	mux.Handle("/garmin-outbound", sentryHandler.Handle(http.HandlerFunc(handlers.NewGarminHandler(repo).CreateOutboundEvent)))
	mux.Handle("/photos", sentryHandler.Handle(http.HandlerFunc(handlers.NewICloudHandler().Photos)))

	return mux
}

func init() {
	var err error
	conf, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Couldnt load config %v", err)
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              conf.SentryDsn,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}
}

func HttpServer(addr string, ctx context.Context) *http.Server {
	if conf.DatabaseURL == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}
	db, err := db.InitializeDB(conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepository(db)
	dayService := service.NewDayService()
	garminClient := clients.NewGarminClient(conf.GarminIPCInbound, conf.GarminDeviceIMEI)

	return &http.Server{
		Addr:         ":" + addr,
		Handler:      newHTTPHandler(repo, dayService, garminClient),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
