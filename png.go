package imgmeta

import (
	"encoding/binary"
	"errors"
	"io"
)

// parsePNG reads the 8-byte PNG signature and the first chunk, which the
// spec requires to be IHDR and to hold the image dimensions.
func parsePNG(r io.Reader) (width, height int, err error) {
	var sig [8]byte
	if _, err := io.ReadFull(r, sig[:]); err != nil {
		return 0, 0, err
	}

	var chunkHeader [8]byte // 4 bytes length + 4 bytes chunk type
	if _, err := io.ReadFull(r, chunkHeader[:]); err != nil {
		return 0, 0, err
	}
	if string(chunkHeader[4:8]) != "IHDR" {
		return 0, 0, errors.New("first PNG chunk is not IHDR")
	}

	var ihdr [8]byte // 4 bytes width + 4 bytes height
	if _, err := io.ReadFull(r, ihdr[:]); err != nil {
		return 0, 0, err
	}
	width = int(binary.BigEndian.Uint32(ihdr[0:4]))
	height = int(binary.BigEndian.Uint32(ihdr[4:8]))
	return width, height, nil
}
