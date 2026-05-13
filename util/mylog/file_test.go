package mylog

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestFileLogWriteAndReadLastNLines(t *testing.T) {
	dir := t.TempDir()
	oldLogger := stdFileLogger
	stdFileLogger = NewFileLogger(dir)
	defer func() {
		stdFileLogger = oldLogger
	}()

	WriteLogFile("access.log", "INFO", "first", F("uid", 1))
	WriteLogFile("access.log", "INFO", "second", F("uid", 2))
	WriteLogFile("access.log", "INFO", "third", F("uid", 3))

	lines, err := ReadLastNLogLines("access.log", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
	}

	var first map[string]any
	if err = json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatal(err)
	}
	if first["message"] != "second" {
		t.Fatalf("want second, got %v", first["message"])
	}
	if first["level"] != "INFO" {
		t.Fatalf("want INFO, got %v", first["level"])
	}
	if first["uid"].(float64) != 2 {
		t.Fatalf("want uid 2, got %v", first["uid"])
	}
	if first["caller"] == "" {
		t.Fatal("caller should not be empty")
	}
}

func TestNewFileLogger(t *testing.T) {
	log := NewFileLogger(t.TempDir())
	log.Write("custom.log", "CUSTOM", "hello", F("ok", true))

	lines, err := log.ReadLastNLines("custom.log", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 {
		t.Fatalf("want 1 line, got %d", len(lines))
	}

	var line map[string]any
	if err = json.Unmarshal([]byte(lines[0]), &line); err != nil {
		t.Fatal(err)
	}
	if line["level"] != "CUSTOM" {
		t.Fatalf("want CUSTOM, got %v", line["level"])
	}
	if line["ok"] != true {
		t.Fatalf("want ok true, got %v", line["ok"])
	}
}

func TestReadLastNLinesMissingFile(t *testing.T) {
	lines, err := ReadLastNLines(filepath.Join(t.TempDir(), "missing.log"), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("want empty lines, got %d", len(lines))
	}
}
