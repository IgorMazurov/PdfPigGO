package page

// PageXmlRoles contains role metadata for regions in a PAGE XML document.
type PageXmlRoles struct {
	// TableCellRole is the data for a region that takes on the role of a table
	// cell within a parent table region.
	TableCellRole *PageXmlTableCellRole `xml:"tableCellRole,omitempty"`
}
