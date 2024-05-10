package deploy

import (
	"fmt"
	"ollie/db"
	"ollie/git"
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
	var selectedStack string
	stacks, err := db.GetStacks()
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(stacks) <= 0 {
		form := huh.NewInput().Title("Enter a stack name").Value(&selectedStack)
		form.Run()
		db.AddStack(selectedStack)
	} else {
		options := []huh.Option[string]{}
		for _, stack := range stacks {
			options = append(options, huh.NewOption(stack, stack))
		}
		options = append(options, huh.NewOption("New Stack", "new"))

		form := huh.NewSelect[string]().Title("Select a stack").
			Options(options...).
			Value(&selectedStack)

		form.Run()

		if selectedStack == "new" {
			form := huh.NewInput().Title("Enter a stack name").Value(&selectedStack)
			form.Run()
		}
		db.AddStack(selectedStack)
	}

	version := git.GetVersion()
	fmt.Println(styles.HighlightStyle.Render(fmt.Sprintf("Deploying version %s to stack %s", version, selectedStack)))

	err = spinner.New().
		Title("Deploying").
		Run()
	if err != nil {
		log.Fatal(err)
	}
}
