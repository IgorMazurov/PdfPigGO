package colors

import "github.com/uglytoad/pdfpig/go/tokens"

// ResourceColorSpace represents a color space definition from a resource dictionary.
type ResourceColorSpace struct {
	// Name is the color space name.
	Name *tokens.NameToken

	// Data is the color space data.
	Data tokens.Token
}

// newResourceColorSpace creates a ResourceColorSpace with both name and data.
func newResourceColorSpace(name *tokens.NameToken, data tokens.Token) ResourceColorSpace {
	return ResourceColorSpace{
		Name: name,
		Data: data,
	}
}

// newResourceColorSpaceNameOnly creates a ResourceColorSpace with only the name (data is nil).
func newResourceColorSpaceNameOnly(name *tokens.NameToken) ResourceColorSpace {
	return ResourceColorSpace{
		Name: name,
	}
}
