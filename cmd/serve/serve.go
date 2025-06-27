package serve

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/taciturnaxolotl/ctfd-alerts/clients"
)

type MonitorState struct {
	LastScoreboard *clients.ScoreboardResponse    `json:"last_scoreboard"`
	LastChallenges *clients.ChallengeListResponse `json:"last_challenges"`
	UserPosition   int                            `json:"user_position"`
}

func getCacheFilePath() string {
	return filepath.Join(".", "cache.json")
}

func loadStateFromCache() *MonitorState {
	cachePath := getCacheFilePath()
	data, err := os.ReadFile(cachePath)
	if err != nil {
		log.Printf("No cache file found or error reading cache: %v", err)
		return &MonitorState{}
	}

	var state MonitorState
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("Error parsing cache file: %v", err)
		return &MonitorState{}
	}

	log.Printf("Loaded state from cache: %s", cachePath)
	return &state
}

func saveStateToCache(state *MonitorState) error {
	cachePath := getCacheFilePath()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling state: %v", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("error writing cache file: %v", err)
	}

	return nil
}

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run monitoring server",
	Long:  "Continuously monitors CTFd for leaderboard changes and new challenges, sending alerts when events occur",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	// Get CTFd client from context
	ctfdClient, ok := ctx.Value("ctfd_client").(clients.CTFdClient)
	if !ok {
		log.Fatal("CTFd client not found in context")
	}

	// Get config from context
	config := ctx.Value("config")

	// Use reflection to access config fields
	configValue := reflect.ValueOf(config).Elem()
	userField := configValue.FieldByName("User").String()
	intervalField := int(configValue.FieldByName("MonitorInterval").Int())

	ntfyConfigField := configValue.FieldByName("NtfyConfig")
	ntfyTopic := ntfyConfigField.FieldByName("Topic").String()
	ntfyApiBase := ntfyConfigField.FieldByName("ApiBase").String()
	ntfyAccessToken := ntfyConfigField.FieldByName("AccessToken").String()

	// Create ntfy client
	ntfyClient := clients.NewNtfyClient(ntfyTopic, ntfyApiBase, ntfyAccessToken)

	// Initialize monitoring state - try to load from cache first
	state := loadStateFromCache()

	// If cache is empty or we want fresh data, get initial state from API
	if state.LastScoreboard == nil || state.LastChallenges == nil {
		log.Println("No cached state found, fetching initial state from API...")
		if err := updateState(ctfdClient, state, userField); err != nil {
			log.Printf("Error getting initial state: %v", err)
		}
	} else {
		log.Println("Using cached state")
		// Still update user position in case it changed
		if state.LastScoreboard != nil {
			state.UserPosition = findUserPosition(state.LastScoreboard, userField)
		}
	}

	log.Printf("Starting monitoring server (interval: %d seconds)", intervalField)
	log.Printf("Monitoring user: %s", userField)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main monitoring loop
	ticker := time.NewTicker(time.Duration(intervalField) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := monitorAndAlert(ctfdClient, ntfyClient, state, userField); err != nil {
				log.Printf("Error during monitoring: %v", err)
			} else {
				// Save state to cache after successful monitoring
				if err := saveStateToCache(state); err != nil {
					log.Printf("Error saving state to cache: %v", err)
				}
			}
		case <-sigChan:
			log.Println("Received shutdown signal, saving state and stopping server...")
			if err := saveStateToCache(state); err != nil {
				log.Printf("Error saving final state to cache: %v", err)
			} else {
				log.Printf("State saved to cache: %s", getCacheFilePath())
			}
			return
		}
	}
}

func updateState(client clients.CTFdClient, state *MonitorState, username string) error {
	// Get scoreboard
	scoreboard, err := client.GetScoreboard()
	if err != nil {
		return fmt.Errorf("failed to get scoreboard: %v", err)
	}
	state.LastScoreboard = scoreboard

	// Find user position
	state.UserPosition = findUserPosition(scoreboard, username)

	// Get challenges
	challenges, err := client.GetChallengeList()
	if err != nil {
		return fmt.Errorf("failed to get challenges: %v", err)
	}
	state.LastChallenges = challenges

	return nil
}

func monitorAndAlert(client clients.CTFdClient, ntfy *clients.NtfyClient, state *MonitorState, username string) error {
	// Get current scoreboard
	currentScoreboard, err := client.GetScoreboard()
	if err != nil {
		return fmt.Errorf("failed to get scoreboard: %v", err)
	}

	// Get current challenges
	currentChallenges, err := client.GetChallengeList()
	if err != nil {
		return fmt.Errorf("failed to get challenges: %v", err)
	}

	// Check for leaderboard bypass
	if state.LastScoreboard != nil {
		currentPosition := findUserPosition(currentScoreboard, username)
		if currentPosition > state.UserPosition && state.UserPosition > 0 {
			// User was bypassed
			msg := ntfy.NewMessage(fmt.Sprintf("🏆 You've been bypassed on the leaderboard! New position: #%d (was #%d)", currentPosition, state.UserPosition))
			msg.Title = "CTFd Leaderboard Alert"
			msg.Tags = []string{"warning", "leaderboard"}
			msg.Priority = 4

			if err := ntfy.SendMessage(msg); err != nil {
				log.Printf("Failed to send bypass alert: %v", err)
			} else {
				log.Printf("Sent bypass alert: %s -> %d", username, currentPosition)
			}
		}
		state.UserPosition = currentPosition
	}

	// Check for new challenges
	if state.LastChallenges != nil {
		newChallenges := findNewChallenges(state.LastChallenges, currentChallenges)
		for _, challenge := range newChallenges {
			msg := ntfy.NewMessage(fmt.Sprintf("🎯 New challenge released: %s (%s) - %d points", challenge.Name, challenge.Category, challenge.Value))
			msg.Title = "New CTFd Challenge"
			msg.Tags = []string{"challenge", "new"}
			msg.Priority = 3

			if err := ntfy.SendMessage(msg); err != nil {
				log.Printf("Failed to send new challenge alert: %v", err)
			} else {
				log.Printf("Sent new challenge alert: %s", challenge.Name)
			}
		}
	}

	// Update state
	state.LastScoreboard = currentScoreboard
	state.LastChallenges = currentChallenges

	return nil
}

func findUserPosition(scoreboard *clients.ScoreboardResponse, username string) int {
	for _, team := range scoreboard.Data {
		if team.Name == username {
			return team.Position
		}
		// Also check team members
		for _, member := range team.Members {
			if member.Name == username {
				return team.Position
			}
		}
	}
	return 0 // User not found
}

func findNewChallenges(oldChallenges, newChallenges *clients.ChallengeListResponse) []clients.Challenge {
	oldMap := make(map[int]bool)
	for _, challenge := range oldChallenges.Data {
		oldMap[challenge.ID] = true
	}

	var newOnes []clients.Challenge
	for _, challenge := range newChallenges.Data {
		if !oldMap[challenge.ID] {
			newOnes = append(newOnes, challenge)
		}
	}

	return newOnes
}
