package store

import (
	"testing"
)

func TestSaveAndLoadTemplate(t *testing.T) {
	s := newTempStore(t)
	tmpl := Template{
		Name:     "base",
		Keys:     []string{"DB_HOST", "DB_PORT"},
		Defaults: map[string]string{"DB_PORT": "5432"},
	}
	if err := s.SaveTemplate(tmpl); err != nil {
		t.Fatal(err)
	}
	got, err := s.LoadTemplate("base")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != tmpl.Name || len(got.Keys) != 2 || got.Defaults["DB_PORT"] != "5432" {
		t.Fatalf("unexpected template: %+v", got)
	}
}

func TestLoadTemplateMissing(t *testing.T) {
	s := newTempStore(t)
	_, err := s.LoadTemplate("nope")
	if err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestDeleteTemplate(t *testing.T) {
	s := newTempStore(t)
	tmpl := Template{Name: "tmp", Keys: []string{"A"}}
	s.SaveTemplate(tmpl)
	if err := s.DeleteTemplate("tmp"); err != nil {
		t.Fatal(err)
	}
	_, err := s.LoadTemplate("tmp")
	if err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound after delete")
	}
}

func TestDeleteTemplateNotFound(t *testing.T) {
	s := newTempStore(t)
	if err := s.DeleteTemplate("ghost"); err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestListTemplates(t *testing.T) {
	s := newTempStore(t)
	names, _ := s.ListTemplates()
	if len(names) != 0 {
		t.Fatal("expected empty list")
	}
	s.SaveTemplate(Template{Name: "a", Keys: []string{"X"}})
	s.SaveTemplate(Template{Name: "b", Keys: []string{"Y"}})
	names, err := s.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(names))
	}
}
