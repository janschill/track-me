package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/janschill/track-me/internal/db"
)

type httpServer struct {
  Events *Events
}

type Events struct {
  events []db.Event
}

type Payload struct {
  Version string `json:"Version"`
  Events  []struct {
    Imei        string `json:"imei"`
    MessageCode int    `json:"messageCode"`
    FreeText    string `json:"freeText"`
    TimeStamp   int64  `json:"timeStamp"`
    Addresses   []struct {
      Address string `json:"address"`
    } `json:"addresses"`
    Point struct {
      Latitude  float64 `json:"latitude"`
      Longitude float64 `json:"longitude"`
      Altitude  int     `json:"altitude"`
      GpsFix    int     `json:"gpsFix"`
      Course    int     `json:"course"`
      Speed     int     `json:"speed"`
    } `json:"point"`
    Status struct {
      Autonomous     int `json:"autonomous"`
      LowBattery     int `json:"lowBattery"`
      IntervalChange int `json:"intervalChange"`
      ResetDetected  int `json:"resetDetected"`
    } `json:"status"`
  } `json:"Events"`
}

func saveEventToStorage(c *Events, event db.Event) error {
  event.ID = int64(len(c.events)) + 1
  c.events = append(c.events, event)

  return nil
}

func (c *Events) save(payload Payload) error {
  for _, pEvent := range payload.Events {
    event := db.Event{
      TripID: 		 1,
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

    err := saveEventToStorage(c, event)
    if err != nil {
      return err
    }
  }
  return nil

}

func (s *httpServer) handleGarminOutbound(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
    return
  }

  var payload Payload
  err := json.NewDecoder(r.Body).Decode(&payload)
  if err != nil {
    http.Error(w, "Error parsing request body", http.StatusInternalServerError)
    return
  }

  s.Events.save(payload)

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Payload received successfully."))
}

func (c *httpServer) handleEvents(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(c.Events.events)
}

func HttpServer(addr string) *http.Server {
  server := &httpServer{
    Events: &Events{},
  }
  router := mux.NewRouter()
  router.HandleFunc("/garmin-outbound", server.handleGarminOutbound).Methods("POST")
  router.HandleFunc("/events", server.handleEvents).Methods("GET")

  return &http.Server{
    Addr:    ":" + addr,
    Handler: router,
  }
}
