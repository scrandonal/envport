package cmd

// Manager defines the operations the CLI commands require from the store layer.
type Manager interface {
	Save(name string, vars map[string]string) error
	Load(name string) (map[string]string, error)
	List() ([]string, error)
	Delete(name string) error
	Rename(src, dst string) error
	Copy(src, dst string, overwrite bool) error
	Merge(src, dst string, overwrite bool) error
	Clone(src, dst string, overwrite bool) error
	Tag(name, tag string) error
	Untag(name, tag string) error
	Search(query string) ([]string, error)
	History(name string) ([]string, error)
	ClearHistory(name string) error
	IsLocked() (bool, error)
	ForceUnlock() error
}
