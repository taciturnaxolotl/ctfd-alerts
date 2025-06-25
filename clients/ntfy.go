package clients

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// NtfyMessage represents a notification message to be sent via ntfy
type NtfyMessage struct {
	Topic    string           `json:"topic"`
	Message  string           `json:"message,omitempty"`
	Title    string           `json:"title,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
	Priority int              `json:"priority,omitempty"`
	Click    string           `json:"click,omitempty"`
	Actions  []map[string]any `json:"actions,omitempty"`
	Attach   string           `json:"attach,omitempty"`
	Filename string           `json:"filename,omitempty"`
	Markdown bool             `json:"markdown,omitempty"`
	Icon     string           `json:"icon,omitempty"`
	Email    string           `json:"email,omitempty"`
	Call     string           `json:"call,omitempty"`
	Delay    string           `json:"delay,omitempty"`
}

// NtfyClient represents a client for sending notifications via ntfy.sh
type NtfyClient struct {
	Topic      string
	ServerURL  string
	HTTPClient *http.Client
	// Optional authentication token
	AccessToken string
}

// NewNtfyClient creates a new ntfy client with the specified topic and server URL.
// It configures an HTTP client with a 10-second timeout and insecure TLS verification.
func NewNtfyClient(topic, serverURL string, accessToken string) *NtfyClient {
	serverURL = strings.TrimSuffix(serverURL, "/")
	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}

	return &NtfyClient{
		Topic:       topic,
		ServerURL:   serverURL,
		AccessToken: accessToken,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

// NewMessage creates a new NtfyMessage with the specified message content
func (c *NtfyClient) NewMessage(messageText string) *NtfyMessage {
	return &NtfyMessage{
		Topic:   c.Topic,
		Message: messageText,
	}
}

// SendMessage sends a structured NtfyMessage
func (c *NtfyClient) SendMessage(msg *NtfyMessage) error {
	// Ensure topic is set
	if msg.Topic == "" {
		msg.Topic = c.Topic
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	req, err := http.NewRequest("POST", c.ServerURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	if c.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
