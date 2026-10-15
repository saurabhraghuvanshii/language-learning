// Package storage provides byte storage. NOTE: it does not import package
// inventory, yet *Memory satisfies inventory.Store. Go interfaces are satisfied
// implicitly, so storage only needs methods with the right signatures, and the
// signatures use only built-in types (string, []byte, error), so no shared
// import is needed either.
package storage

import (
	"errors"
	"sort"
)

// ErrNotFound is exported so callers can compare with errors.Is.
var ErrNotFound = errors.New("storage: key not found")

// Memory is an in-memory key/value store.
type Memory struct {
	data map[string][]byte // unexported: callers must use the methods
}

// NewMemory builds a ready-to-use store (maps must be made before writing).
func NewMemory() *Memory {
	return &Memory{data: make(map[string][]byte)}
}

// Put stores a copy of value under key.
func (m *Memory) Put(key string, value []byte) error {
	m.data[key] = append([]byte(nil), value...) // copy so the caller cannot mutate our data later
	return nil
}

// Get returns the value for key or ErrNotFound.
func (m *Memory) Get(key string) ([]byte, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// Keys returns all keys sorted (unexported helper would be "keys"; this one is public).
func (m *Memory) Keys() []string {
	out := make([]string, 0, len(m.data))
	for k := range m.data {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// helper is unexported: invisible outside this package. Uncommenting a call to
// storage.helper in main.go gives: "undefined: storage.helper" / "cannot refer to unexported name".
func helper() {}
