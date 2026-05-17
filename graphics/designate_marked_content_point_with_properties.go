// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// DesignateMarkedContentPointWithProperties designates a single marked-content point in the content stream with an associated property list.
type DesignateMarkedContentPointWithProperties struct {
	Name                   *tokens.NameToken
	PropertyDictionaryName *tokens.NameToken
	Properties             *tokens.DictionaryToken
}

const designateMarkedContentPointWithPropertiesSymbol = "DP"

// NewDesignateMarkedContentPointWithPropertiesName creates a new DesignateMarkedContentPointWithProperties with a property dictionary name.
func NewDesignateMarkedContentPointWithPropertiesName(name, propertyDictionaryName *tokens.NameToken) (*DesignateMarkedContentPointWithProperties, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	if propertyDictionaryName == nil {
		return nil, fmt.Errorf("property dictionary name cannot be nil")
	}
	return &DesignateMarkedContentPointWithProperties{
		Name:                   name,
		PropertyDictionaryName: propertyDictionaryName,
	}, nil
}

// NewDesignateMarkedContentPointWithPropertiesDict creates a new DesignateMarkedContentPointWithProperties with inline properties.
func NewDesignateMarkedContentPointWithPropertiesDict(name *tokens.NameToken, properties *tokens.DictionaryToken) (*DesignateMarkedContentPointWithProperties, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	if properties == nil {
		return nil, fmt.Errorf("properties cannot be nil")
	}
	return &DesignateMarkedContentPointWithProperties{
		Name:       name,
		Properties: properties,
	}, nil
}

// Operator returns the operator symbol for this operation.
func (o DesignateMarkedContentPointWithProperties) Operator() string {
	return designateMarkedContentPointWithPropertiesSymbol
}

// Run executes the operation on the given context. This is a no-op as per PDF spec.
func (o *DesignateMarkedContentPointWithProperties) Run(ctx OperationContext) {
}

// Write writes the operator to the given writer.
func (o DesignateMarkedContentPointWithProperties) Write(w io.Writer) error {
	var propPart string
	if o.PropertyDictionaryName != nil {
		propPart = o.PropertyDictionaryName.String()
	} else if o.Properties != nil {
		propPart = o.Properties.String()
	}

	if _, err := fmt.Fprintf(w, "%s %s %s\n", o.Name, propPart, designateMarkedContentPointWithPropertiesSymbol); err != nil {
		return err
	}
	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (DesignateMarkedContentPointWithProperties) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o DesignateMarkedContentPointWithProperties) String() string {
	var propPart string
	if o.PropertyDictionaryName != nil {
		propPart = o.PropertyDictionaryName.String()
	} else if o.Properties != nil {
		propPart = o.Properties.String()
	}
	return fmt.Sprintf("%s %s %s", o.Name, propPart, designateMarkedContentPointWithPropertiesSymbol)
}
