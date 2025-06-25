package status

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/taciturnaxolotl/ctfd-alerts/clients"
)

var (
	// Colors
	purple    = lipgloss.AdaptiveColor{Light: "#9D8EFF", Dark: "#7D56F4"}
	gray      = lipgloss.AdaptiveColor{Light: "#BEBEBE", Dark: "#4A4A4A"}
	lightGray = lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#3A3A3A"}
	green     = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	red       = lipgloss.AdaptiveColor{Light: "#F52D4F", Dark: "#FF5F7A"}

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)

	containerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true).
			Align(lipgloss.Center)

	oddRowStyle = lipgloss.NewStyle()

	evenRowStyle = lipgloss.NewStyle()
	// Status indicators
	solvedStyle   = lipgloss.NewStyle().Foreground(green).SetString("✓")
	unsolvedStyle = lipgloss.NewStyle().Foreground(red).SetString("✗")
)

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func createTable(headers []string, rows [][]string) string {

	t := table.New().
		Border(lipgloss.ASCIIBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lightGray)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers(headers...).
		Rows(rows...)

	return t.Render()
}

func runDashboard(cmd *cobra.Command, args []string) {
	// Get CTFd client from root command context
	ctfdClient := cmd.Context().Value("ctfd_client").(CTFdClient)

	// Get scoreboard data
	scoreboard, err := ctfdClient.GetScoreboard()
	if err != nil {
		log.Fatalf("Error fetching scoreboard: %v", err)
	}

	// Prepare scoreboard data
	scoreboardHeaders := []string{"Pos", "Team", "Score", "Members"}
	scoreboardRows := make([][]string, len(scoreboard.Data))

	for i, team := range scoreboard.Data {
		memberNames := make([]string, 0, len(team.Members))
		for _, member := range team.Members {
			memberNames = append(memberNames, member.Name)
		}
		scoreboardRows[i] = []string{
			fmt.Sprintf("%d", team.Position),
			truncateString(team.Name, 24),
			fmt.Sprintf("%d", team.Score),
			truncateString(strings.Join(memberNames, ", "), 39),
		}
	}

	// Get challenge list
	challenges, err := ctfdClient.GetChallengeList()
	if err != nil {
		log.Fatalf("Error fetching challenges: %v", err)
	}

	// Prepare challenge data
	challengeHeaders := []string{"ID", "Name", "Category", "Value", "Solves", "Solved"}
	challengeRows := make([][]string, len(challenges.Data))

	for i, challenge := range challenges.Data {
		solvedStatus := unsolvedStyle.String()
		if challenge.SolvedByMe {
			solvedStatus = solvedStyle.String()
		}
		challengeRows[i] = []string{
			fmt.Sprintf("%d", challenge.ID),
			truncateString(challenge.Name, 24),
			truncateString(challenge.Category, 14),
			fmt.Sprintf("%d", challenge.Value),
			fmt.Sprintf("%d", challenge.Solves),
			solvedStatus,
		}
	}

	// Build and render the complete dashboard
	var dashboard strings.Builder

	// Scoreboard section
	dashboard.WriteString(titleStyle.Render(fmt.Sprintf("CTFd Scoreboard [%d]", len(scoreboard.Data))))
	dashboard.WriteString("\n")
	dashboard.WriteString(createTable(scoreboardHeaders, scoreboardRows))

	// Challenges section
	dashboard.WriteString("\n\n")
	dashboard.WriteString(titleStyle.Render(fmt.Sprintf("CTFd Challenges [%d]", len(challenges.Data))))
	dashboard.WriteString("\n")
	dashboard.WriteString(createTable(challengeHeaders, challengeRows))

	// Render the final output
	fmt.Print("\n")
	fmt.Print(dashboard.String())
	fmt.Print("\n")
}

// CTFdClient alias for the client interface from the clients package
type CTFdClient = clients.CTFdClient
