package imgmeta

import (
	"bytes"
	"testing"
)

// minimalPNG builds the signature and IHDR chunk header/dimensions.
// parsePNG never reads past the width/height fields, so the trailing IHDR
// fields (bit depth, color type, ...) and the chunk CRC are omitted.
func minimalPNG(width, height uint32) []byte {
	buf := new(bytes.Buffer)
	buf.Write(pngSignature)
	buf.Write([]byte{0x00, 0x00, 0x00, 0x0D}) // chunk length (13)
	buf.WriteString("IHDR")
	writeUint32(buf, width)
	writeUint32(buf, height)
	return buf.Bytes()
}

func writeUint32(buf *bytes.Buffer, v uint32) {
	buf.Write([]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)})
}

func TestParsePNG(t *testing.T) {
	data := minimalPNG(200, 100)
	width, height, err := parsePNG(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("parsePNG returned error: %v", err)
	}
	if width != 200 || height != 100 {
		t.Fatalf("got %dx%d, want 200x100", width, height)
	}
}

func TestParsePNGWrongFirstChunk(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write(pngSignature)
	buf.Write([]byte{0x00, 0x00, 0x00, 0x04})
	buf.WriteString("gAMA")
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00})

	if _, _, err := parsePNG(buf); err == nil {
		t.Fatal("expected error when first chunk is not IHDR, got nil")
	}
}

func TestParsePNGTruncated(t *testing.T) {
	data := pngSignature[:4]
	if _, _, err := parsePNG(bytes.NewReader(data)); err == nil {
		t.Fatal("expected error for truncated PNG signature, got nil")
	}
}
