package page

// PageXmlTableCellRole represents data for a region that takes on the role of
// a table cell within a parent table region in PAGE XML documents.
type PageXmlTableCellRole struct {
	// RowIndex is the cell position in table starting with row 0.
	RowIndex int `xml:"rowIndex,attr"`

	// ColumnIndex is the cell position in table starting with column 0.
	ColumnIndex int `xml:"columnIndex,attr"`

	// RowSpan is the number of rows the cell spans (optional; default is 1).
	RowSpan *int `xml:"rowSpan,attr,omitempty"`

	// ColSpan is the number of columns the cell spans (optional; default is 1).
	ColSpan *int `xml:"colSpan,attr,omitempty"`

	// Header indicates whether the cell is a column or row header.
	Header *bool `xml:"header,attr,omitempty"`
}
