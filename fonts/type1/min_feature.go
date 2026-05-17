package type1

import "fmt"

// MinFeature represents the Type1PrivateDictionary MinFeature entry which is required for compatibility
// and is required by the specification to have the value 16, 16.
type MinFeature struct {
	First  int
	Second int
}

// DefaultMinFeature is the required default value of MinFeature (16, 16).
var DefaultMinFeature = MinFeature{First: 16, Second: 16}

// NewMinFeature creates a new MinFeature with the given values.
func NewMinFeature(first, second int) MinFeature {
	return MinFeature{First: first, Second: second}
}

// String returns the string representation of MinFeature in PDF array format.
func (m MinFeature) String() string {
	return fmt.Sprintf("{%d %d}", m.First, m.Second)
}
