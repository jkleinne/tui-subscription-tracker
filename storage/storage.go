package storage

import (
	"fmt"
	"subscription-tracker/models"
	"sync"
)

type Storage interface {
	AddSubscription(sub *models.Subscription) error
	GetSubscriptions() []*models.Subscription
	UpdateSubscription(name string, updatedSub *models.Subscription) error
	DeleteSubscription(name string) error
}

type MemoryStorage struct {
	subscriptions []*models.Subscription
	mutex         sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		subscriptions: make([]*models.Subscription, 0),
	}
}

func (s *MemoryStorage) AddSubscription(sub *models.Subscription) error {
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
	return nil
}

func (s *MemoryStorage) GetSubscriptions() []*models.Subscription {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return a copy of the subscriptions slice to prevent external modifications
	result := make([]*models.Subscription, len(s.subscriptions))
	copy(result, s.subscriptions)
	return result
}

func (s *MemoryStorage) UpdateSubscription(name string, updatedSub *models.Subscription) error {
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
			s.subscriptions[i] = updatedSub
			return nil
		}
	}
	return fmt.Errorf("subscription with name '%s' not found", name)
}

func (s *MemoryStorage) DeleteSubscription(name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, sub := range s.subscriptions {
		if sub.Name() == name {
			// Remove the subscription by slicing
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("subscription with name '%s' not found", name)
}
