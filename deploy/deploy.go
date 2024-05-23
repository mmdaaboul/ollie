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
		log.Fatalf("Error opening the git repo: %s", err)
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
		deployStack()
	case "staging":
		deployStaging()
	case "production":
		deployProd()
	default:
		log.Fatal("Invalid environment level")
	}
}

func deployStack() {
	stack, err := stacks.SelectStack()
	if err != nil {
		log.Fatal("There was an issue getting the stack", err)
	}

	log.Debug(styles.HighlightStyle.Render(fmt.Sprintf("Deploying to stack %s", stack)))

	version := git.GetVersion()
	log.Debug(styles.HighlightStyle.Render(fmt.Sprintf("Current version is %s", version)))

	bump := versionBump()
	newVersion, err := git.VersionBump(version, bump, false, false)

	interfaces := releaseInterfaces()

	log.Debug(styles.HighlightStyle.Render(fmt.Sprintf("New version is %s, stack is %s", newVersion, stack)))
	err = spinner.New().
		Title("Deploying...").
		Action(func() { TagAndPush(newVersion, stack, interfaces) }).
		Run()
	if err != nil {
		log.Fatal(err)
	}
}

func deployStaging() {
	version := git.GetVersion()
	bump := versionBump()
	newVersion, err := git.VersionBump(version, bump, false, true)
	if err != nil {
		log.Fatalf("Unable to complete version bump: %s", err)
	}
	interfaces := releaseInterfaces()

	log.Debugf("New version is %s, going to staging", newVersion)

	err = spinner.New().
		Title("Deploying...").
		Action(func() { TagAndPush(newVersion, "staging", interfaces) }).
		Run()

	if err != nil {
		log.Fatal(err)
	}
}

func deployProd() {
	version := git.GetVersion()
	bump := versionBump()
	newVersion, err := git.VersionBump(version, bump, true, false)
	if err != nil {
		log.Fatalf("Unable to complete version bump: %s", err)
	}
	interfaces := releaseInterfaces()

	log.Debugf("New version is %s, going to prod", newVersion)

	err = spinner.New().
		Title("Deploying...").
		Action(func() { TagAndPush(newVersion, "", interfaces) }).
		Run()

	if err != nil {
		log.Fatalf("There was an issue with TagAndPush: %s", err)
	}
}

func TagAndPush(tag string, stack string, release bool) {
	var message string
	if stack != "" {
		message = fmt.Sprintf("%s|%s", stack, tag)
	} else {
		message = getReleaseDoc()
	}

	if release {
		message += "|ri"
	}
	log.Debug(fmt.Sprintf("Tag: %s, Message: %s ", tag, message))
	err := git.TagAndPush(tag, message)
	if err != nil {
		log.Fatalf("Error tagging and pushing: %s", err)
	}
}

func releaseInterfaces() bool {
	var release bool

	setRelease := huh.NewConfirm().
		Title("Release Interfaces").
		Value(&release)

	setRelease.Run()
	return release
}

func versionBump() string {
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

	return bump
}

func getReleaseDoc() string {
	var selectedDoc string
	docs, err := db.GetReleaseDocs()
	if err != nil {
		log.Fatalf("Error returned from GetReleaseDocs: %s", err)
		return ""
	}

	if len(docs) <= 0 {
		form := huh.NewInput().Title("Enter Release Doc").Value(&selectedDoc)
		form.Run()
		db.AddReleaseDoc(selectedDoc)
	} else {
		options := []huh.Option[string]{}
		for _, doc := range docs {
			options = append(options, huh.NewOption(doc, doc))
		}
		options = append(options, huh.NewOption("New Doc", "new"))

		form := huh.NewSelect[string]().Title("Select a doc").
			Options(options...).
			Value(&selectedDoc)

		form.Run()
		var newDoc string

		if selectedDoc == "new" {
			form := huh.NewInput().Title("Enter a doc URL").Value(&newDoc)
			form.Run()
			selectedDoc = newDoc
			db.AddReleaseDoc(newDoc)
			log.Debug(fmt.Sprintf("Added %s to the local db", newDoc))
		}
	}

	return selectedDoc
}
