package clients

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// CTFdClient interface defines the methods required for interacting with CTFd
type CTFdClient interface {
	GetScoreboard() (*ScoreboardResponse, error)
	GetChallengeList() (*ChallengeListResponse, error)
}

// ctfdClient represents a CTFd API client implementation
type ctfdClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// ScoreboardResponse represents the top-level response from the CTFd API for scoreboard
type ScoreboardResponse struct {
	Success bool           `json:"success"`
	Data    []TeamStanding `json:"data"`
}

// TeamStanding represents a team's standing on the scoreboard
type TeamStanding struct {
	Position    int      `json:"pos"`
	AccountID   int      `json:"account_id"`
	AccountURL  string   `json:"account_url"`
	AccountType string   `json:"account_type"`
	OAuthID     *string  `json:"oauth_id"`
	Name        string   `json:"name"`
	Score       int      `json:"score"`
	BracketID   *string  `json:"bracket_id"`
	BracketName *string  `json:"bracket_name"`
	Members     []Member `json:"members"`
}

// Member represents a team member
type Member struct {
	ID          int     `json:"id"`
	OAuthID     *string `json:"oauth_id"`
	Name        string  `json:"name"`
	Score       int     `json:"score"`
	BracketID   *string `json:"bracket_id"`
	BracketName *string `json:"bracket_name"`
}

// ChallengeListResponse represents the top-level response from the CTFd API for challenges
type ChallengeListResponse struct {
	Success bool        `json:"success"`
	Data    []Challenge `json:"data"`
}

// Challenge represents a CTFd challenge
type Challenge struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Attribution    string         `json:"attribution"`
	ConnectionInfo string         `json:"connection_info"`
	NextID         int            `json:"next_id"`
	MaxAttempts    int            `json:"max_attempts"`
	Value          int            `json:"value"`
	Category       string         `json:"category"`
	Type           string         `json:"type"`
	State          string         `json:"state"`
	Requirements   map[string]any `json:"requirements"`
	Solves         int            `json:"solves"`
	SolvedByMe     bool           `json:"solved_by_me"`
}

// NewCTFdClient creates a new CTFd client with the specified base URL and API token.
// It configures an HTTP client with a 10-second timeout and insecure TLS verification.
func NewCTFdClient(baseURL, apiToken string) CTFdClient {
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &ctfdClient{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

// GetScoreboard fetches the CTFd scoreboard data from the API.
// Returns a ScoreboardResponse containing team standings or an error if the request fails.
func (c *ctfdClient) GetScoreboard() (*ScoreboardResponse, error) {
	endpoint := "/scoreboard"

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	var scoreboard ScoreboardResponse
	if err := json.Unmarshal(body, &scoreboard); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	if !scoreboard.Success {
		return nil, fmt.Errorf("API returned success=false")
	}

	return &scoreboard, nil
}

// GetChallengeList fetches the list of challenges from the CTFd API.
// Returns a ChallengeListResponse containing all challenges sorted by ID or an error if the request fails.
func (c *ctfdClient) GetChallengeList() (*ChallengeListResponse, error) {
	endpoint := "/challenges"

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token "+c.apiToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %s", string(body))
	}

	var challengeList ChallengeListResponse
	if err := json.Unmarshal(body, &challengeList); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	sort.Slice(challengeList.Data, func(i, j int) bool {
		return challengeList.Data[i].ID < challengeList.Data[j].ID
	})

	if !challengeList.Success {
		return nil, fmt.Errorf("API returned success=false")
	}

	return &challengeList, nil
}
