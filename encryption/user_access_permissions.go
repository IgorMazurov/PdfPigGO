package encryption

// UserAccessPermissions defines the access permissions granted to a user when
// opening an encrypted PDF document. Each flag corresponds to a specific bit
// in the 32-bit permissions field of the PDF encryption dictionary.
type UserAccessPermissions int64

const (
	// Print allows printing the document. In revision 2 this is unrestricted.
	// In revision 3 or greater it may not allow highest quality printing;
	// see PrintHighQuality.
	Print UserAccessPermissions = 1 << 2

	// Modify allows modifying the contents of the document by operations other
	// than those controlled by AddOrModifyTextAnnotationsAndFillFormFields,
	// FillExistingFormFields, and AssembleDocument.
	Modify UserAccessPermissions = 1 << 3

	// CopyTextAndGraphics allows copying or otherwise extracting text and
	// graphics from the document, including for accessibility purposes. In
	// revision 2 this is unrestricted. In revision 3 or greater it does not
	// include operations controlled by ExtractTextAndGraphics.
	CopyTextAndGraphics UserAccessPermissions = 1 << 4

	// AddOrModifyTextAnnotationsAndFillFormFields allows adding or modifying
	// text annotations, filling in interactive form fields, and, if Modify is
	// also set, creating or modifying interactive form fields (including
	// signature fields).
	AddOrModifyTextAnnotationsAndFillFormFields UserAccessPermissions = 1 << 5

	// FillExistingFormFields allows filling in existing interactive form fields
	// (including signature fields) even if AddOrModifyTextAnnotationsAndFillFormFields
	// is clear. Available in revision 3 or greater.
	FillExistingFormFields UserAccessPermissions = 1 << 8

	// ExtractTextAndGraphics allows extracting text and graphics in support of
	// accessibility to users with disabilities or for other purposes. Available
	// in revision 3 or greater.
	ExtractTextAndGraphics UserAccessPermissions = 1 << 9

	// AssembleDocument allows assembling the document (insert, rotate, or delete
	// pages and create bookmarks or thumbnail images) even if Modify is clear.
	// Available in revision 3 or greater.
	AssembleDocument UserAccessPermissions = 1 << 10

	// PrintHighQuality allows printing the document to a representation from
	// which a faithful digital copy of the PDF content could be generated. When
	// this flag is clear (and Print is set), printing is limited to a low-level
	// representation, possibly of degraded quality. Available in revision 3 or
	// greater.
	PrintHighQuality UserAccessPermissions = 1 << 12
)
