package alto

// AltoItemsChoice represents the variation of tag types available in ALTO format.
type AltoItemsChoice string

const (
	// AltoItemsChoiceLayoutTag indicates a layout-related tag element.
	AltoItemsChoiceLayoutTag AltoItemsChoice = "LayoutTag"
	// AltoItemsChoiceNamedEntityTag indicates a named entity tag element.
	AltoItemsChoiceNamedEntityTag AltoItemsChoice = "NamedEntityTag"
	// AltoItemsChoiceOtherTag indicates an other type of tag element.
	AltoItemsChoiceOtherTag AltoItemsChoice = "OtherTag"
	// AltoItemsChoiceRoleTag indicates a role-related tag element.
	AltoItemsChoiceRoleTag AltoItemsChoice = "RoleTag"
	// AltoItemsChoiceStructureTag indicates a structure-related tag element.
	AltoItemsChoiceStructureTag AltoItemsChoice = "StructureTag"
)

// String returns the XML string representation of this items choice type.
func (c AltoItemsChoice) String() string {
	return string(c)
}
