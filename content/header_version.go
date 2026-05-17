package content

import "fmt"

// HeaderVersion holds the PDF header version information parsed from a file's first line.
type HeaderVersion struct {
	// Version is the numeric version value (e.g., 1.4).
	Version float64

	// VersionString is the raw version string as found in the file (e.g., "1.4").
	VersionString string

	// OffsetInFile is the byte offset from the start of the file to the start of the version comment.
	OffsetInFile int64
}

// NewHeaderVersion creates a new HeaderVersion with the given parameters.
func NewHeaderVersion(version float64, versionString string, offsetInFile int64) (*HeaderVersion, error) {
	if offsetInFile < 0 {
		return nil, fmt.Errorf("invalid offset for header version, must be positive. Got: %d", offsetInFile)
	}

	return &HeaderVersion{
		Version:       version,
		VersionString: versionString,
		OffsetInFile:  offsetInFile,
	}, nil
}

// String returns the string representation of the header version.
func (h *HeaderVersion) String() string {
	return fmt.Sprintf("Version: %s", h.VersionString)
}

var _ fmt.Stringer = (*HeaderVersion)(nil)
