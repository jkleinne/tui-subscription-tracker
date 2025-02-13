package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"subscription-tracker/models"
	"subscription-tracker/storage"
	"time"
)

type UI struct {
	app           *tview.Application
	pages         *tview.Pages
	storage       storage.Storage
	subscriptions *tview.List
	form          *tview.Form
}

func NewUI(app *tview.Application) *UI {
	ui := &UI{
		app:     app,
		pages:   tview.NewPages(),
		storage: storage.NewMemoryStorage(),
	}

	ui.setupPages()
	return ui
}

func (ui *UI) setupPages() {
	// Create main menu
	menu := tview.NewList().
		AddItem("Add Subscription", "Add a new subscription", 'a', ui.showAddForm).
		AddItem("List Subscriptions", "View all subscriptions", 'l', ui.showSubscriptions).
		AddItem("Quit", "Exit the application", 'q', func() {
			ui.app.Stop()
		})
	menu.SetBorder(true).SetTitle(" Main Menu ").SetTitleAlign(tview.AlignLeft)

	// Create subscription form
	ui.form = tview.NewForm()
	ui.form.
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Cost", "", 20, tview.InputFieldFloat, nil).
		AddInputField("Payment Frequency (daily/weekly/monthly/yearly)", "", 20, nil, nil).
		AddInputField("Next Payment Date (YYYY-MM-DD)", "", 20, nil, nil).
		AddInputField("Total Payments", "", 10, tview.InputFieldInteger, nil).
		AddButton("Save", ui.saveSubscription).
		AddButton("Cancel", func() {
			ui.pages.SwitchToPage("menu")
		})
	ui.form.SetBorder(true).SetTitle(" Add Subscription ").SetTitleAlign(tview.AlignLeft)

	// Create subscriptions list
	ui.subscriptions = tview.NewList().
		AddItem("Back to Menu", "Return to main menu", 'b', func() {
			ui.pages.SwitchToPage("menu")
		})
	ui.subscriptions.SetBorder(true).SetTitle(" Active Subscriptions ").SetTitleAlign(tview.AlignLeft)

	// Add pages
	ui.pages.AddPage("menu", menu, true, true)
	ui.pages.AddPage("form", ui.form, true, false)
	ui.pages.AddPage("list", ui.subscriptions, true, false)
}

func (ui *UI) showAddForm() {
	ui.form.Clear(true)
	ui.form.
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Cost", "", 20, tview.InputFieldFloat, nil).
		AddInputField("Payment Frequency (daily/weekly/monthly/yearly)", "", 20, nil, nil).
		AddInputField("Next Payment Date (YYYY-MM-DD)", "", 20, nil, nil).
		AddInputField("Total Payments", "", 10, tview.InputFieldInteger, nil).
		AddButton("Save", ui.saveSubscription).
		AddButton("Cancel", func() {
			ui.pages.SwitchToPage("menu")
		})
	ui.pages.SwitchToPage("form")
}

func (ui *UI) saveSubscription() {
	name := ui.form.GetFormItem(0).(*tview.InputField).GetText()
	costStr := ui.form.GetFormItem(1).(*tview.InputField).GetText()
	frequency := ui.form.GetFormItem(2).(*tview.InputField).GetText()
	dateStr := ui.form.GetFormItem(3).(*tview.InputField).GetText()
	totalPaymentsStr := ui.form.GetFormItem(4).(*tview.InputField).GetText()

	var cost float64
	fmt.Sscanf(costStr, "%f", &cost)

	var totalPayments int
	fmt.Sscanf(totalPaymentsStr, "%d", &totalPayments)

	nextPayment, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ui.showError("Invalid date format. Please use YYYY-MM-DD")
		return
	}

	sub, err := models.NewSubscription(name, cost, frequency, nextPayment, totalPayments)
	if err != nil {
		ui.showError(err.Error())
		return
	}

	if err := ui.storage.AddSubscription(sub); err != nil {
		ui.showError(err.Error())
		return
	}

	ui.pages.SwitchToPage("menu")
}

func (ui *UI) showSubscriptions() {
	ui.subscriptions.Clear()

	subs := ui.storage.GetSubscriptions()
	for _, sub := range subs {
		timeLeft := sub.FormattedTimeUntilNextPayment()
		description := fmt.Sprintf("Cost: $%.2f | Frequency: %s | Next Payment: %s (%s) | %s",
			sub.Cost,
			sub.PaymentFrequency,
			sub.NextPaymentDate.Format("2006-01-02"),
			timeLeft,
			sub.Status())

		// Create a copy of sub for the closure
		currentSub := sub
		ui.subscriptions.AddItem(sub.Name, description, 0, func() {
			// Show context menu for the selected subscription
			contextMenu := tview.NewModal().
				SetText(fmt.Sprintf("Selected: %s\nWhat would you like to do?", currentSub.Name)).
				AddButtons([]string{"Edit", "Delete", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					switch buttonLabel {
					case "Edit":
						ui.showEditForm(currentSub)
					case "Delete":
						ui.showDeleteConfirmation(currentSub)
					}
					ui.pages.RemovePage("context")
				})
			ui.pages.AddPage("context", contextMenu, false, true)
		})
	}

	ui.subscriptions.AddItem("Back to Menu", "Return to main menu", 'b', func() {
		ui.pages.SwitchToPage("menu")
	})

	ui.pages.SwitchToPage("list")
}

func (ui *UI) showEditForm(sub *models.Subscription) {
	form := tview.NewForm()
	form.
		AddInputField("Name", sub.Name, 30, nil, nil).
		AddInputField("Cost", fmt.Sprintf("%.2f", sub.Cost), 20, tview.InputFieldFloat, nil).
		AddInputField("Payment Frequency (daily/weekly/monthly/yearly)", sub.PaymentFrequency, 20, nil, nil).
		AddInputField("Next Payment Date (YYYY-MM-DD)", sub.NextPaymentDate.Format("2006-01-02"), 20, nil, nil).
		AddInputField("Total Payments", fmt.Sprintf("%d", sub.TotalPayments), 10, tview.InputFieldInteger, nil).
		AddButton("Save", func() {
			name := form.GetFormItem(0).(*tview.InputField).GetText()
			costStr := form.GetFormItem(1).(*tview.InputField).GetText()
			frequency := form.GetFormItem(2).(*tview.InputField).GetText()
			dateStr := form.GetFormItem(3).(*tview.InputField).GetText()
			totalPaymentsStr := form.GetFormItem(4).(*tview.InputField).GetText()

			var cost float64
			fmt.Sscanf(costStr, "%f", &cost)

			var totalPayments int
			fmt.Sscanf(totalPaymentsStr, "%d", &totalPayments)

			nextPayment, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				ui.showError("Invalid date format. Please use YYYY-MM-DD")
				return
			}

			updatedSub, err := models.NewSubscription(name, cost, frequency, nextPayment, totalPayments)
			if err != nil {
				ui.showError(err.Error())
				return
			}

			if err := ui.storage.UpdateSubscription(sub.Name, updatedSub); err != nil {
				ui.showError(err.Error())
				return
			}

			ui.pages.RemovePage("edit")
			ui.showSubscriptions()
		}).
		AddButton("Cancel", func() {
			ui.pages.RemovePage("edit")
		})

	form.SetBorder(true).SetTitle(" Edit Subscription ").SetTitleAlign(tview.AlignLeft)
	ui.pages.AddPage("edit", form, true, true)
}

func (ui *UI) showDeleteConfirmation(sub *models.Subscription) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete the subscription '%s'?", sub.Name)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				if err := ui.storage.DeleteSubscription(sub.Name); err != nil {
					ui.showError(err.Error())
				} else {
					ui.showSubscriptions()
				}
			}
			ui.pages.RemovePage("delete")
		})
	ui.pages.AddPage("delete", modal, false, true)
}

func (ui *UI) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.pages.RemovePage("error")
		})

	ui.pages.AddPage("error", modal, false, true)
}

func (ui *UI) Run() error {
	ui.app.SetRoot(ui.pages, true)
	return ui.app.Run()
}
