// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// BeginMarkedContentWithProperties begins a marked-content sequence with an associated property list
// terminated by a balancing EndMarkedContent operator.
type BeginMarkedContentWithProperties struct {
	Name                    *tokens.NameToken
	PropertyDictionaryName  *tokens.NameToken
	Properties              *tokens.DictionaryToken
}

const beginMarkedContentWithPropertiesSymbol = "BDC"

// NewBeginMarkedContentWithPropertiesName creates a new BeginMarkedContentWithProperties using a property dictionary name.
func NewBeginMarkedContentWithPropertiesName(name *tokens.NameToken, propertyDictionaryName *tokens.NameToken) (*BeginMarkedContentWithProperties, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	if propertyDictionaryName == nil {
		return nil, fmt.Errorf("propertyDictionaryName cannot be nil")
	}
	return &BeginMarkedContentWithProperties{
		Name:                   name,
		PropertyDictionaryName: propertyDictionaryName,
	}, nil
}

// NewBeginMarkedContentWithPropertiesDict creates a new BeginMarkedContentWithProperties using inline properties.
func NewBeginMarkedContentWithPropertiesDict(name *tokens.NameToken, properties *tokens.DictionaryToken) (*BeginMarkedContentWithProperties, error) {
	if name == nil {
		return nil, fmt.Errorf("name cannot be nil")
	}
	if properties == nil {
		return nil, fmt.Errorf("properties cannot be nil")
	}
	return &BeginMarkedContentWithProperties{
		Name:       name,
		Properties: properties,
	}, nil
}

// Operator returns the operator symbol for this operation.
func (o BeginMarkedContentWithProperties) Operator() string {
	return beginMarkedContentWithPropertiesSymbol
}

// Run executes the operation on the given context, starting a marked content section with properties.
func (o *BeginMarkedContentWithProperties) Run(ctx OperationContext) {
	ctx.BeginMarkedContent(o.Name, o.PropertyDictionaryName, o.Properties)
}

// Write writes the operator to the given writer.
func (o BeginMarkedContentWithProperties) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s ", o.Name); err != nil {
		return err
	}

	if o.PropertyDictionaryName != nil {
		if _, err := fmt.Fprintf(w, "%s", o.PropertyDictionaryName); err != nil {
			return err
		}
	} else if o.Properties != nil {
		if _, err := fmt.Fprintf(w, "%s", o.Properties); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, " %s\n", beginMarkedContentWithPropertiesSymbol); err != nil {
		return err
	}

	return nil
}

// IsGraphicsOp marks this type as implementing content.GraphicsStateOperation.
func (BeginMarkedContentWithProperties) IsGraphicsOp() {}

// String returns the string representation of this operation.
func (o BeginMarkedContentWithProperties) String() string {
	var propPart string
	if o.PropertyDictionaryName != nil {
		propPart = o.PropertyDictionaryName.String()
	} else if o.Properties != nil {
		propPart = o.Properties.String()
	}
	return fmt.Sprintf("%s %s %s", o.Name, propPart, beginMarkedContentWithPropertiesSymbol)
}
