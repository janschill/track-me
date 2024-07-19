# IPC Outbound API Documentation

## Overview

```go
type IPCOutbound struct {
  Version string  `json:"Version"`
  Events  []Event `json:"Events"`
}

// Event object containing data from an inReach device
type Event struct {
  IMEI        string     `json:"imei"`
  MessageCode int        `json:"messageCode"`
  FreeText    string     `json:"freeText"`
  TimeStamp   int64      `json:"timeStamp"`
  Addresses   []Address  `json:"addresses"`
  Point       Point      `json:"point"`
  Status      Status     `json:"status"`
}

// Address object within an Event
type Address struct {
  Address string `json:"address"`
}

// Point object within an Event
type Point struct {
  Latitude  float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
  Altitude  float64 `json:"altitude"`
  GPSFix    int     `json:"gpsFix"`
  Course    float64 `json:"course"`
  Speed     float64 `json:"speed"`
}

// Status object within an Event
type Status struct {
  Autonomous    int `json:"autonomous"`
  LowBattery    int `json:"lowBattery"`
  IntervalChange int `json:"intervalChange"`
  ResetDetected int `json:"resetDetected"`
}
```
