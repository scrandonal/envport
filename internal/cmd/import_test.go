package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportCmdDotenv(t *testing.T) {
	m := &mockManager{}
	cmd := newImportCmd(m)

	f := writeTempFile(t, "KEY1=value1\nKEY2=value2\n")
	cmd.SetArgs([]string{"mysnap", "--file", f})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "mysnap", m.savedName)
	assert.Equal(t, "value1", m.savedSnap.Vars["KEY1"])
	assert.Equal(t, "value2", m.savedSnap.Vars["KEY2"])
	assert.Contains(t, buf.String(), "Imported 2 variables")
}

func TestImportCmdShellFormat(t *testing.T) {
	m := &mockManager{}
	cmd := newImportCmd(m)

	f := writeTempFile(t, "export FOO=bar\nexport BAZ=qux\n")
	cmd.SetArgs([]string{"shellsnap", "--file", f, "--format", "shell"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "bar", m.savedSnap.Vars["FOO"])
	assert.Equal(t, "qux", m.savedSnap.Vars["BAZ"])
}

func TestImportCmdDestExists(t *testing.T) {
	m := &mockManager{listNames: []string{"existing"}}
	cmd := newImportCmd(m)

	f := writeTempFile(t, "KEY=val\n")
	cmd.SetArgs([]string{"existing", "--file", f})
	cmd.SetOut(&bytes.Buffer{})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestImportCmdOverwrite(t *testing.T) {
	m := &mockManager{listNames: []string{"existing"}}
	cmd := newImportCmd(m)

	f := writeTempFile(t, "KEY=val\n")
	cmd.SetArgs([]string{"existing", "--file", f, "--overwrite"})
	cmd.SetOut(&bytes.Buffer{})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "existing", m.savedName)
}

func TestImportCmdEmptyFile(t *testing.T) {
	m := &mockManager{}
	cmd := newImportCmd(m)

	f := writeTempFile(t, "# just a comment\n")
	cmd.SetArgs([]string{"snap", "--file", f})
	cmd.SetOut(&bytes.Buffer{})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no valid environment variables")
}

func TestImportCmdRequiresArg(t *testing.T) {
	m := &mockManager{}
	cmd := newImportCmd(m)
	cmd.SetArgs([]string{})
	cmd.SetOut(&bytes.Buffer{})
	err := cmd.Execute()
	assert.Error(t, err)
}

// TestImportCmdMissingFile verifies that an error is returned when the
// specified input file does not exist.
func TestImportCmdMissingFile(t *testing.T) {
	m := &mockManager{}
	cmd := newImportCmd(m)

	cmd.SetArgs([]string{"snap", "--file", "/nonexistent/path/input.env"})
	cmd.SetOut(&bytes.Buffer{})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Nil(t, m.savedSnap)
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "input.env")
	require.NoError(t, os.WriteFile(p, []byte(content), 0600))
	return p
}
