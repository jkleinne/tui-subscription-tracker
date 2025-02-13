package models

import (
	"fmt"
	"time"
)

// Valid payment frequencies
const (
	FrequencyDaily   = "daily"
	FrequencyWeekly  = "weekly"
	FrequencyMonthly = "monthly"
	FrequencyYearly  = "yearly"
)

var ValidFrequencies = map[string]bool{
	FrequencyDaily:   true,
	FrequencyWeekly:  true,
	FrequencyMonthly: true,
	FrequencyYearly:  true,
}

type Subscription struct {
	Name              string    `json:"name"`
	Cost              float64   `json:"cost"`
	PaymentFrequency  string    `json:"payment_frequency"`
	NextPaymentDate   time.Time `json:"next_payment_date"`
	RemainingPayments int       `json:"remaining_payments"`
	TotalPayments     int       `json:"total_payments"`
}

func NewSubscription(name string, cost float64, frequency string, nextPayment time.Time, totalPayments int) (*Subscription, error) {
	if name == "" {
		return nil, fmt.Errorf("subscription name cannot be empty")
	}
	if cost <= 0 {
		return nil, fmt.Errorf("cost must be greater than 0")
	}
	if !ValidFrequencies[frequency] {
		return nil, fmt.Errorf("invalid payment frequency: must be one of daily, weekly, monthly, or yearly")
	}
	if nextPayment.Before(time.Now()) {
		return nil, fmt.Errorf("next payment date must be in the future")
	}
	if totalPayments <= 0 {
		return nil, fmt.Errorf("total payments must be greater than 0")
	}

	return &Subscription{
		Name:              name,
		Cost:              cost,
		PaymentFrequency:  frequency,
		NextPaymentDate:   nextPayment,
		RemainingPayments: totalPayments,
		TotalPayments:     totalPayments,
	}, nil
}

func (s *Subscription) TimeUntilNextPayment() time.Duration {
	return time.Until(s.NextPaymentDate)
}

func (s *Subscription) FormattedTimeUntilNextPayment() string {
	duration := s.TimeUntilNextPayment()
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

func (s *Subscription) ProcessPayment() error {
	if s.RemainingPayments <= 0 {
		return fmt.Errorf("subscription has ended")
	}

	s.RemainingPayments--

	// Calculate next payment date based on frequency
	switch s.PaymentFrequency {
	case FrequencyDaily:
		s.NextPaymentDate = s.NextPaymentDate.AddDate(0, 0, 1)
	case FrequencyWeekly:
		s.NextPaymentDate = s.NextPaymentDate.AddDate(0, 0, 7)
	case FrequencyMonthly:
		s.NextPaymentDate = s.NextPaymentDate.AddDate(0, 1, 0)
	case FrequencyYearly:
		s.NextPaymentDate = s.NextPaymentDate.AddDate(1, 0, 0)
	}

	return nil
}

func (s *Subscription) Status() string {
	if s.RemainingPayments <= 0 {
		return "Completed"
	}
	return fmt.Sprintf("Active (%d/%d payments remaining)", s.RemainingPayments, s.TotalPayments)
}
