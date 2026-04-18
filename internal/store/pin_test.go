package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPinAndPinned(t *testing.T) {
	s := newTempStore(t)
	snap := makeSnapshot(map[string]string{"KEY": "val"})
	require.NoError(t, s.Save("mysnap", snap))

	require.NoError(t, s.Pin("mysnap"))

	name, err := s.Pinned()
	require.NoError(t, err)
	assert.Equal(t, "mysnap", name)
}

func TestPinNotFound(t *testing.T) {
	s := newTempStore(t)
	err := s.Pin("ghost")
	assert.Error(t, err)
}

func TestPinnedMissing(t *testing.T) {
	s := newTempStore(t)
	_, err := s.Pinned()
	assert.ErrorIs(t, err, ErrNotPinned)
}

func TestUnpin(t *testing.T) {
	s := newTempStore(t)
	snap := makeSnapshot(map[string]string{"A": "1"})
	require.NoError(t, s.Save("snap", snap))
	require.NoError(t, s.Pin("snap"))

	require.NoError(t, s.Unpin())

	_, err := s.Pinned()
	assert.ErrorIs(t, err, ErrNotPinned)
}

func TestUnpinWhenNotPinned(t *testing.T) {
	s := newTempStore(t)
	err := s.Unpin()
	assert.ErrorIs(t, err, ErrNotPinned)
}

func TestPinOverwrite(t *testing.T) {
	s := newTempStore(t)
	for _, n := range []string{"a", "b"} {
		require.NoError(t, s.Save(n, makeSnapshot(map[string]string{"X": n})))
	}
	require.NoError(t, s.Pin("a"))
	require.NoError(t, s.Pin("b"))

	name, err := s.Pinned()
	require.NoError(t, err)
	assert.Equal(t, "b", name)
}
