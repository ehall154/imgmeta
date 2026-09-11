# imgmeta

A Go library for reading basic image metadata straight from file bytes,
without decoding pixel data. Given a JPEG or PNG file, it tells you the
format and dimensions by walking the file's own structure (marker segments
for JPEG, the IHDR chunk for PNG).

The reason this exists instead of just calling `image.DecodeConfig`: I
wanted metadata reading that stays useful outside a decoder, with its own
output shapes, so it can grow into a real EXIF/metadata toolkit without
being tied to Go's `image` package registration and its side-effecting
blank imports.

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/ehall154/imgmeta"
)

func main() {
	m, err := imgmeta.Read("photo.jpg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(m)
}
```

Human-readable output:

```
photo.jpg: jpeg 4032x3024
```

## JSON output

Every entry point supports a JSON mode alongside the human-readable one,
so a caller (a CLI, a batch job, whatever) can pick the shape it needs
without reimplementing formatting:

```go
m, err := imgmeta.Read("photo.jpg")
if err != nil {
	// handle err
}

// asJSON would typically come from a --json flag in a calling program.
if err := imgmeta.Write(os.Stdout, m, asJSON); err != nil {
	// handle err
}
```

With `asJSON == true`:

```json
{
  "path": "photo.jpg",
  "format": "jpeg",
  "width": 4032,
  "height": 3024
}
```

## What's supported now

- JPEG: dimensions via the first start-of-frame marker.
- PNG: dimensions via the IHDR chunk.

## What's not here yet

- EXIF tag extraction (orientation, camera model, GPS, timestamps).
- Other formats (GIF, WebP, TIFF, HEIC).
- A command-line front end.

## License

MIT, see LICENSE.
