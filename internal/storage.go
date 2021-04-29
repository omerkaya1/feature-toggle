package internal

import (
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryStorage struct {
	// Imitates transaction functionality
	sync.RWMutex
	// Imitates DB
	store map[string]Feature
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		RWMutex: sync.RWMutex{},
		store:   make(map[string]Feature),
	}
}

func (ims *InMemoryStorage) Close() error {
	return nil
}

func (ims *InMemoryStorage) GetFeatureByName(_ context.Context, name string) (Feature, error) {
	ims.RLock()
	defer ims.RUnlock()
	if ft, ok := ims.store[name]; ok {
		return ft, nil
	}
	return Feature{}, errors.New("entry does not exist")
}

func contains(s []string, k string) bool {
	for i := range s {
		if s[i] == k {
			return true
		}
	}
	return false
}

func (ims *InMemoryStorage) GetUserFeaturesByStatus(ctx context.Context, customer string, active bool) ([]Feature, error) {
	ims.RLock()
	defer ims.RUnlock()
	var result []Feature
	for _, v := range ims.store {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if contains(v.CustomerIDs, customer) && v.Active == active {
				result = append(result, v)
			}
		}
	}
	return result, nil
}

func (ims *InMemoryStorage) GetUserFeatures(ctx context.Context, customer string) ([]Feature, error) {
	ims.RLock()
	defer ims.RUnlock()
	var result []Feature
	for _, v := range ims.store {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if contains(v.CustomerIDs, customer) {
				result = append(result, v)
			}
		}
	}
	return result, nil
}

func (ims *InMemoryStorage) GetFeatures(ctx context.Context) ([]Feature, error) {
	ims.RLock()
	defer ims.RUnlock()
	var result []Feature
	for _, v := range ims.store {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			result = append(result, v)
		}
	}
	return result, nil
}

func (ims *InMemoryStorage) CreateFeature(_ context.Context, inverted, active bool, displayName, techName, description string, customerIDs []string, expires time.Time) (int64, error) {
	// Check the presence of the feature in the storage
	ims.Lock()
	defer ims.Unlock()
	if _, ok := ims.store[techName]; ok {
		return 0, errors.New("conflicting tech name")
	}
	id := int64(len(ims.store)) + 1
	// Create a new entry in the "DB"
	ims.store[techName] = Feature{
		Inverted:      inverted,
		Active:        active,
		ID:            id,
		DisplayName:   displayName,
		TechnicalName: techName,
		Description:   description,
		CustomerIDs:   customerIDs,
		Expires:       expires,
	}
	return id, nil
}

// UpdateFeature .
// NOTE: following assumptions have been made:
//  - the techName cannot be changed, as it should be unique for each feature
//  - the only fields that seem reasonable to update are: expiration date, description, display name, status
func (ims *InMemoryStorage) UpdateFeature(_ context.Context, techName, description, displayName string, ed time.Time, status bool) error {
	ims.Lock()
	defer ims.Unlock()
	if v, ok := ims.store[techName]; ok {
		v.Description = description
		v.DisplayName = displayName
		v.Expires = ed
		v.Active = status
		ims.store[techName] = v
		return nil
	}
	return errors.New("entry does not exist")
}

func (ims *InMemoryStorage) DeleteFeature(_ context.Context, techName string) error {
	ims.Lock()
	defer ims.Unlock()
	if _, ok := ims.store[techName]; ok {
		delete(ims.store, techName)
		return nil
	}
	return errors.New("entry does not exist")
}
