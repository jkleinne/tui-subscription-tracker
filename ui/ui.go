package ui

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"subscription-tracker/models"
	"subscription-tracker/storage"
	"time"

	"github.com/rivo/tview"
)

type UI struct {
	app           *tview.Application
	pages         *tview.Pages
	storage       storage.Storage
	subscriptions *tview.List
	form          *tview.Form
}

func initializeStorage() (storage.Storage, error) {
	dataFilePath := filepath.Join("data", "subscriptions.json")
	return storage.NewJSONStorage(dataFilePath)
}

func NewUI(app *tview.Application) *UI {
	storage, err := initializeStorage()
	if err != nil {
		log.Printf("Failed to initialize storage: %v", err)
		// Show error modal and provide option to retry or exit
		modal := tview.NewModal().
			SetText(fmt.Sprintf("Failed to initialize storage: %v\nWould you like to retry?", err)).
			AddButtons([]string{"Retry", "Exit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Retry" {
					NewUI(app)
				} else {
					app.Stop()
				}
			})
		app.SetRoot(modal, true)
		return nil
	}

	ui := &UI{
		app:     app,
		pages:   tview.NewPages(),
		storage: storage,
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
	ui.setupForm(ui.form, "", ui.saveSubscription)

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

func (ui *UI) setupForm(form *tview.Form, title string, saveFunc func()) {
	form.Clear(true)
	form.
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Cost", "", 20, tview.InputFieldFloat, nil).
		AddInputField("Payment Frequency (daily/weekly/monthly/yearly)", "", 20, nil, nil).
		AddInputField("Next Payment Date (YYYY-MM-DD)", "", 20, nil, nil).
		AddInputField("Total Payments", "", 10, tview.InputFieldInteger, nil).
		AddButton("Save", saveFunc).
		AddButton("Cancel", func() {
			ui.pages.SwitchToPage("menu")
		})
	if title != "" {
		form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	}
}

func (ui *UI) showAddForm() {
	ui.setupForm(ui.form, " Add Subscription ", ui.saveSubscription)
	ui.pages.SwitchToPage("form")
}

func (ui *UI) validateFormInput(name, costStr, frequency, dateStr, totalPaymentsStr string) (float64, time.Time, int, error) {
	var validationErrors []string

	if name == "" {
		validationErrors = append(validationErrors, "Name cannot be empty")
	}

	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil || cost <= 0 {
		validationErrors = append(validationErrors, "Cost must be a positive number")
	}

	if !models.ValidFrequencies[frequency] {
		validationErrors = append(validationErrors, "Invalid payment frequency: must be daily, weekly, monthly, or yearly")
	}

	nextPayment, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		validationErrors = append(validationErrors, "Invalid date format. Please use YYYY-MM-DD")
	}

	totalPayments, err := strconv.Atoi(totalPaymentsStr)
	if err != nil || totalPayments <= 0 {
		validationErrors = append(validationErrors, "Total payments must be a positive number")
	}

	if len(validationErrors) > 0 {
		return 0, time.Time{}, 0, fmt.Errorf(strings.Join(validationErrors, "\n"))
	}

	return cost, nextPayment, totalPayments, nil
}

func (ui *UI) saveSubscription() {
	name := ui.form.GetFormItem(0).(*tview.InputField).GetText()
	costStr := ui.form.GetFormItem(1).(*tview.InputField).GetText()
	frequency := ui.form.GetFormItem(2).(*tview.InputField).GetText()
	dateStr := ui.form.GetFormItem(3).(*tview.InputField).GetText()
	totalPaymentsStr := ui.form.GetFormItem(4).(*tview.InputField).GetText()

	cost, nextPayment, totalPayments, err := ui.validateFormInput(name, costStr, frequency, dateStr, totalPaymentsStr)
	if err != nil {
		ui.showError(err.Error())
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

	ui.showSuccess("Subscription added successfully")
	ui.pages.SwitchToPage("menu")
}

func (ui *UI) showSubscriptions() {
	ui.subscriptions.Clear()

	subs := ui.storage.GetSubscriptions()
	for _, sub := range subs {
		timeLeft := sub.FormattedTimeUntilNextPayment()
		description := fmt.Sprintf("Cost: $%.2f | Frequency: %s | Next Payment: %s (%s) | %s",
			sub.Cost(),
			sub.PaymentFrequency(),
			sub.NextPaymentDate().Format("2006-01-02"),
			timeLeft,
			sub.Status())

		// Create a copy of sub for the closure
		currentSub := sub
		ui.subscriptions.AddItem(sub.Name(), description, 0, func() {
			ui.showSubscriptionMenu(currentSub)
		})
	}

	ui.subscriptions.AddItem("Back to Menu", "Return to main menu", 'b', func() {
		ui.pages.SwitchToPage("menu")
	})

	ui.pages.SwitchToPage("list")
}

func (ui *UI) showSubscriptionMenu(sub *models.Subscription) {
	contextMenu := tview.NewModal().
		SetText(fmt.Sprintf("Selected: %s\nWhat would you like to do?", sub.Name())).
		AddButtons([]string{"Edit", "Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Edit":
				ui.showEditForm(sub)
			case "Delete":
				ui.showDeleteConfirmation(sub)
			}
			ui.pages.RemovePage("context")
		})
	ui.pages.AddPage("context", contextMenu, false, true)
}

func (ui *UI) showEditForm(sub *models.Subscription) {
	form := tview.NewForm()
	form.
		AddInputField("Name", sub.Name(), 30, nil, nil).
		AddInputField("Cost", fmt.Sprintf("%.2f", sub.Cost()), 20, tview.InputFieldFloat, nil).
		AddInputField("Payment Frequency (daily/weekly/monthly/yearly)", sub.PaymentFrequency(), 20, nil, nil).
		AddInputField("Next Payment Date (YYYY-MM-DD)", sub.NextPaymentDate().Format("2006-01-02"), 20, nil, nil).
		AddInputField("Total Payments", fmt.Sprintf("%d", sub.TotalPayments()), 10, tview.InputFieldInteger, nil).
		AddButton("Save", func() {
			name := form.GetFormItem(0).(*tview.InputField).GetText()
			costStr := form.GetFormItem(1).(*tview.InputField).GetText()
			frequency := form.GetFormItem(2).(*tview.InputField).GetText()
			dateStr := form.GetFormItem(3).(*tview.InputField).GetText()
			totalPaymentsStr := form.GetFormItem(4).(*tview.InputField).GetText()

			cost, nextPayment, totalPayments, err := ui.validateFormInput(name, costStr, frequency, dateStr, totalPaymentsStr)
			if err != nil {
				ui.showError(err.Error())
				return
			}

			updatedSub, err := models.NewSubscription(name, cost, frequency, nextPayment, totalPayments)
			if err != nil {
				ui.showError(err.Error())
				return
			}

			if err := ui.storage.UpdateSubscription(sub.Name(), updatedSub); err != nil {
				ui.showError(err.Error())
				return
			}

			ui.showSuccess("Subscription updated successfully")
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
		SetText(fmt.Sprintf("Are you sure you want to delete the subscription '%s'?", sub.Name())).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				if err := ui.storage.DeleteSubscription(sub.Name()); err != nil {
					ui.showError(err.Error())
				} else {
					ui.showSuccess("Subscription deleted successfully")
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

func (ui *UI) showSuccess(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.pages.RemovePage("success")
		})

	ui.pages.AddPage("success", modal, false, true)
}

func (ui *UI) Run() error {
	if ui == nil {
		return fmt.Errorf("UI not properly initialized")
	}

	ui.app.SetRoot(ui.pages, true)
	return ui.app.Run()
}
