package content

import (
	"strconv"
	"strings"
	"time"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// DocumentInformation holds metadata for a PDF document.
type DocumentInformation struct {
	// DocumentInformationDictionary is the underlying document information PDF dictionary from the document.
	DocumentInformationDictionary *tokens.DictionaryToken

	// Title is the title of this document if applicable.
	Title string

	// Author is the name of the person who created this document if applicable.
	Author string

	// Subject is the subject of this document if applicable.
	Subject string

	// Keywords are any keywords associated with this document if applicable.
	Keywords string

	// Creator is the name of the application which created the original document before it was converted to PDF if applicable.
	Creator string

	// Producer is the name of the application used to convert the original document to PDF if applicable.
	Producer string

	// CreationDate is the date and time the document was created.
	CreationDate string

	// ModifiedDate is the date and time the document was most recently modified.
	ModifiedDate string

	representation string
}

// DefaultDocumentInformation is the default instance with all fields empty.
var DefaultDocumentInformation = NewDocumentInformation(nil, nil, nil, nil, nil, nil, nil, nil, nil)

// NewDocumentInformation creates a new DocumentInformation instance.
func NewDocumentInformation(
	documentInformationDictionary *tokens.DictionaryToken,
	title *string,
	author *string,
	subject *string,
	keywords *string,
	creator *string,
	producer *string,
	creationDate *string,
	modifiedDate *string,
) *DocumentInformation {
	if documentInformationDictionary == nil {
		documentInformationDictionary, _ = tokens.WithMap(make(map[string]tokens.Token))
	}

	rep := new(strings.Builder)
	appendPart(rep, "Title", derefStr(title))
	appendPart(rep, "Author", derefStr(author))
	appendPart(rep, "Subject", derefStr(subject))
	appendPart(rep, "Keywords", derefStr(keywords))
	appendPart(rep, "Creator", derefStr(creator))
	appendPart(rep, "Producer", derefStr(producer))
	appendPart(rep, "CreationDate", derefStr(creationDate))
	appendPart(rep, "ModifiedDate", derefStr(modifiedDate))

	return &DocumentInformation{
		DocumentInformationDictionary: documentInformationDictionary,
		Title:                         derefStr(title),
		Author:                        derefStr(author),
		Subject:                       derefStr(subject),
		Keywords:                      derefStr(keywords),
		Creator:                       derefStr(creator),
		Producer:                      derefStr(producer),
		CreationDate:                  derefStr(creationDate),
		ModifiedDate:                  derefStr(modifiedDate),
		representation:                rep.String(),
	}
}

// GetCreatedDateTimeOffset parses the CreationDate field and returns it as a time.Time,
// or zero value with false if parsing fails.
func (d *DocumentInformation) GetCreatedDateTimeOffset() (time.Time, bool) {
	if d.CreationDate == "" {
		return time.Time{}, false
	}

	return parsePdfDate(d.CreationDate)
}

// GetModifiedDateTimeOffset parses the ModifiedDate field and returns it as a time.Time,
// or zero value with false if parsing fails.
func (d *DocumentInformation) GetModifiedDateTimeOffset() (time.Time, bool) {
	if d.ModifiedDate == "" {
		return time.Time{}, false
	}

	return parsePdfDate(d.ModifiedDate)
}

// String returns a string representation of this document information.
// Empty entries are not shown.
func (d *DocumentInformation) String() string {
	return d.representation
}

func appendPart(builder *strings.Builder, name, value string) {
	if value == "" {
		return
	}

	builder.WriteString(name)
	builder.WriteString(": ")
	builder.WriteString(value)
	builder.WriteString("; ")
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// parsePdfDate parses a PDF-formatted date string into a time.Time.
// PDF date values follow the format D:YYYYMMDDHHmmSSOHH'mm as defined in ISO/IEC 8824.
// This is a local copy to avoid an import cycle with the util package.
func parsePdfDate(s string) (time.Time, bool) {
	defer func() { recover() }()

	if s == "" || len(s) < 4 {
		return time.Time{}, false
	}

	loc := 0
	if s[0] == 'D' && s[1] == ':' {
		loc = 2
	}

	hasRemaining := func(pos, length int) bool {
		return pos+length <= len(s)
	}

	isAtEnd := func(pos int) bool {
		return pos == len(s)
	}

	inRange := func(val, min, max int) bool {
		return val >= min && val <= max
	}

	if !hasRemaining(loc, 4) {
		return time.Time{}, false
	}

	year, err := strconv.Atoi(s[loc : loc+4])
	if err != nil {
		return time.Time{}, false
	}

	loc += 4

	if !hasRemaining(loc, 2) {
		if !isAtEnd(loc) {
			return time.Time{}, false
		}

		return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC), true
	}

	month, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || !inRange(month, 1, 12) {
		return time.Time{}, false
	}

	loc += 2

	if !hasRemaining(loc, 2) {
		if !isAtEnd(loc) {
			return time.Time{}, false
		}

		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), true
	}

	day, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || !inRange(day, 1, 31) {
		return time.Time{}, false
	}

	loc += 2

	if !hasRemaining(loc, 2) {
		if !isAtEnd(loc) {
			return time.Time{}, false
		}

		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), true
	}

	hour, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || !inRange(hour, 0, 23) {
		return time.Time{}, false
	}

	loc += 2

	if !hasRemaining(loc, 2) {
		if !isAtEnd(loc) {
			return time.Time{}, false
		}

		return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC), true
	}

	minute, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || !inRange(minute, 0, 59) {
		return time.Time{}, false
	}

	loc += 2

	if !hasRemaining(loc, 2) {
		if !isAtEnd(loc) {
			return time.Time{}, false
		}

		return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC), true
	}

	second, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || !inRange(second, 0, 59) {
		return time.Time{}, false
	}

	loc += 2

	if !hasRemaining(loc, 1) {
		if !isAtEnd(loc) {
			return time.Time{}, false
		}

		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), true
	}

	o := s[loc]
	loc++

	var sign int
	switch o {
	case '-':
		sign = -1
	case '+':
		sign = 1
	case 'Z':
		sign = 0
	default:
		return time.Time{}, false
	}

	if isAtEnd(loc) {
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), true
	}

	if !hasRemaining(loc, 3) {
		return time.Time{}, false
	}

	hoursOffset, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || s[loc+2] != '\'' || !inRange(hoursOffset, 0, 23) {
		return time.Time{}, false
	}

	loc += 3

	if isAtEnd(loc) {
		offset := time.Duration(hoursOffset*sign) * time.Hour
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.FixedZone("", int(offset.Seconds()))), true
	}

	if !hasRemaining(loc, 3) {
		return time.Time{}, false
	}

	minutesOffset, err := strconv.Atoi(s[loc : loc+2])
	if err != nil || s[loc+2] != '\'' || !inRange(minutesOffset, 0, 59) {
		return time.Time{}, false
	}

	loc += 3

	if isAtEnd(loc) {
		offset := time.Duration(hoursOffset*sign)*time.Hour + time.Duration(minutesOffset*sign)*time.Minute
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.FixedZone("", int(offset.Seconds()))), true
	}

	return time.Time{}, false
}
