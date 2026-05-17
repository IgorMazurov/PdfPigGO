// Package writer provides XMP metadata generation for PDF documents.
package writer

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/tokens"
)

const (
	xmpTk = "Adobe XMP Core 5.6-c014 79.156797, 2014/08/20-09:53:02        "

	rdfNamespace      = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	xmpMetaPrefix     = "x"
	xmpMetaNamespace  = "adobe:ns:meta/"
	dcPrefix          = "dc"
	dcNamespace       = "http://purl.org/dc/elements/1.1/"
	xmpBasicPrefix    = "xmp"
	xmpBasicNamespace = "http://ns.adobe.com/xap/1.0/"
	adobePdfPrefix    = "pdf"
	adobePdfNamespace = "http://ns.adobe.com/pdf/1.3/"
	pdfaidPrefix      = "pdfaid"
	pdfaidNamespace   = "http://www.aiim.org/pdfa/ns/id/"
)

// schemaMapper maps a document information field to an XMP element name.
type schemaMapper struct {
	name      string
	valueFunc func(info *content.DocumentInformation, version float64) string
}

// GenerateXmpStream creates the XMP metadata stream token for PDF/A compliance.
func GenerateXmpStream(
	info *content.DocumentInformation,
	version float64,
	standard PdfAStandard,
	additionalXmpMetadata *string,
) *tokens.StreamToken {
	var sb strings.Builder

	sb.WriteString(`<?xpacket begin="`)
	sb.WriteRune('\ufeff')
	sb.WriteString(`" id="W5M0MpCehiHzreSzNTczkc9d"?>` + "\n")

	sb.WriteString(`<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="` + xmpTk + `">`)
	sb.WriteString(`<rdf:RDF xmlns:rdf="` + rdfNamespace + `">`)

	// Collect all namespace declarations for the main rdf:Description element.
	var nsAttrs strings.Builder
	nsAttrs.WriteString(` xmlns:` + dcPrefix + `="` + dcNamespace + `"`)
	nsAttrs.WriteString(` xmlns:` + xmpBasicPrefix + `="` + xmpBasicNamespace + `"`)
	nsAttrs.WriteString(` xmlns:` + adobePdfPrefix + `="` + adobePdfNamespace + `"`)

	if part, _ := pdfAPartConformance(standard); part > 0 {
		nsAttrs.WriteString(` xmlns:` + pdfaidPrefix + `="` + pdfaidNamespace + `"`)
	}

	sb.WriteString(`<rdf:Description rdf:about=""` + nsAttrs.String() + `>`)

	// Dublin Core Schema elements (inside single rdf:Description)
	writeSchemaElements(&sb, dcPrefix, info, version, []schemaMapper{
		{"format", func(_ *content.DocumentInformation, _ float64) string { return "application/pdf" }},
		{"creator", func(i *content.DocumentInformation, _ float64) string { return i.Author }},
		{"description", func(i *content.DocumentInformation, _ float64) string { return i.Subject }},
		{"title", func(i *content.DocumentInformation, _ float64) string { return i.Title }},
	})

	// XMP Basic Schema elements
	writeSchemaElements(&sb, xmpBasicPrefix, info, version, []schemaMapper{
		{"CreatorTool", func(i *content.DocumentInformation, _ float64) string { return i.Creator }},
	})

	// Adobe PDF Schema elements
	writeSchemaElements(&sb, adobePdfPrefix, info, version, []schemaMapper{
		{"PDFVersion", func(_ *content.DocumentInformation, v float64) string { return fmt.Sprintf("%.1f", v) }},
		{"Producer", func(i *content.DocumentInformation, _ float64) string { return i.Producer }},
	})

	sb.WriteString(`</rdf:Description>`)

	// Separate rdf:Description for PDF/A identification
	if part, conformance := pdfAPartConformance(standard); part > 0 {
		writePdfAIdElement(&sb, part, conformance)
	}

	sb.WriteString(`</rdf:RDF>`)
	sb.WriteString(`</x:xmpmeta>`)

	xml := sb.String()
	xml = strings.ReplaceAll(xml, "\r\n", "\n")

	if additionalXmpMetadata != nil && *additionalXmpMetadata != "" {
		xml = mergeXmpDocuments(xml, *additionalXmpMetadata)
	}

	sb.Reset()
	sb.WriteString(`<?xpacket begin="`)
	sb.WriteRune('\ufeff')
	sb.WriteString(`" id="W5M0MpCehiHzreSzNTczkc9d"?>` + "\n")
	sb.WriteString(xml)
	sb.WriteString("\n<?xpacket end=\"r\"?>")

	xmlBytes := []byte(sb.String())

	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Type:    tokens.Metadata,
		tokens.Subtype: tokens.Xml,
		tokens.Length:  tokens.NewNumericTokenFromInt(len(xmlBytes)),
	})

	stream, _ := tokens.NewStreamToken(dict, xmlBytes)
	return stream
}

// writeSchemaElements writes child elements for a schema inside an existing
// rdf:Description. The namespace must already be declared on the parent element.
func writeSchemaElements(sb *strings.Builder, prefix string, info *content.DocumentInformation, version float64, mappers []schemaMapper) {
	for _, mapper := range mappers {
		val := mapper.valueFunc(info, version)
		if val == "" {
			continue
		}
		sb.WriteString(`<` + prefix + `:` + mapper.name + `>`)
		sb.WriteString(escapeXml(val))
		sb.WriteString(`</` + prefix + `:` + mapper.name + `>`)
	}
}

// writePdfAIdElement writes a separate rdf:Description for PDF/A identification.
func writePdfAIdElement(sb *strings.Builder, part int, conformance string) {
	sb.WriteString(`<rdf:Description rdf:about="" xmlns:` + pdfaidPrefix + `="` + pdfaidNamespace + `">`)
	sb.WriteString(fmt.Sprintf(`<`+pdfaidPrefix+`:part>%d</`+pdfaidPrefix+`:part>`, part))
	sb.WriteString(fmt.Sprintf(`<`+pdfaidPrefix+`:conformance>%s</`+pdfaidPrefix+`:conformance>`, conformance))
	sb.WriteString(`</rdf:Description>`)
}

// pdfAPartConformance returns the part number and conformance letter for a PDF/A standard.
func pdfAPartConformance(standard PdfAStandard) (int, string) {
	switch standard {
	case PdfA1B:
		return 1, "B"
	case PdfA1A:
		return 1, "A"
	case PdfA2B:
		return 2, "B"
	case PdfA2A:
		return 2, "A"
	case PdfA3B:
		return 3, "B"
	case PdfA3A:
		return 3, "A"
	default:
		return 0, ""
	}
}

// escapeXml escapes special XML characters in the given string.
func escapeXml(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// mergeXmpDocuments merges two XMP XML documents, avoiding duplicate elements.
// Elements already present in the base document are not duplicated from additional.
func mergeXmpDocuments(base, additional string) string {
	rdfTag := "rdf:RDF"
	baseOpenIdx := strings.Index(base, "<"+rdfTag)
	addOpenIdx := strings.Index(additional, "<"+rdfTag)

	if baseOpenIdx == -1 || addOpenIdx == -1 {
		return base
	}

	// Collect element names already present in the base document.
	existingNames := collectDescendantNames(base, baseOpenIdx)

	baseCloseIdx := strings.LastIndex(base, "</"+rdfTag+">")
	if baseCloseIdx == -1 {
		baseCloseIdx = len(base) - 1
	}

	// Extract rdf:Description blocks from the additional document.
	addContentStart := addOpenIdx + len("<" + rdfTag + ">")
	addCloseIdx := strings.LastIndex(additional[addContentStart:], "</"+rdfTag+">")
	if addCloseIdx == -1 {
		return base
	}

	var newDescriptions strings.Builder
	segment := additional[addContentStart : addContentStart+addCloseIdx]

	for {
		descOpen := strings.Index(segment, "<rdf:Description")
		if descOpen == -1 {
			break
		}

		closeDesc := strings.Index(segment[descOpen:], "</rdf:Description>")
		if closeDesc == -1 {
			break
		}

		descEnd := descOpen + closeDesc + len("</rdf:Description>")
		descContent := segment[descOpen:descEnd]

		filtered := filterDuplicateElements(descContent, existingNames)
		newDescriptions.WriteString(filtered)

		segment = segment[descEnd:]
	}

	var sb strings.Builder
	sb.WriteString(base[:baseCloseIdx])
	sb.WriteString(newDescriptions.String())
	sb.WriteString("</" + rdfTag + ">")
	// Include everything after </rdf:RDF> from base (e.g., </x:xmpmeta>)
	closeTagLen := len("</" + rdfTag + ">")
	if baseCloseIdx+closeTagLen <= len(base) {
		sb.WriteString(base[baseCloseIdx+closeTagLen:])
	}
	return sb.String()
}

// collectDescendantNames collects all element names from rdf:Description descendants in the given XML.
func collectDescendantNames(xml string, rdfIdx int) map[string]bool {
	names := make(map[string]bool)

	contentStart := rdfIdx + len("<rdf:RDF>")
	closeIdx := strings.LastIndex(xml[contentStart:], "</rdf:RDF>")
	if closeIdx == -1 {
		return names
	}

	contentEnd := contentStart + closeIdx
	segment := xml[contentStart:contentEnd]

	for {
		descIdx := strings.Index(segment, "<rdf:Description")
		if descIdx == -1 {
			break
		}

		closeDesc := strings.Index(segment[descIdx:], "</rdf:Description>")
		if closeDesc == -1 {
			break
		}

		descContent := segment[descIdx : descIdx+closeDesc]
		idx := 0
		for idx < len(descContent) {
			tagIdx := strings.Index(descContent[idx:], "<")
			if tagIdx == -1 {
				break
			}

			absIdx := idx + tagIdx
			closeTag := strings.Index(descContent[absIdx:], ">")
			if closeTag == -1 {
				break
			}

			tagEnd := absIdx + closeTag
			tagContent := descContent[absIdx : tagEnd+1]

			if !strings.HasPrefix(tagContent, "<rdf:") && !strings.Contains(tagContent, "xmlns:") {
				nameStart := strings.Index(tagContent, ":")
				if nameStart > 0 {
					elemName := tagContent[nameStart+1:closeTag]
					names[elemName] = true
				}
			}

			idx = tagEnd + 1
		}

		segment = segment[descIdx+closeDesc:]
	}

	return names
}

// filterDuplicateElements removes child elements from an rdf:Description that already exist in the given set.
func filterDuplicateElements(desc string, existingNames map[string]bool) string {
	var sb strings.Builder

	openEnd := strings.Index(desc, ">")
	if openEnd == -1 {
		return desc
	}

	sb.WriteString(desc[:openEnd+1])

	contentStart := openEnd + 1
	closeIdx := strings.LastIndex(desc, "</rdf:Description>")
	if closeIdx == -1 {
		closeIdx = len(desc)
	}

	if contentStart >= closeIdx {
		sb.WriteString(desc[contentStart:])
		return sb.String()
	}

	segment := desc[contentStart:closeIdx]
	idx := 0
	for idx < len(segment) {
		tagIdx := strings.Index(segment[idx:], "<")
		if tagIdx == -1 || idx+tagIdx >= len(segment) {
			sb.WriteString(segment[idx:])
			break
		}

		absIdx := idx + tagIdx

		// Preserve text content between previous position and this tag.
		if absIdx > idx {
			sb.WriteString(segment[idx:absIdx])
		}

		closeTag := strings.Index(segment[absIdx:], ">")
		if closeTag == -1 {
			sb.WriteString(segment[absIdx:])
			break
		}

		tagEnd := absIdx + closeTag
		tagContent := segment[absIdx : tagEnd+1]

		if !strings.HasPrefix(tagContent, "<rdf:") && !strings.HasPrefix(tagContent, "</") {
			nameStart := strings.Index(tagContent, ":")
			if nameStart > 0 {
				// Extract local element name (after first colon) for existing check.
				localEnd := len(tagContent) - 1 // before ">"
				spacePos := strings.Index(tagContent, " ")
				if spacePos > 0 && spacePos < localEnd {
					localEnd = spacePos
				}
				elemName := tagContent[nameStart+1 : localEnd]

				if existingNames[elemName] {
					// Use full element identifier (prefix:name) for closing tag search.
					fullIdentifier := tagContent[1:localEnd] // e.g., "pdfaid:part"
					closeElem := `</` + fullIdentifier + `>`
					fullCloseIdx := strings.Index(segment[tagEnd:], closeElem)
					if fullCloseIdx != -1 {
						idx = tagEnd + fullCloseIdx + len(closeElem)
						continue
					}
				}
			}
		}

		sb.WriteString(tagContent)
		idx = tagEnd + 1
	}

	sb.WriteString("</rdf:Description>")
	return sb.String()
}
