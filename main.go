package main

import (
	"flag"
	"fmt"

	"ollie/deploy"
	"ollie/git"
	logo "ollie/logo"
	"ollie/setup"
	style "ollie/styles"
	"ollie/zookeeper"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

func main() {
	var logLevel string
	flag.StringVar(&logLevel, "log", "error", "set log level (debug, info, warn, error)")
	flag.Parse()

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}

	fmt.Println(style.HighlightStyle.Render(logo.Logo))

	var task string

	// Check if a non-flag argument is provided
	if len(flag.Args()) == 0 {
		form := huh.NewSelect[string]().
			Title("Select an option").
			Options(
				huh.NewOption("Deploy", "deploy"),
				huh.NewOption("Zookeeper", "zookeeper"),
				huh.NewOption("Print Git Tags", "printTags"),
			).
			Value(&task)

		form.Run()
	} else if flag.Args()[0] == "deploy" {
		// Skip the form if the first non-flag argument is "deploy"
		task = "deploy"
	} else if flag.Args()[0] == "zookeeper" {
		task = "zookeeper"
	} else if flag.Args()[0] == "printTags" {
		task = "printTags"
	}

	_, err := setup.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch task {
	case "deploy":
		deploy.Deploy()
	case "zookeeper":
		zookeeper.EnterZookeeper()
	case "printTags":
		git.PrettyPrintTags()
	default:
		log.Error("Invalid task")
	}
}
