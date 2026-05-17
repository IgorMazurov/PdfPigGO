package truetypeparser

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/names"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
	"golang.org/x/text/encoding"
)

// NameTableParser parses the name table of a TrueType font.
type NameTableParser struct{}

// Parse reads and interprets the name (name) table from raw TrueType data.
func (p *NameTableParser) Parse(header truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, register *TableRegisterBuilder) (*tables.NameTable, error) {
	if _, err := data.Seek(int64(header.Offset), 0); err != nil {
		return nil, err
	}
	_ = data.ReadUnsignedShort() // format, unused
	count := int(data.ReadUnsignedShort())
	stringOffset := data.ReadUnsignedShort()

	nameRecords := make([]nameRecordBuilder, count)
	for i := 0; i < count; i++ {
		nameRecords[i] = readNameRecordBuilder(data)
	}

	strings := make([]names.TrueTypeNameRecord, count)
	offset := header.Offset + uint32(stringOffset)
	for i := 0; i < count; i++ {
		if rec := getTrueTypeNameRecord(nameRecords[i], data, offset); rec != nil {
			strings[i] = *rec
		}
	}

	return tables.NewNameTable(
		header,
		getName(4, strings),
		getName(1, strings),
		getName(2, strings),
		strings,
	), nil
}

func getTrueTypeNameRecord(nameRec nameRecordBuilder, data *TrueTypeDataBytes, offset uint32) *names.TrueTypeNameRecord {
	var enc encoding.Encoding = iso88591Encoding

	switch nameRec.platformId {
	case names.Windows:
		platformEncoding := names.TrueTypeWindowsEncodingIdentifier(nameRec.platformEncodingId)
		if platformEncoding == names.Symbol || platformEncoding == names.UnicodeBmp {
			enc = utf16BEEncoding
		}
	case names.Unicode:
		enc = utf16BEEncoding
	case names.Iso:
		switch nameRec.platformEncodingId {
		case 0:
			enc = asciiEncoding
		case 1:
			enc = utf16BEEncoding
		}
	}

	position := offset + uint32(nameRec.offset)
	if position >= uint32(data.Length()) {
		return nil
	}

	if _, err := data.Seek(int64(position), 0); err != nil {
		return nil
	}

	str, ok := data.TryReadString(int(nameRec.length), enc)
	if !ok {
		return nil
	}

	return names.NewTrueTypeNameRecord(
		nameRec.platformId,
		nameRec.platformEncodingId,
		nameRec.languageId,
		nameRec.nameId,
		str,
	)
}

func getName(nameID uint16, nameRecords []names.TrueTypeNameRecord) string {
	windowsEnUs := uint16(409)
	var windowsVal string
	var anyVal string

	for _, record := range nameRecords {
		if record.NameId != nameID {
			continue
		}

		if record.PlatformId == names.Windows && record.LanguageId == windowsEnUs {
			return record.Value
		}

		if record.PlatformId == names.Windows {
			windowsVal = record.Value
		}

		anyVal = record.Value
	}

	if windowsVal != "" {
		return windowsVal
	}
	return anyVal
}

/* ------------------------------------------------------------------ */
/*  nameRecordBuilder (internal struct mirroring C# NameRecordBuilder) */
/* ------------------------------------------------------------------ */

type nameRecordBuilder struct {
	platformId           names.TrueTypePlatformIdentifier
	platformEncodingId   uint16
	languageId           uint16
	nameId               uint16
	length               uint16
	offset               uint16
}

func readNameRecordBuilder(data *TrueTypeDataBytes) nameRecordBuilder {
	return nameRecordBuilder{
		platformId:         names.TrueTypePlatformIdentifier(data.ReadUnsignedShort()),
		platformEncodingId: data.ReadUnsignedShort(),
		languageId:         data.ReadUnsignedShort(),
		nameId:             data.ReadUnsignedShort(),
		length:             data.ReadUnsignedShort(),
		offset:             data.ReadUnsignedShort(),
	}
}
