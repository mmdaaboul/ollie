package git

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrettyPrintTags() {
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

	// Collect all tags in a slice
	var allTags []*object.Tag
	err = tags.ForEach(func(t *object.Tag) error {
		allTags = append(allTags, t)
		return nil
	})
	if err != nil {
		fmt.Println("Error iterating tags:", err)
		return
	}

	// Sort tags by date in descending order
	sort.Slice(allTags, func(i, j int) bool {
		return allTags[i].Tagger.When.After(allTags[j].Tagger.When)
	})

	// Select the last 10 tags
	if len(allTags) > 10 {
		allTags = allTags[:10]
	}

	// Define table headers
	headers := table.Row{"Tag", "Message", "Date", "Author"}

	// Create the table writer
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(headers)

	// Build table rows
	for _, tag := range allTags {
		row := table.Row{
			tag.Name,
			tag.Message, // No truncation
			formatDate(tag.Tagger.When),
			formatCommitter(tag.Tagger.Name),
		}
		tw.AppendRow(row)
	}

	// Render the table
	tw.Render()
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatCommitter(fullName string) string {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return fullName
	}
	firstName := parts[0]
	if len(parts) > 1 {
		lastNameInitial := string(parts[1][0])
		return fmt.Sprintf("%s %s.", firstName, lastNameInitial)
	}
	return firstName
}
