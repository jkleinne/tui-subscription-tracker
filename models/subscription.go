package models

import (
	"fmt"
	"strings"
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

// ValidationError represents multiple validation errors
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return strings.Join(v.Errors, "; ")
}

// Subscription represents a subscription with thread-safe operations
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
	var validationErrors []string

	if name == "" {
		validationErrors = append(validationErrors, "subscription name cannot be empty")
	}
	if cost <= 0 {
		validationErrors = append(validationErrors, "cost must be greater than 0")
	}
	if !ValidFrequencies[frequency] {
		validationErrors = append(validationErrors, "invalid payment frequency: must be one of daily, weekly, monthly, or yearly")
	}
	if nextPayment.Before(time.Now()) {
		validationErrors = append(validationErrors, "next payment date must be in the future")
	}
	if totalPayments <= 0 {
		validationErrors = append(validationErrors, "total payments must be greater than 0")
	}

	if len(validationErrors) > 0 {
		return nil, &ValidationError{Errors: validationErrors}
	}

	return &Subscription{
		name:              name,
		cost:              cost,
		paymentFrequency:  frequency,
		nextPaymentDate:   nextPayment,
		remainingPayments: totalPayments,
		totalPayments:     totalPayments,
	}, nil
}

// Getters with read locks
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

// Setters with write locks and validation
func (s *Subscription) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.name = name
	return nil
}

func (s *Subscription) SetCost(cost float64) error {
	if cost <= 0 {
		return fmt.Errorf("cost must be greater than 0")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cost = cost
	return nil
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

// calculateNextPaymentDate handles edge cases in date calculations
func (s *Subscription) calculateNextPaymentDate() time.Time {
	current := s.nextPaymentDate
	switch s.paymentFrequency {
	case FrequencyDaily:
		return current.AddDate(0, 0, 1)
	case FrequencyWeekly:
		return current.AddDate(0, 0, 7)
	case FrequencyMonthly:
		// Handle month end cases
		year, month, day := current.Date()
		month++
		if month > 12 {
			year++
			month = 1
		}
		// Adjust for months with fewer days
		lastDay := time.Date(year, month+1, 0, current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location()).Day()
		if day > lastDay {
			day = lastDay
		}
		return time.Date(year, month, day, current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location())
	case FrequencyYearly:
		// Handle leap year cases
		year, month, day := current.Date()
		if month == 2 && day == 29 && !isLeapYear(year+1) {
			day = 28
		}
		return time.Date(year+1, month, day, current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location())
	default:
		return current
	}
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func (s *Subscription) ProcessPayment() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.remainingPayments <= 0 {
		return fmt.Errorf("subscription has ended")
	}

	s.remainingPayments--
	s.nextPaymentDate = s.calculateNextPaymentDate()
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
