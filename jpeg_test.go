package imgmeta

import (
	"bytes"
	"testing"
)

// minimalJPEG builds a JPEG byte stream with an APP0 (JFIF) segment ahead of
// the start-of-frame marker, so the test also exercises the marker-skipping
// path in nextMarker/parseJPEG, not just the SOF0 read itself.
func minimalJPEG(width, height int) []byte {
	return []byte{
		0xFF, 0xD8, // SOI
		0xFF, 0xE0, // APP0
		0x00, 0x10, // segment length (16)
		'J', 'F', 'I', 'F', 0x00,
		0x01, 0x01, // version
		0x00,       // density units
		0x00, 0x01, // x density
		0x00, 0x01, // y density
		0x00, // thumbnail width
		0x00, // thumbnail height
		0xFF, 0xC0, // SOF0
		0x00, 0x11, // segment length (17)
		0x08, // sample precision
		byte(height >> 8), byte(height),
		byte(width >> 8), byte(width),
	}
}

func TestParseJPEG(t *testing.T) {
	data := minimalJPEG(64, 48)
	width, height, err := parseJPEG(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("parseJPEG returned error: %v", err)
	}
	if width != 64 || height != 48 {
		t.Fatalf("got %dx%d, want 64x48", width, height)
	}
}

func TestParseJPEGMissingSOI(t *testing.T) {
	data := []byte{0x00, 0x00, 0xFF, 0xC0}
	if _, _, err := parseJPEG(bytes.NewReader(data)); err == nil {
		t.Fatal("expected error for missing start-of-image marker, got nil")
	}
}

func TestParseJPEGTruncated(t *testing.T) {
	data := []byte{0xFF, 0xD8, 0xFF, 0xC0, 0x00, 0x11, 0x08}
	if _, _, err := parseJPEG(bytes.NewReader(data)); err == nil {
		t.Fatal("expected error for truncated SOF0 segment, got nil")
	}
}
