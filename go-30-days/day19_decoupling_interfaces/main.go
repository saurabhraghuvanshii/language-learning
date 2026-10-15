/*
DAY 19: DECOUPLING WITH CONSUMER-SIDE INTERFACES
RUN:  go run ./day19_decoupling_interfaces

THEORY:
"Accept interfaces, return structs." The code that USES a dependency (the
consumer) declares a small interface listing only the methods it calls. The
provider just exposes a concrete type with methods; it does not declare or
know about the interface. Because satisfaction is implicit, any number of
implementations fit: a real database, an in-memory map, a fake that fails on
purpose. The consumer never imports the provider's package, so the dependency
arrow points from the concrete code toward the abstraction.
This also keeps interfaces tiny: each consumer needs only the 1-3 methods it
uses, instead of one giant "UserRepository" with 15 methods.
REAL WORLD: HTTP handlers depend on a small Store interface; unit tests pass a
fake that returns errors, so you can test the error path without a database
or network.
EXPLAIN-IT ANSWER: Why define interfaces on the consumer side?
Only the consumer knows which methods it actually needs, so only it can
write the smallest correct interface. Defining it there removes the import of
the provider (no cycles, no coupling), lets you add fakes in tests, and lets
you swap implementations without editing the consumer. An interface declared
next to the implementation tends to grow big and force every user to depend
on the whole thing.
PITFALLS:
  - Declaring a big interface "just in case" or one per struct (Java style).
  - Returning interfaces from constructors; return the concrete type instead.
  - Forgetting that the fake must be swappable: depend on the interface in
    the struct field, not on *MemoryStore.

EXERCISES:
 1. Add a method Delete(id) to the UserStore interface; fix both MemoryStore
    and FailingStore until it compiles.
 2. Add a Rename(id, name) method on UserService that uses Find then Save.
 3. Make FailingStore fail only on the second Save call (add a counter) and
    show that Register of the first user succeeds while the second fails.
*/
package main

import (
	"errors"
	"fmt"
	"strings"
)

// User is the domain object.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserStore is declared HERE, next to the consumer UserService.
// It has only the two methods the service calls.
type UserStore interface {
	Save(u User) error
	Find(id int) (User, error)
}

// ErrNotFound is returned by stores when no user has that id.
var ErrNotFound = errors.New("user not found")

// UserService contains the business rules. It holds the interface, not a concrete type.
type UserService struct {
	store  UserStore
	nextID int
}

// NewUserService takes the dependency as an argument (dependency injection).
func NewUserService(store UserStore) *UserService {
	return &UserService{store: store, nextID: 1}
}

// Register validates the input and saves a new user.
func (s *UserService) Register(name, email string) (User, error) {
	if strings.TrimSpace(name) == "" {
		return User{}, errors.New("name is required")
	}
	if !strings.Contains(email, "@") {
		return User{}, fmt.Errorf("invalid email %q", email)
	}
	u := User{ID: s.nextID, Name: name, Email: strings.ToLower(email)}
	if err := s.store.Save(u); err != nil {
		// Wrap with %w so the caller can still inspect the cause.
		return User{}, fmt.Errorf("register %s: %w", name, err)
	}
	s.nextID++ // only consume the id after a successful save
	return u, nil
}

// Greeting looks a user up and builds a message.
func (s *UserService) Greeting(id int) (string, error) {
	u, err := s.store.Find(id)
	if err != nil {
		return "", fmt.Errorf("greeting: %w", err)
	}
	return "Hello, " + u.Name + "!", nil
}

// MemoryStore is the real in-memory implementation. It never mentions UserStore.
type MemoryStore struct {
	users map[int]User
}

// NewMemoryStore returns the concrete type (accept interfaces, return structs).
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{users: map[int]User{}}
}

// Save stores the user by id.
func (m *MemoryStore) Save(u User) error {
	m.users[u.ID] = u
	return nil
}

// Find returns the user or ErrNotFound.
func (m *MemoryStore) Find(id int) (User, error) {
	u, ok := m.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

// FailingStore is a test fake: it fails every call with a chosen error.
type FailingStore struct {
	Err error
}

// Save always fails.
func (f FailingStore) Save(User) error { return f.Err }

// Find always fails.
func (f FailingStore) Find(int) (User, error) { return User{}, f.Err }

func main() {
	fmt.Println("== Real in-memory store")
	svc := NewUserService(NewMemoryStore()) // *MemoryStore satisfies UserStore
	u, err := svc.Register("Priya Shah", "Priya@Example.com")
	fmt.Printf("registered: %+v err=%v\n", u, err)
	msg, err := svc.Greeting(u.ID)
	fmt.Println("greeting:", msg, "err:", err)
	_, err = svc.Greeting(99)
	fmt.Println("unknown user:", err, "| is ErrNotFound:", errors.Is(err, ErrNotFound))

	fmt.Println("== Validation errors (no store involved)")
	_, err = svc.Register("", "x@y.com")
	fmt.Println("empty name:", err)
	_, err = svc.Register("Bob", "not-an-email")
	fmt.Println("bad email:", err)

	fmt.Println("== Failure test with the fake store")
	boom := errors.New("disk full")
	failing := NewUserService(FailingStore{Err: boom}) // FailingStore (value) satisfies UserStore too
	_, err = failing.Register("Carol", "carol@example.com")
	fmt.Println("error:", err)
	fmt.Println("is disk full:", errors.Is(err, boom))
	fmt.Println("id not consumed after failure:", failing.nextID == 1)
}
