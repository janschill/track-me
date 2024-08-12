package clients

import (
	"fmt"
	"log"
	"net/http"
)

type GarminClient struct {
	httpClient *http.Client
	address    string
	imei       string
	email    string
	password string
}

type GarminMessage struct {
	Sender     string
	Message    string
	Timestamp  string
	Recipients []string
}

type GarminMessageRequest struct {
	Messages []GarminMessage `json:"Messages"`
}

type ClientConfig struct {
	Address  string
	Imei     string
	Email    string
	Password string
}

func NewGarminClient(config ClientConfig) *GarminClient {
	return &GarminClient{
		httpClient: &http.Client{},
		address:    config.Address,
		imei:       config.Imei,
		email:      config.Email,
		password:   config.Password,
	}
}

func (c *GarminClient) SendMessage(payload GarminMessage) {
	url := fmt.Sprintf("%s/Messaging.svc/Message", c.address)
	payload.Recipients = []string{c.imei}

	log.Print("Sending message", payload)
}
