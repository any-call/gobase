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
	if first["content"] != "second" {
		t.Fatalf("want second, got %v", first["content"])
	}
	if first["timestamp"] == "" {
		t.Fatal("timestamp should not be empty")
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
	defer log.Close()

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

func TestKeepLastNLines(t *testing.T) {
	log := NewFileLogger(t.TempDir())
	defer log.Close()

	log.Write("keep.log", "INFO", "first")
	log.Write("keep.log", "INFO", "second")
	log.Write("keep.log", "INFO", "third")

	if err := log.KeepLastNLines("keep.log", 2); err != nil {
		t.Fatal(err)
	}

	lines, err := log.ReadLastNLines("keep.log", 10)
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
	if first["content"] != "second" {
		t.Fatalf("want second, got %v", first["content"])
	}
}

func TestFileLoggerClose(t *testing.T) {
	log := NewFileLogger(t.TempDir())
	log.Write("close.log", "INFO", "first")

	if len(log.files) != 1 {
		t.Fatalf("want 1 cached file, got %d", len(log.files))
	}
	if err := log.Close(); err != nil {
		t.Fatal(err)
	}
	if len(log.files) != 0 {
		t.Fatalf("want empty cached files, got %d", len(log.files))
	}

	log.Write("close.log", "INFO", "second")
	lines, err := log.ReadLastNLines("close.log", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
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
