package stacks

import (
	"ollie/db"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func SelectStack() (string, error) {
	var selectedStack string
	stacks, err := db.GetStacks()
	if err != nil {
		log.Error(err)
		return "", err
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
}
