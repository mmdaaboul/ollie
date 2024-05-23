package git

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetVersion() string {
	// Fetch and pull latest changes
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Error opening git repository: %s", err)
	}

	FetchAndPull(repo)
	tags, err := repo.Tags()
	if err != nil {
		log.Fatalf("Error getting tags: %s", err)
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
		log.Fatalf("Error getting worktree: %s", err)
	}

	w.Pull(&git.PullOptions{RemoteName: "origin"})
}

func TagAndPush(tagName, message string) error {
	// Fetch and pull latest changes
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Errorf("Error opening git repository: %s", err)
		return err
	}

	created, err := setTag(repo, tagName, message)
	if err != nil {
		log.Printf("create tag error: %s", err)
		return err
	}

	if created {
		err = pushTags(repo)
		if err != nil {
			log.Errorf("push tag error: %s", err)
			return err
		}
	}

	return nil
}

func tagExists(tag string, r *git.Repository) bool {
	tagFoundErr := "tag was found"
	log.Info("git show-ref --tag")
	tags, err := r.TagObjects()
	if err != nil {
		log.Debugf("get tags error: %s", err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			log.Errorf(tagFoundErr)
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		log.Debugf("iterate tags error: %s", err)
		return false
	}
	return res
}

func setTag(r *git.Repository, tag, message string) (bool, error) {
	if tagExists(tag, r) {
		log.Debugf("tag %s already exists", tag)
		return false, nil
	}
	log.Debugf("Set tag %s", tag)
	h, err := r.Head()
	if err != nil {
		log.Debugf("get HEAD error: %s", err)
		return false, err
	}
	log.Infof("git tag -a %s %s -m \"%s\"", tag, h.Hash(), message)
	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Message: message,
	})

	if err != nil {
		log.Debugf("create tag error: %s", err)
		return false, err
	}

	return true, nil
}

func pushTags(r *git.Repository) error {
	cmd := exec.Command("git", "push", "--tags")

	err := cmd.Run()
	if err != nil {
		log.Errorf("Error Pushing Tags: %s", err)
	} else {
		log.Debug("Successfully Pushed Tags")
	}
	return err
}
