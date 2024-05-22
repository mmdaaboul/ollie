package git

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func GetVersion() string {
	// Fetch and pull latest changes
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error opening git repository: %s", err))
	}

	FetchAndPull(repo)
	tags, err := repo.Tags()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error getting tags: %s", err))
	}

	// Get the latest tag
	var latestTag string
	if err := tags.ForEach(func(ref *plumbing.Reference) error {
		obj, err := repo.TagObject(ref.Hash())
		switch err {
		case nil:
			// Tag object present
			latestTag = obj.Name
			return nil
		case plumbing.ErrObjectNotFound:
			// Not a tag object
			return nil
		default:
			// Some other error
			return err
		}
	}); err != nil {
		// Handle outer iterator error
	}
	return latestTag
}

func FetchAndPull(r *git.Repository) {
	r.Fetch(&git.FetchOptions{RemoteName: "origin"})
	w, err := r.Worktree()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error getting worktree: %s", err))
	}

	w.Pull(&git.PullOptions{RemoteName: "origin"})
}
