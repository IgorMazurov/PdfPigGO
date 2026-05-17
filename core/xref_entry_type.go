package core

// XrefEntryType indicates where an object is located in the cross-reference section.
type XrefEntryType byte

const (
	// XrefEntryTypeFree represents a free (deleted) object.
	XrefEntryTypeFree XrefEntryType = 0

	// XrefEntryTypeFile represents an object located as a regular object in the file.
	XrefEntryTypeFile XrefEntryType = 1

	// XrefEntryTypeObjectStream represents an object located in a compressed object stream.
	XrefEntryTypeObjectStream XrefEntryType = 2
)
