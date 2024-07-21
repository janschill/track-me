package db

import "time"

type Trip struct {
	ID          int64
	StartTime   time.Time
	EndTime     time.Time
	Description string
	Events      []Event
}

type Event struct {
	ID          int64
	TripID      int64
	Imei        string
	MessageCode int
	FreeText    string
	TimeStamp   int64
	Addresses   []Address
	Point       Point
	Status      Status
}

type Address struct {
	Address string
}

type Point struct {
	Latitude  float64
	Longitude float64
	Altitude  int
	GpsFix    int
	Course    int
	Speed     int
}

type Status struct {
	Autonomous     int
	LowBattery     int
	IntervalChange int
	ResetDetected  int
}
