package store

import (
	"testing"
)

func newFavoriteStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestAddAndListFavorite(t *testing.T) {
	s := newFavoriteStore(t)
	if err := s.Save("alpha", map[string]string{"K": "1"}); err != nil {
		t.Fatal(err)
	}
	if err := AddFavorite(s.Root, "alpha"); err != nil {
		t.Fatal(err)
	}
	favs, err := ListFavorites(s.Root)
	if err != nil {
		t.Fatal(err)
	}
	if len(favs) != 1 || favs[0] != "alpha" {
		t.Fatalf("expected [alpha], got %v", favs)
	}
}

func TestAddFavoriteNotFound(t *testing.T) {
	s := newFavoriteStore(t)
	if err := AddFavorite(s.Root, "ghost"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAddFavoriteDuplicate(t *testing.T) {
	s := newFavoriteStore(t)
	s.Save("alpha", map[string]string{"K": "1"})
	AddFavorite(s.Root, "alpha")
	AddFavorite(s.Root, "alpha")
	favs, _ := ListFavorites(s.Root)
	if len(favs) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(favs))
	}
}

func TestRemoveFavorite(t *testing.T) {
	s := newFavoriteStore(t)
	s.Save("alpha", map[string]string{"K": "1"})
	AddFavorite(s.Root, "alpha")
	if err := RemoveFavorite(s.Root, "alpha"); err != nil {
		t.Fatal(err)
	}
	favs, _ := ListFavorites(s.Root)
	if len(favs) != 0 {
		t.Fatalf("expected empty, got %v", favs)
	}
}

func TestRemoveFavoriteIdempotent(t *testing.T) {
	s := newFavoriteStore(t)
	if err := RemoveFavorite(s.Root, "ghost"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
