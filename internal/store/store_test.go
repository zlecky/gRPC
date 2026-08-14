package store

import (
	"errors"
	"testing"
)

func TestMemoryStoreCRUD(t *testing.T) {
	s := NewMemoryStore()

	u, err := s.Create("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u.ID == "" || u.Name != "Alice" {
		t.Fatalf("unexpected user: %+v", u)
	}

	got, err := s.Get(u.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Email != "alice@example.com" {
		t.Fatalf("email mismatch: %s", got.Email)
	}

	updated, err := s.Update(u.ID, "Alice W", "alice.w@example.com")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "Alice W" || updated.Email != "alice.w@example.com" {
		t.Fatalf("unexpected update: %+v", updated)
	}

	if err := s.Delete(u.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get(u.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

func TestMemoryStoreDuplicateEmail(t *testing.T) {
	s := NewMemoryStore()
	if _, err := s.Create("A", "a@example.com"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create("B", "a@example.com"); !errors.Is(err, ErrEmailExists) {
		t.Fatalf("expected ErrEmailExists, got %v", err)
	}
}

func TestMemoryStoreListPagination(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 5; i++ {
		email := "u" + itoa(i+1) + "@example.com"
		if _, err := s.Create("U"+itoa(i+1), email); err != nil {
			t.Fatal(err)
		}
	}

	page1, next, err := s.List(2, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 2 || next == "" {
		t.Fatalf("page1=%d next=%q", len(page1), next)
	}

	page2, next2, err := s.List(2, next)
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 {
		t.Fatalf("page2=%d", len(page2))
	}

	page3, next3, err := s.List(2, next2)
	if err != nil {
		t.Fatal(err)
	}
	if len(page3) != 1 || next3 != "" {
		t.Fatalf("page3=%d next=%q", len(page3), next3)
	}
}

func TestMemoryStoreValidation(t *testing.T) {
	s := NewMemoryStore()
	if _, err := s.Create("", "a@example.com"); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
	if _, err := s.Get(""); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}
