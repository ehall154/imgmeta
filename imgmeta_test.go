package imgmeta

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("writing fixture: %v", err)
	}
	return path
}

func TestReadJPEG(t *testing.T) {
	path := writeFixture(t, "photo.jpg", minimalJPEG(64, 48))

	m, err := Read(path)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if m.Format != "jpeg" || m.Width != 64 || m.Height != 48 {
		t.Fatalf("got %+v, want format jpeg, 64x48", m)
	}
	if m.Path != path {
		t.Fatalf("got path %q, want %q", m.Path, path)
	}
}

func TestReadPNG(t *testing.T) {
	path := writeFixture(t, "photo.png", minimalPNG(200, 100))

	m, err := Read(path)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if m.Format != "png" || m.Width != 200 || m.Height != 100 {
		t.Fatalf("got %+v, want format png, 200x100", m)
	}
}

func TestReadUnrecognizedFormat(t *testing.T) {
	path := writeFixture(t, "notes.txt", []byte("just some text, not an image"))

	if _, err := Read(path); err == nil {
		t.Fatal("expected error for unrecognized format, got nil")
	}
}

func TestReadMissingFile(t *testing.T) {
	if _, err := Read(filepath.Join(t.TempDir(), "does-not-exist.jpg")); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
