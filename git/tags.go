package git

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func PrettyPrintTags() {
	re := lipgloss.NewRenderer(os.Stdout)
	// Open the git repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println("Error opening git repository:", err)
		return
	}

	// Get all tags
	tags, err := repo.TagObjects()
	if err != nil {
		fmt.Println("Error getting tags:", err)
		return
	}

	// Define table headers
	headers := []string{"Tag", "Message", "Date", "Committer"}

	var rows [][]string
	// Build table rows
	err = tags.ForEach(func(t *object.Tag) error {
		row := []string{
			t.Name,
			truncateString(t.Message, 30),
			formatDate(t.Tagger.When),
			t.Tagger.Name,
		}
		rows = append(rows, row)
		return nil
	})
	if err != nil {
		fmt.Println("Error iterating tags:", err)
		return
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Rows(rows...)

	// Render the table
	fmt.Print(t)
}

func truncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength] + "..."
	}
	return str
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
