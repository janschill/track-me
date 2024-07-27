package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type IPCOutbound struct {
	Version string  `json:"Version"`
	Events  []Event `json:"Events"`
}

type Event struct {
	Imei        string    `json:"imei"`
	MessageCode int       `json:"messageCode"`
	FreeText    string    `json:"freeText"`
	TimeStamp   int64     `json:"timeStamp"`
	Addresses   []Address `json:"addresses"`
	Point       Point     `json:"point"`
	Status      Status    `json:"status"`
}

type Address struct {
	Address string `json:"address"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  int     `json:"altitude"`
	GpsFix    int     `json:"gpsFix"`
	Course    int     `json:"course"`
	Speed     int     `json:"speed"`
}

type Status struct {
	Autonomous     int `json:"autonomous"`
	LowBattery     int `json:"lowBattery"`
	IntervalChange int `json:"intervalChange"`
	ResetDetected  int `json:"resetDetected"`
}

var lastLatitude float64 = 31.3331459241914
var lastLongitude float64 = -108.530207744702
var lastAltitude int = 1421

func generateMockEvent() Event {
	latitudeDeviation := randomFloatInRange(-0.0001, 0.0001)
	longitudeDeviation := randomFloatInRange(-0.0001, 0.0001)
	altitudeDeviation := rand.Intn(3) - 1 // -1, 0, or 1

	// Update last known coordinates
	lastLatitude += latitudeDeviation
	lastLongitude += longitudeDeviation
	lastAltitude += altitudeDeviation

	return Event{
		Imei:        "fake-imei",
		MessageCode: rand.Intn(100),
		FreeText:    "Example text",
		TimeStamp:   time.Now().Unix(),
		Addresses:   []Address{{Address: "example@example.com"}},
		Point: Point{
			Latitude:  lastLatitude,
			Longitude: lastLongitude,
			Altitude:  lastAltitude,
			GpsFix:    rand.Intn(5),
			Course:    rand.Intn(360),
			Speed:     rand.Intn(120),
		},
		Status: Status{
			Autonomous:     0,
			LowBattery:     0,
			IntervalChange: 0,
			ResetDetected:  0,
		},
	}
}

func randomFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

var eventQueue []Event

func sendMockEvents(events []Event, url string) {
	ipcOutbound := IPCOutbound{
		Version: "1.0",
		Events:  events,
	}

	jsonData, err := json.Marshal(ipcOutbound)
	if err != nil {
		fmt.Println("Error marshalling event:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer fooo")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	eventQueue = nil
	fmt.Println("POST request sent. Status Code:", resp.StatusCode)
}

func main() {
	targetURL := "http://localhost:8080/garmin-outbound"
	ticker := time.NewTicker(2 * time.Second)

	for range ticker.C {
		mockEvent := generateMockEvent()
		eventQueue = append(eventQueue, mockEvent)

		// Randomly decide to simulate missed communication
		if rand.Intn(10) < 8 { // 80% chance to send
			sendMockEvents(eventQueue, targetURL)
		} else {
			fmt.Println("Simulating missed communication...")
		}
	}
}
