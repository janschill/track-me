package garmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockProcessPayload(payload OutboundPayload) error {
	return nil
}

func mockProcessPayloadWithError(payload OutboundPayload) error {
	return fmt.Errorf("mock error")
}

func TestCreateOutboundEvent_Success(t *testing.T) {
	handler := NewOutboundHandler(mockProcessPayload)

	payload := OutboundPayload{
		Version: "1.0",
		Events: []Event{
			{
				Imei:        "123456789012345",
				MessageCode: 1,
				FreeText:    "Test message",
				TimeStamp:   1622547800,
				Addresses:   []Address{{Address: "test@example.com"}},
				Point: Point{
					Latitude:  37.7749,
					Longitude: -122.4194,
					Altitude:  30.0,
					GpsFix:    1,
					Course:    0.0,
					Speed:     0.0,
				},
				Status: Status{
					Autonomous:     1,
					LowBattery:     0,
					IntervalChange: 0,
					ResetDetected:  0,
				},
				Payload: "Test payload",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "/garmin-outbound", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.CreateOutboundEvent(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Payload received successfully."
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateOutboundEvent_Error(t *testing.T) {
	handler := NewOutboundHandler(mockProcessPayloadWithError)

	payload := OutboundPayload{
		Version: "1.0",
		Events: []Event{
			{
				Imei:        "123456789012345",
				MessageCode: 1,
				FreeText:    "Test message",
				TimeStamp:   1622547800,
				Addresses:   []Address{{Address: "test@example.com"}},
				Point: Point{
					Latitude:  37.7749,
					Longitude: -122.4194,
					Altitude:  30.0,
					GpsFix:    1,
					Course:    0.0,
					Speed:     0.0,
				},
				Status: Status{
					Autonomous:     1,
					LowBattery:     0,
					IntervalChange: 0,
					ResetDetected:  0,
				},
				Payload: "Test payload",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "/garmin-outbound", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.CreateOutboundEvent(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "Error processing payload\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestCreateOutboundEvent_MethodNotAllowed(t *testing.T) {
	handler := NewOutboundHandler(mockProcessPayload)

	req, err := http.NewRequest("GET", "/garmin-outbound", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.CreateOutboundEvent(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	expected := "Method is not supported.\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateOutboundEvent_InvalidJSON(t *testing.T) {
	handler := NewOutboundHandler(mockProcessPayload)

	invalidJSON := "{invalid json}"

	req, err := http.NewRequest("POST", "/garmin-outbound", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.CreateOutboundEvent(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "Error parsing request body\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
