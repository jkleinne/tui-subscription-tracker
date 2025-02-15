package models

import (
	"fmt"
	"log"
	"sync"
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
	mu                sync.RWMutex
	name              string    `json:"name"`
	cost              float64   `json:"cost"`
	paymentFrequency  string    `json:"payment_frequency"`
	nextPaymentDate   time.Time `json:"next_payment_date"`
	remainingPayments int       `json:"remaining_payments"`
	totalPayments     int       `json:"total_payments"`
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

	sub := &Subscription{
		name:              name,
		cost:              cost,
		paymentFrequency:  frequency,
		nextPaymentDate:   nextPayment,
		remainingPayments: totalPayments,
		totalPayments:     totalPayments,
	}

	log.Printf("Created new subscription: %s with %d total payments", name, totalPayments)
	return sub, nil
}

// Getter methods
func (s *Subscription) Name() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.name
}

func (s *Subscription) Cost() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cost
}

func (s *Subscription) PaymentFrequency() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.paymentFrequency
}

func (s *Subscription) NextPaymentDate() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nextPaymentDate
}

func (s *Subscription) RemainingPayments() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.remainingPayments
}

func (s *Subscription) TotalPayments() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalPayments
}

func (s *Subscription) TimeUntilNextPayment() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Until(s.nextPaymentDate)
}

func (s *Subscription) FormattedTimeUntilNextPayment() string {
	duration := s.TimeUntilNextPayment()
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

func (s *Subscription) calculateNextPaymentDate(current time.Time) time.Time {
	switch s.paymentFrequency {
	case FrequencyDaily:
		return current.AddDate(0, 0, 1)
	case FrequencyWeekly:
		return current.AddDate(0, 0, 7)
	case FrequencyMonthly:
		// Handle edge cases for months with different lengths
		nextMonth := current.AddDate(0, 1, 0)
		if current.Day() != nextMonth.Day() {
			// If the day changed (e.g., Jan 31 -> Feb 28), adjust to end of month
			if current.Day() > nextMonth.Day() {
				return time.Date(nextMonth.Year(), nextMonth.Month()+1, 1,
					current.Hour(), current.Minute(), current.Second(),
					current.Nanosecond(), current.Location()).Add(-time.Second)
			}
		}
		return nextMonth
	case FrequencyYearly:
		// Handle leap years
		nextYear := current.AddDate(1, 0, 0)
		if current.Month() == time.February && current.Day() == 29 {
			// If it's Feb 29 and next year is not a leap year, use Feb 28
			if nextYear.Month() == time.March {
				return time.Date(nextYear.Year(), time.February, 28,
					current.Hour(), current.Minute(), current.Second(),
					current.Nanosecond(), current.Location())
			}
		}
		return nextYear
	default:
		// This shouldn't happen due to validation in NewSubscription
		log.Printf("Warning: Invalid payment frequency %s for subscription %s",
			s.paymentFrequency, s.name)
		return current
	}
}

func (s *Subscription) ProcessPayment() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.remainingPayments <= 0 {
		return fmt.Errorf("subscription has ended")
	}

	s.remainingPayments--
	oldDate := s.nextPaymentDate
	s.nextPaymentDate = s.calculateNextPaymentDate(oldDate)

	log.Printf("Processed payment for subscription %s. Remaining payments: %d/%d",
		s.name, s.remainingPayments, s.totalPayments)

	return nil
}

func (s *Subscription) Status() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.remainingPayments <= 0 {
		return "Completed"
	}
	return fmt.Sprintf("Active (%d/%d payments remaining)", s.remainingPayments, s.totalPayments)
}
