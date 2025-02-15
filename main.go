package main

import (
	"fmt"
	"github.com/rivo/tview"
	"log"
	"os"
	"subscription-tracker/ui"
)

func main() {
	// Initialize logging
	logFile, err := os.OpenFile("subscription-tracker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Starting subscription tracker application")

	// Initialize application
	app := tview.NewApplication()
	defer func() {
		app.Stop()
		log.Println("Application stopped")
	}()

	ui := ui.NewUI(app)

	if err := ui.Run(); err != nil {
		log.Printf("Error running application: %v\n", err)
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
