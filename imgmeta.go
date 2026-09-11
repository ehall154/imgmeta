// Package imgmeta reads structural metadata (format, dimensions) directly
// from image file bytes, without decoding pixel data.
package imgmeta

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Metadata describes what was found in a single image file.
type Metadata struct {
	Path   string `json:"path"`
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

var pngSignature = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

// Read opens the file at path and extracts its metadata. The format is
// determined by the file's magic bytes, not its extension.
func Read(path string) (*Metadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	magic := make([]byte, 8)
	n, err := io.ReadFull(f, magic)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("imgmeta: %s: %w", path, err)
	}
	magic = magic[:n]
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("imgmeta: %s: %w", path, err)
	}

	var format string
	var width, height int

	switch {
	case len(magic) >= 2 && magic[0] == 0xFF && magic[1] == 0xD8:
		format = "jpeg"
		width, height, err = parseJPEG(f)
	case len(magic) >= 8 && bytes.Equal(magic, pngSignature):
		format = "png"
		width, height, err = parsePNG(f)
	default:
		return nil, fmt.Errorf("imgmeta: %s: unrecognized image format", path)
	}
	if err != nil {
		return nil, fmt.Errorf("imgmeta: %s: %w", path, err)
	}

	return &Metadata{Path: path, Format: format, Width: width, Height: height}, nil
}
