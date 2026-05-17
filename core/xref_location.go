package core

// XrefLocation holds information about where an object is located in the file
// according to the cross-reference section or brute force parsing.
type XrefLocation struct {
	// Type indicates which kind of location this represents.
	Type XrefEntryType

	// Value1 is the byte offset when Type is XrefEntryTypeFile,
	// or the object stream number when Type is XrefEntryTypeObjectStream.
	Value1 int64

	// Value2 is the index of the object within the stream; only used when
	// Type is XrefEntryTypeObjectStream.
	Value2 int
}

// File creates a location mapped to a byte offset in the file.
func File(offset int64) XrefLocation {
	return XrefLocation{
		Type:   XrefEntryTypeFile,
		Value1: offset,
		Value2: 0,
	}
}

// Stream creates a location mapped to an index inside an object stream.
func Stream(objStream int64, index int) XrefLocation {
	return XrefLocation{
		Type:   XrefEntryTypeObjectStream,
		Value1: objStream,
		Value2: index,
	}
}
