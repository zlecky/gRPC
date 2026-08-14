package store

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrInvalidArgument = errors.New("invalid argument")
)

// User is the persistence model.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemoryStore is a thread-safe in-memory user store for demos and tests.
type MemoryStore struct {
	mu     sync.RWMutex
	users  map[string]*User
	emails map[string]string // email -> id
	seq    int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:  make(map[string]*User),
		emails: make(map[string]string),
	}
}

func (s *MemoryStore) Create(name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, ErrInvalidArgument
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.emails[email]; ok {
		return nil, ErrEmailExists
	}

	s.seq++
	now := time.Now().UTC()
	u := &User{
		ID:        formatID(s.seq),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.users[u.ID] = u
	s.emails[email] = u.ID
	return clone(u), nil
}

func (s *MemoryStore) Get(id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidArgument
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(u), nil
}

func (s *MemoryStore) List(pageSize int, pageToken string) ([]*User, string, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.users))
	for id := range s.users {
		ids = append(ids, id)
	}
	sortIDs(ids)

	start := 0
	if pageToken != "" {
		found := false
		for i, id := range ids {
			if id == pageToken {
				start = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, "", ErrInvalidArgument
		}
	}

	if start >= len(ids) {
		return []*User{}, "", nil
	}

	end := start + pageSize
	if end > len(ids) {
		end = len(ids)
	}

	out := make([]*User, 0, end-start)
	for _, id := range ids[start:end] {
		out = append(out, clone(s.users[id]))
	}

	next := ""
	if end < len(ids) {
		next = ids[end-1]
	}
	return out, next, nil
}

func (s *MemoryStore) Update(id, name, email string) (*User, error) {
	if id == "" || name == "" || email == "" {
		return nil, ErrInvalidArgument
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	if other, ok := s.emails[email]; ok && other != id {
		return nil, ErrEmailExists
	}

	delete(s.emails, u.Email)
	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now().UTC()
	s.emails[email] = id
	return clone(u), nil
}

func (s *MemoryStore) Delete(id string) error {
	if id == "" {
		return ErrInvalidArgument
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.emails, u.Email)
	delete(s.users, id)
	return nil
}

func (s *MemoryStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users)
}

func clone(u *User) *User {
	cp := *u
	return &cp
}

func formatID(n int) string {
	return "usr_" + itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

func sortIDs(ids []string) {
	for i := 1; i < len(ids); i++ {
		j := i
		for j > 0 && ids[j-1] > ids[j] {
			ids[j-1], ids[j] = ids[j], ids[j-1]
			j--
		}
	}
}
