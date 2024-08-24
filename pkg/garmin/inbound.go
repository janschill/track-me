package garmin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	httpClient  *http.Client
	address     string
	imei        string
	email       string
	password    string
	rateLimiter *RateLimiter
}

type Message struct {
	Sender     string   `json:"Sender"`
	Message    string   `json:"Message"`
	Timestamp  string   `json:"Timestamp"`
	Recipients []string `json:"Recipients"`
}

type MessageRequestPayload struct {
	Messages []Message `json:"Messages"`
}

type Config struct {
	Address  string
	Imei     string
	Email    string
	Password string
	Limit    int
	Interval time.Duration
}

func NewClient(config Config) *Client {
	return &Client{
		httpClient:  &http.Client{},
		address:     config.Address,
		imei:        config.Imei,
		email:       config.Email,
		password:    config.Password,
		rateLimiter: NewRateLimiter(config.Limit, config.Interval),
	}
}

func (c *Client) newRequest(method, endpoint string, body []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.address, endpoint)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.email, c.password)))
	req.Header.Set("Authorization", "Basic "+auth)

	return req, nil
}

func logRequest(req *http.Request) {
	log.Printf("Request URL: %s", req.URL.String())
	log.Printf("Request Method: %s", req.Method)
	log.Printf("Request Headers: %v", req.Header)
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset the body after reading
		log.Printf("Request Body: %s", string(bodyBytes))
	}
}

func logResponse(res *http.Response) {
	log.Printf("Response Status: %v", res.StatusCode)
	log.Printf("Response Headers: %v", res.Header)
	if res.Body != nil {
		bodyBytes, _ := io.ReadAll(res.Body)
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset the body after reading
		log.Printf("Response Body: %s", string(bodyBytes))
	}
}

func (c *Client) SendMessage(sender string, message string) error {
	if len(message) > 160 {
		return fmt.Errorf("message length exceeds limit")
	}

	if !c.rateLimiter.Allow(c.imei) {
		log.Print("rate limit exceeded")
		return fmt.Errorf("rate limit exceeded")
	}

	endpoint := "/Messaging.svc/Message"
	timestampStr := fmt.Sprintf("/Date(%d)/", time.Now().Unix()*1000)

	payload := MessageRequestPayload{
		Messages: []Message{
			{
				Sender:     sender,
				Message:    message,
				Timestamp:  timestampStr,
				Recipients: []string{c.imei},
			},
		},
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Print("failed to marshal request body: %w", err)
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := c.newRequest("POST", endpoint, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	log.Printf("Message sent successfully: %v", payload)

	return nil
}
