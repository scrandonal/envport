package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagCmdAdd(t *testing.T) {
	m := newMockManager()
	snap := makeSnapshot(map[string]string{"FOO": "bar"})
	m.snapshots["dev"] = snap

	root := newTestRoot(m)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"tag", "dev", "production"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "tagged \"dev\" with \"production\"")
	assert.Contains(t, m.snapshots["dev"].Tags, "production")
}

func TestTagCmdRemove(t *testing.T) {
	m := newMockManager()
	snap := makeSnapshot(map[string]string{"FOO": "bar"})
	snap.Tags = []string{"production", "v1"}
	m.snapshots["dev"] = snap

	root := newTestRoot(m)
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"tag", "--remove", "dev", "production"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "removed tag \"production\" from \"dev\"")
	assert.NotContains(t, m.snapshots["dev"].Tags, "production")
	assert.Contains(t, m.snapshots["dev"].Tags, "v1")
}

func TestTagCmdNotFound(t *testing.T) {
	m := newMockManager()
	root := newTestRoot(m)
	root.SetArgs([]string{"tag", "missing", "sometag"})

	err := root.Execute()
	assert.Error(t, err)
}

func TestTagCmdRequiresTwoArgs(t *testing.T) {
	m := newMockManager()
	root := newTestRoot(m)
	root.SetArgs([]string{"tag", "only-one"})

	err := root.Execute()
	assert.Error(t, err)
}
