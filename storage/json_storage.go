package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"subscription-tracker/models"
	"sync"
	"time"
)

type subscriptionJSON struct {
	Name              string  `json:"name"`
	Cost              float64 `json:"cost"`
	PaymentFrequency  string  `json:"payment_frequency"`
	NextPaymentDate   string  `json:"next_payment_date"`
	RemainingPayments int     `json:"remaining_payments"`
	TotalPayments     int     `json:"total_payments"`
}

type JSONStorage struct {
	filePath      string
	subscriptions []*models.Subscription
	mutex         sync.RWMutex
}

func NewJSONStorage(filePath string) (*JSONStorage, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	storage := &JSONStorage{
		filePath:      filePath,
		subscriptions: make([]*models.Subscription, 0),
	}

	// Load existing data if file exists
	if _, err := os.Stat(filePath); err == nil {
		if err := storage.loadFromFile(); err != nil {
			return nil, fmt.Errorf("failed to load subscriptions: %v", err)
		}
	}

	return storage, nil
}

func (s *JSONStorage) loadFromFile() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var jsonSubs []subscriptionJSON
	if err := json.Unmarshal(data, &jsonSubs); err != nil {
		return err
	}

	s.subscriptions = make([]*models.Subscription, 0, len(jsonSubs))
	for _, jsonSub := range jsonSubs {
		// Parse the date string
		date, err := time.Parse(time.RFC3339, jsonSub.NextPaymentDate)
		if err != nil {
			return fmt.Errorf("invalid date format for subscription %s: %v", jsonSub.Name, err)
		}

		sub, err := models.NewSubscription(
			jsonSub.Name,
			jsonSub.Cost,
			jsonSub.PaymentFrequency,
			date,
			jsonSub.TotalPayments,
		)
		if err != nil {
			return fmt.Errorf("failed to create subscription from JSON: %v", err)
		}
		s.subscriptions = append(s.subscriptions, sub)
	}

	return nil
}

func (s *JSONStorage) saveToFile() error {
	jsonSubs := make([]subscriptionJSON, len(s.subscriptions))
	for i, sub := range s.subscriptions {
		jsonSubs[i] = subscriptionJSON{
			Name:              sub.Name(),
			Cost:              sub.Cost(),
			PaymentFrequency:  sub.PaymentFrequency(),
			NextPaymentDate:   sub.NextPaymentDate().Format(time.RFC3339),
			RemainingPayments: sub.RemainingPayments(),
			TotalPayments:     sub.TotalPayments(),
		}
	}

	data, err := json.MarshalIndent(jsonSubs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *JSONStorage) AddSubscription(sub *models.Subscription) error {
	if sub == nil {
		return fmt.Errorf("subscription cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check for duplicate names
	for _, existing := range s.subscriptions {
		if existing.Name() == sub.Name() {
			return fmt.Errorf("subscription with name '%s' already exists", sub.Name())
		}
	}

	s.subscriptions = append(s.subscriptions, sub)

	if err := s.saveToFile(); err != nil {
		// Remove the subscription if save fails
		s.subscriptions = s.subscriptions[:len(s.subscriptions)-1]
		return fmt.Errorf("failed to save subscription: %v", err)
	}

	return nil
}

func (s *JSONStorage) GetSubscriptions() []*models.Subscription {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return a copy of the subscriptions slice to prevent external modifications
	result := make([]*models.Subscription, len(s.subscriptions))
	copy(result, s.subscriptions)
	return result
}

func (s *JSONStorage) UpdateSubscription(name string, updatedSub *models.Subscription) error {
	if updatedSub == nil {
		return fmt.Errorf("updated subscription cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, sub := range s.subscriptions {
		if sub.Name() == name {
			// If the name is being changed, check for duplicates
			if name != updatedSub.Name() {
				for _, existing := range s.subscriptions {
					if existing.Name() == updatedSub.Name() {
						return fmt.Errorf("subscription with name '%s' already exists", updatedSub.Name())
					}
				}
			}

			// Store old subscription in case save fails
			oldSub := s.subscriptions[i]
			s.subscriptions[i] = updatedSub

			if err := s.saveToFile(); err != nil {
				// Restore old subscription if save fails
				s.subscriptions[i] = oldSub
				return fmt.Errorf("failed to save subscription update: %v", err)
			}

			return nil
		}
	}
	return fmt.Errorf("subscription with name '%s' not found", name)
}

func (s *JSONStorage) DeleteSubscription(name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, sub := range s.subscriptions {
		if sub.Name() == name {
			// Store subscription and index in case save fails
			oldSub := sub
			oldIndex := i

			// Remove the subscription
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)

			if err := s.saveToFile(); err != nil {
				// Restore subscription if save fails
				s.subscriptions = append(s.subscriptions[:oldIndex], append([]*models.Subscription{oldSub}, s.subscriptions[oldIndex:]...)...)
				return fmt.Errorf("failed to save subscription deletion: %v", err)
			}

			return nil
		}
	}
	return fmt.Errorf("subscription with name '%s' not found", name)
}
