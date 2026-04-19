package store

import (
	"testing"
	"time"
)

func TestAppendAndReadAudit(t *testing.T) {
	s := newTempStore(t)
	entry := AuditEntry{Time: time.Now(), Operation: "save", Name: "dev", Detail: "created"}
	if err := s.AppendAudit(entry); err != nil {
		t.Fatal(err)
	}
	entries, err := s.ReadAudit()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Operation != "save" || entries[0].Name != "dev" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestAuditEmptyOnMissingFile(t *testing.T) {
	s := newTempStore(t)
	entries, err := s.ReadAudit()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(entries))
	}
}

func TestAuditClear(t *testing.T) {
	s := newTempStore(t)
	_ = s.AppendAudit(AuditEntry{Time: time.Now(), Operation: "delete", Name: "old"})
	if err := s.ClearAudit(); err != nil {
		t.Fatal(err)
	}
	entries, _ := s.ReadAudit()
	if len(entries) != 0 {
		t.Errorf("expected empty after clear, got %d", len(entries))
	}
}

func TestAuditClearIdempotent(t *testing.T) {
	s := newTempStore(t)
	if err := s.ClearAudit(); err != nil {
		t.Errorf("expected no error on missing file, got %v", err)
	}
}

func TestAuditMultipleEntries(t *testing.T) {
	s := newTempStore(t)
	ops := []string{"save", "load", "delete"}
	for _, op := range ops {
		_ = s.AppendAudit(AuditEntry{Time: time.Now(), Operation: op, Name: "snap"})
	}
	entries, err := s.ReadAudit()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}
