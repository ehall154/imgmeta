package imgmeta

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMetadataString(t *testing.T) {
	m := &Metadata{Path: "photo.jpg", Format: "jpeg", Width: 64, Height: 48}
	want := "photo.jpg: jpeg 64x48"
	if got := m.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestMetadataJSON(t *testing.T) {
	m := &Metadata{Path: "photo.jpg", Format: "jpeg", Width: 64, Height: 48}
	b, err := m.JSON()
	if err != nil {
		t.Fatalf("JSON returned error: %v", err)
	}

	var got Metadata
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshaling JSON output: %v", err)
	}
	if got != *m {
		t.Fatalf("got %+v, want %+v", got, *m)
	}
}

func TestWrite(t *testing.T) {
	m := &Metadata{Path: "photo.jpg", Format: "jpeg", Width: 64, Height: 48}

	var plain bytes.Buffer
	if err := Write(&plain, m, false); err != nil {
		t.Fatalf("Write (plain) returned error: %v", err)
	}
	if want := m.String() + "\n"; plain.String() != want {
		t.Fatalf("got %q, want %q", plain.String(), want)
	}

	var asJSON bytes.Buffer
	if err := Write(&asJSON, m, true); err != nil {
		t.Fatalf("Write (json) returned error: %v", err)
	}
	var got Metadata
	if err := json.Unmarshal(asJSON.Bytes(), &got); err != nil {
		t.Fatalf("unmarshaling JSON output: %v", err)
	}
	if got != *m {
		t.Fatalf("got %+v, want %+v", got, *m)
	}
}
