package imgmeta

import (
	"encoding/json"
	"fmt"
	"io"
)

// String renders the metadata as a single human-readable line.
func (m *Metadata) String() string {
	return fmt.Sprintf("%s: %s %dx%d", m.Path, m.Format, m.Width, m.Height)
}

// JSON renders the metadata as an indented JSON object.
func (m *Metadata) JSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// Write formats m to w in whichever mode the caller wants. This is the
// hook a CLI wraps around a --json flag: pass the flag value straight
// through as asJSON.
func Write(w io.Writer, m *Metadata, asJSON bool) error {
	if asJSON {
		b, err := m.JSON()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(b))
		return err
	}
	_, err := fmt.Fprintln(w, m.String())
	return err
}
