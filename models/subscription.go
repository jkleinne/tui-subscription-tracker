package models

import (
	"fmt"
	"time"
)

type Subscription struct {
	Name            string    `json:"name"`
	Cost            float64   `json:"cost"`
	PaymentFrequency string   `json:"payment_frequency"`
	NextPaymentDate time.Time `json:"next_payment_date"`
}

func NewSubscription(name string, cost float64, frequency string, nextPayment time.Time) (*Subscription, error) {
	if name == "" {
		return nil, fmt.Errorf("subscription name cannot be empty")
	}
	if cost <= 0 {
		return nil, fmt.Errorf("cost must be greater than 0")
	}
	if frequency == "" {
		return nil, fmt.Errorf("payment frequency cannot be empty")
	}
	if nextPayment.Before(time.Now()) {
		return nil, fmt.Errorf("next payment date must be in the future")
	}

	return &Subscription{
		Name:            name,
		Cost:            cost,
		PaymentFrequency: frequency,
		NextPaymentDate: nextPayment,
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
