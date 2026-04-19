package store

import (
	"testing"
)

func newBookmarkStore(t *testing.T) string {
	t.Helper()
	s := newTempStore(t)
	// create a dummy snapshot so AddBookmark can find it
	snap := map[string]string{"KEY": "val"}
	if err := s.Save("mysnap", snap); err != nil {
		t.Fatal(err)
	}
	return s.base
}

func TestAddAndListBookmark(t *testing.T) {
	base := newBookmarkStore(t)
	if err := AddBookmark(base, "mysnap", "my favourite"); err != nil {
		t.Fatal(err)
	}
	bm, err := ListBookmarks(base)
	if err != nil {
		t.Fatal(err)
	}
	if len(bm) != 1 || bm[0].Name != "mysnap" || bm[0].Label != "my favourite" {
		t.Fatalf("unexpected bookmarks: %+v", bm)
	}
}

func TestAddBookmarkNotFound(t *testing.T) {
	base := newBookmarkStore(t)
	if err := AddBookmark(base, "ghost", "label"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAddBookmarkDuplicate(t *testing.T) {
	base := newBookmarkStore(t)
	_ = AddBookmark(base, "mysnap", "first")
	if err := AddBookmark(base, "mysnap", "second"); err == nil {
		t.Fatal("expected error on duplicate bookmark")
	}
}

func TestRemoveBookmark(t *testing.T) {
	base := newBookmarkStore(t)
	_ = AddBookmark(base, "mysnap", "label")
	if err := RemoveBookmark(base, "mysnap"); err != nil {
		t.Fatal(err)
	}
	bm, _ := ListBookmarks(base)
	if len(bm) != 0 {
		t.Fatal("expected empty bookmarks")
	}
}

func TestRemoveBookmarkNotFound(t *testing.T) {
	base := newBookmarkStore(t)
	if err := RemoveBookmark(base, "nobody"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
