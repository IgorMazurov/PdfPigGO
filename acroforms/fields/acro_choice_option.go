package fields

import "fmt"

// AcroChoiceOption represents an option in a choice field (combo box or list box).
type AcroChoiceOption struct {
	Index        int
	IsSelected   bool
	Name         string
	ExportValue  *string
	HasExportVal bool
}

// NewAcroChoiceOption creates a new AcroChoiceOption.
func NewAcroChoiceOption(index int, isSelected bool, name string, exportValue *string) *AcroChoiceOption {
	hasExport := exportValue != nil
	return &AcroChoiceOption{
		Index:        index,
		IsSelected:   isSelected,
		Name:         name,
		ExportValue:  exportValue,
		HasExportVal: hasExport,
	}
}

// String returns a string representation of the option.
func (o *AcroChoiceOption) String() string {
	return fmt.Sprintf("%d: %s (%t).", o.Index, o.Name, o.IsSelected)
}
