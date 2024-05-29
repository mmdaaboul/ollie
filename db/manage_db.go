package db

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func ManageDb() {
	// TODO: Other things to add later, delete/add release docs, edit for stacks and release docs
	var opt string

	form := huh.NewSelect[string]().
		Title("Select What to do").
		Value(&opt).
		Options(huh.NewOption("Delete Stacks", "delete"), huh.NewOption("Add Stack", "add"))

	form.Run()

	switch opt {
	case "delete":
		deleteStacks()
	case "add":
		addStack()
	default:
		log.Fatal("Bad management selection")
	}
}

func deleteStacks() {
	stacks, err := GetStacks()
	if err != nil {
		log.Fatalf("Error getting stacks in deleteStacks: %s", err)
	}

	var selected []string
	var options []huh.Option[string]

	for _, stack := range stacks {
		log.Debugf("Adding %s as an option", stack)
		option := huh.NewOption(stack, stack)
		options = append(options, option)
	}

	form := huh.NewMultiSelect[string]().
		Options(options...).
		Title("Select Stacks to Delete").
		Value(&selected)

	form.Run()

	for _, selection := range selected {
		log.Debugf("Deleting %s", selection)
		DeleteStack(selection)
	}

	log.Print("Deletion Complete")
}

func addStack() {
	var newStack string
	form := huh.NewInput().
		Title("Enter new name for the stack").
		Value(&newStack)

	form.Run()

	alreadyInDb, err := HasStack(newStack)
	if err != nil {
		log.Fatalf("Error checking database: %s", err)
	}
	if alreadyInDb {
		log.Printf("%s already in database, updating instead", newStack)
		UpdateStack(newStack)
		return
	}

	AddStack(newStack)
	log.Printf("Added %s to database", newStack)
}
