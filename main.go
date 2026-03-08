package main

import (
	"os"

	"github.com/tailscale/walk"
)

func main() {
	_, err := walk.InitApp()
	failOnError(err, "Failed to initialize app")

	ui, err := createUI()
	failOnError(err, "Failed to create UI")

	err = ui.loadSessions()
	failOnError(err, "Failed to load PuTTY sessions")

	ui.Run()
}

func failOnError(err error, msg string) {
	if err != nil {
		showTaskDialog(nil, "Error", msg+" Error message: "+err.Error())
		os.Exit(1)
	}
}
