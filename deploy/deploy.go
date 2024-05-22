package deploy

import (
	"fmt"
	"ollie/db"
	"ollie/git"
	"ollie/stacks"
	"ollie/styles"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
	gogit "github.com/go-git/go-git/v5"
)

func Deploy() {
	_, err := gogit.PlainOpen(".")
	if err == gogit.ErrRepositoryNotExists {
		log.Fatal("This directory is not a git repository")
		return
	} else if err != nil {
		log.Fatal(err)
		return
	}
	var level string

	form := huh.NewSelect[string]().Title("Select an environment level").
		Options(huh.NewOption("Dev Stack", "stack"),
			huh.NewOption("Staging", "staging"),
			huh.NewOption("Production", "production"),
		).
		Value(&level)

	form.Run()

	switch level {
	case "stack":
		DeployStack()
	case "staging":
		fmt.Println(styles.ErrorStyle.Render("Not yet implemented"))
	case "production":
		fmt.Println(styles.ErrorStyle.Render("Not yet implemented"))
	default:
		log.Fatal("Invalid environment level")
	}
}

func DeployStack() {
	stack := stacks.SelectStack()

	version := git.GetVersion()
	fmt.Println(styles.HighlightStyle.Render(fmt.Sprintf("Current version is %s", version)))

	var bump string

	versionBump := huh.NewSelect[string]().
		Title("How big of a bump?").
		Options(
			huh.NewOption("Major", "major"),
			huh.NewOption("Minor", "minor"),
			huh.NewOption("Patch", "patch"),
			huh.NewOption("Don't bump", "same")).
		Value(&bump)

	versionBump.Run()

	newVersion, err := git.VersionBump(version, bump, false)

	fmt.Println(styles.HighlightStyle.Render(fmt.Sprintf("New version is %s", newVersion)))
	err = spinner.New().
		Title("Deploying").
		Run()
	if err != nil {
		log.Fatal(err)
	}
}
