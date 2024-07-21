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

func generateMockEvent() IPCOutbound {
	return IPCOutbound{
		Version: "1.0",
		Events: []Event{
			{
				Imei:        "100000000000001",
				MessageCode: rand.Intn(100), // Example: Random message code
				FreeText:    "Example text",
				TimeStamp:   time.Now().Unix(),
				Addresses:   []Address{{Address: "example@example.com"}},
				Point: Point{
					Latitude:  randomFloatInRange(-90, 90),
					Longitude: randomFloatInRange(-180, 180),
					Altitude:  rand.Intn(10000),
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
			},
		},
	}
}

func randomFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func sendMockEvent(event IPCOutbound, url string) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshalling event:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("POST request sent. Status Code:", resp.StatusCode)
}

func main() {
	targetURL := "http://localhost:8080/garmin-outbound"
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		mockEvent := generateMockEvent()
		sendMockEvent(mockEvent, targetURL)
	}
}
