package page

// PageXmlTextDataSimpleType represents XML Schema data types for text content in PAGE XML documents.
type PageXmlTextDataSimpleType byte

const (
	// TextDataXsddouble represents xsd:double - Examples: "123.456", "+1234.456", "-1234.456", "-.456", "-456"
	TextDataXsddouble PageXmlTextDataSimpleType = iota

	// TextDataXsdFloat represents xsd:float - Examples: "123.456", "+1234.456", "-1.2344e56", "-.45E-6", "INF", "-INF", "NaN"
	TextDataXsdFloat

	// TextDataXsdInteger represents xsd:integer - Examples: "123456", "+00000012", "-1", "-456"
	TextDataXsdInteger

	// TextDataXsdBoolean represents xsd:boolean - Examples: "true", "false", "1", "0"
	TextDataXsdBoolean

	// TextDataXsdDate represents xsd:date - Examples: "2001-10-26", "2001-10-26+02:00", "2001-10-26Z"
	TextDataXsdDate

	// TextDataXsdTime represents xsd:time - Examples: "21:32:52", "21:32:52+02:00", "19:32:52Z"
	TextDataXsdTime

	// TextDataXsdDateTime represents xsd:dateTime - Examples: "2001-10-26T21:32:52", "2001-10-26T21:32:52+02:00"
	TextDataXsdDateTime

	// TextDataXsdString represents xsd:string - Generic text string
	TextDataXsdString

	// TextDataOther represents an XSD type not listed above or a custom type (use dataTypeDetails attribute)
	TextDataOther
)
