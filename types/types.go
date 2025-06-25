package types

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
