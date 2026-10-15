// Package inventory holds the business logic. It DEFINES the interface it needs
// (Store) on the consumer side. It does not import package storage.
package inventory

import (
	"encoding/json"
	"fmt"

	"go30/day18_packages_layout/internal/idgen" // allowed: we are inside day18_packages_layout/
)

// Store is the consumer-side interface: only the two methods inventory calls.
// Any type with these methods (storage.Memory, a fake, a database) fits.
type Store interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

// Item is exported (capital I) with exported fields so encoding/json can see them.
type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}

// Service is exported, but its fields are not: users must go through New.
type Service struct {
	store Store
	ids   *idgen.Generator
}

// New receives the Store interface (dependency injection). It never names a
// concrete storage type.
func New(s Store) *Service {
	return &Service{store: s, ids: idgen.New("SKU")}
}

// Add creates an item with a new id and saves it as JSON.
func (s *Service) Add(name string, stock int, price float64) (Item, error) {
	it := Item{ID: s.ids.Next(), Name: name, Stock: stock, Price: price}
	b, err := json.Marshal(it)
	if err != nil {
		return Item{}, fmt.Errorf("encode item: %w", err)
	}
	if err := s.store.Put(it.ID, b); err != nil {
		return Item{}, fmt.Errorf("save item %s: %w", it.ID, err)
	}
	return it, nil
}

// Find loads an item by id.
func (s *Service) Find(id string) (Item, error) {
	b, err := s.store.Get(id)
	if err != nil {
		return Item{}, fmt.Errorf("find %s: %w", id, err) // keeps storage.ErrNotFound visible to errors.Is
	}
	var it Item
	if err := json.Unmarshal(b, &it); err != nil {
		return Item{}, fmt.Errorf("decode item %s: %w", id, err)
	}
	return it, nil
}

// value is unexported: only package inventory can call it.
func value(it Item) float64 { return float64(it.Stock) * it.Price }

// TotalValue is exported and uses the unexported helper.
func (s *Service) TotalValue(ids []string) (float64, error) {
	total := 0.0
	for _, id := range ids {
		it, err := s.Find(id)
		if err != nil {
			return 0, err
		}
		total += value(it)
	}
	return total, nil
}
