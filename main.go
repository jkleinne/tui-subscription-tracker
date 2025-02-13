package main

import (
	"github.com/rivo/tview"
	"subscription-tracker/ui"
)

func main() {
	app := tview.NewApplication()
	ui := ui.NewUI(app)
	
	if err := ui.Run(); err != nil {
		panic(err)
	}
}
