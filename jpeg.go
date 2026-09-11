package imgmeta

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// parseJPEG walks the marker segments of a JPEG stream until it finds a
// start-of-frame marker, which carries the pixel dimensions. It never
// decodes the compressed scan data.
func parseJPEG(r io.Reader) (width, height int, err error) {
	br := bufio.NewReader(r)

	var soi [2]byte
	if _, err := io.ReadFull(br, soi[:]); err != nil {
		return 0, 0, err
	}
	if soi[0] != 0xFF || soi[1] != 0xD8 {
		return 0, 0, errors.New("missing JPEG start-of-image marker")
	}

	for {
		marker, err := nextMarker(br)
		if err != nil {
			return 0, 0, err
		}

		// Markers with no payload: restart markers, TEM, and SOI/EOI.
		if marker == 0x01 || marker == 0xD9 || (marker >= 0xD0 && marker <= 0xD7) {
			continue
		}

		if isStartOfFrame(marker) {
			var seg [7]byte
			if _, err := io.ReadFull(br, seg[:]); err != nil {
				return 0, 0, err
			}
			height = int(seg[3])<<8 | int(seg[4])
			width = int(seg[5])<<8 | int(seg[6])
			return width, height, nil
		}

		var lenBuf [2]byte
		if _, err := io.ReadFull(br, lenBuf[:]); err != nil {
			return 0, 0, err
		}
		length := int(lenBuf[0])<<8 | int(lenBuf[1])
		if length < 2 {
			return 0, 0, fmt.Errorf("invalid marker segment length %d", length)
		}
		if _, err := br.Discard(length - 2); err != nil {
			return 0, 0, err
		}
	}
}

// nextMarker scans forward to the next 0xFF marker byte and returns the
// byte that follows it, skipping padding fill bytes (0xFF) and stuffed
// zero bytes that appear inside entropy-coded data.
func nextMarker(br *bufio.Reader) (byte, error) {
	for {
		b, err := br.ReadByte()
		if err != nil {
			return 0, err
		}
		if b != 0xFF {
			continue
		}
		var m byte
		for {
			m, err = br.ReadByte()
			if err != nil {
				return 0, err
			}
			if m != 0xFF {
				break
			}
		}
		if m == 0x00 {
			continue
		}
		return m, nil
	}
}

func isStartOfFrame(marker byte) bool {
	switch marker {
	case 0xC0, 0xC1, 0xC2, 0xC3, 0xC5, 0xC6, 0xC7, 0xC9, 0xCA, 0xCB, 0xCD, 0xCE, 0xCF:
		return true
	default:
		return false
	}
}
