package writer_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/writer"
)

func init() {
	testutil.IntegrationDocumentsRoot = "../testdata/integration/Documents"
}

const (
	pdfaidNamespace   = "http://www.aiim.org/pdfa/ns/id/"
	ftxFormsNamespace = "http://ns.ftx.com/forms/1.0/"
	ftxControlNs      = "http://ns.ftx.com/forms/1.0/controldata/"
)

type xmpNode struct {
	XMLName xml.Name
	Content string   `xml:",chardata"`
	Child   []xmpNode `xml:",any"`
}

func parseXmp(xmlStr string) ([]xmpNode, error) {
	var nodes []xmpNode
	dec := xml.NewDecoder(strings.NewReader(xmlStr))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if se, ok := tok.(xml.StartElement); ok {
			var n xmpNode
			n.XMLName = se.Name
			if err := dec.DecodeElement(&n, &se); err != nil {
				continue
			}
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

func findDescendant(nodes []xmpNode, namespace, localName string) []string {
	var results []string
	for _, n := range nodes {
		if n.XMLName.Space == namespace && n.XMLName.Local == localName {
			results = append(results, strings.TrimSpace(n.Content))
			continue
		}
		if found := findDescendant(n.Child, namespace, localName); len(found) > 0 {
			results = append(results, found...)
		}
	}
	return results
}

func buildPdfA2aDocument(t *testing.T) []byte {
	t.Helper()
	docPath := testutil.GetDocumentPath("Single Page Simple - from inkscape.pdf", true)

	doc, err := parser.OpenFile(docPath, nil)
	if err != nil {
		t.Fatalf("OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	builder := writer.NewPdfDocumentBuilder()
	builder.SetArchiveStandard(writer.PdfA2A)

	_, err = builder.AddPageWithOptions(doc, 1, nil)
	if err != nil {
		t.Fatalf("AddPageWithOptions: %v", err)
	}

	built, err := builder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	err = builder.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}

	return built
}

func TestXmpInfoIsWrittenToPdfADocument(t *testing.T) {
	pdfBytes := buildPdfA2aDocument(t)

	doc, err := parser.OpenMemory(pdfBytes, nil)
	if err != nil {
		t.Fatalf("OpenMemory(pdfA2a): %v", err)
	}
	defer doc.Close()

	xmpMeta, ok, err := doc.TryGetXmpMetadata()
	if err != nil {
		t.Fatalf("TryGetXmpMetadata: %v", err)
	}
	if !ok {
		t.Fatal("expected XMP metadata to be present")
	}

	xmlStr := xmpMeta.GetXmlString()

	nodes, err := parseXmp(xmlStr)
	if err != nil {
		t.Fatalf("parseXmp: %v", err)
	}

	parts := findDescendant(nodes, pdfaidNamespace, "part")
	if len(parts) == 0 || parts[0] != "2" {
		t.Errorf("expected pdfaid:part = \"2\", got %v", parts)
	}

	conformances := findDescendant(nodes, pdfaidNamespace, "conformance")
	if len(conformances) == 0 || conformances[0] != "A" {
		t.Errorf("expected pdfaid:conformance = \"A\", got %v", conformances)
	}
}

func TestCustomXmpInfoIsMergedIntoPdfADocumentXmp(t *testing.T) {
	simpleDoc, err := parser.OpenFile(
		testutil.GetDocumentPath("Single Page Simple - from inkscape.pdf", true), nil,
	)
	if err != nil {
		t.Fatalf("OpenFile(simple): %v", err)
	}
	defer simpleDoc.Close()

	customXmp := `<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Adobe XMP Core 5.6-c014 79.156797, 2014/08/20-09:53:02        ">
  <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
    <rdf:Description
         xmlns:ftx="http://ns.ftx.com/forms/1.0/"
         xmlns:control="http://ns.ftx.com/forms/1.0/controldata/"
         xmlns:pdfaid="http://www.aiim.org/pdfa/ns/id/"
         xmlns:pdf="http://ns.adobe.com/pdf/1.3/">
      <ftx:ControlData rdf:parseType="Resource">
        <control:Anzahl_Zeichen_Titel>0</control:Anzahl_Zeichen_Titel>
        <control:Anzahl_Zeichen_Vorname>0</control:Anzahl_Zeichen_Vorname>
        <control:Anzahl_Zeichen_Namenszusatz>0</control:Anzahl_Zeichen_Namenszusatz>
        <control:Anzahl_Zeichen_Hausnummer>0</control:Anzahl_Zeichen_Hausnummer>
        <control:Anzahl_Zeichen_Postleitzahl>0</control:Anzahl_Zeichen_Postleitzahl>
        <control:Anzahl_Zeichen_Wohnsitzlaendercode>0</control:Anzahl_Zeichen_Wohnsitzlaendercode>
        <control:Auftragsnummer_Einsender>0</control:Auftragsnummer_Einsender>
        <control:Formularnummer>10</control:Formularnummer>
        <control:Formularversion>10.2020</control:Formularversion>
        <control:Technische_Version>6</control:Technische_Version>
      </ftx:ControlData>
      <pdfaid:part>1</pdfaid:part>
      <pdfaid:conformance>B</pdfaid:conformance>
    </rdf:Description>
  </rdf:RDF>
</x:xmpmeta>`

	builder := writer.NewPdfDocumentBuilder()
	builder.SetArchiveStandard(writer.PdfA2A)
	builder.SetIncludeDocumentInformation(true)
	builder.SetXmpMetadata(&customXmp)

	_, err = builder.AddPageWithOptions(simpleDoc, 1, nil)
	if err != nil {
		t.Fatalf("AddPageWithOptions: %v", err)
	}

	built, err := builder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	err = builder.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}

	xmpDoc, err := parser.OpenMemory(built, nil)
	if err != nil {
		t.Fatalf("OpenMemory(xmp): %v", err)
	}
	defer xmpDoc.Close()

	xmpMeta, ok, err := xmpDoc.TryGetXmpMetadata()
	if err != nil {
		t.Fatalf("TryGetXmpMetadata: %v", err)
	}
	if !ok {
		t.Fatal("expected XMP metadata to be present")
	}

	xmlStr := xmpMeta.GetXmlString()

	nodes, err := parseXmp(xmlStr)
	if err != nil {
		t.Fatalf("parseXmp: %v", err)
	}

	parts := findDescendant(nodes, pdfaidNamespace, "part")
	if len(parts) == 0 || parts[0] != "2" {
		t.Errorf("expected pdfaid:part = \"2\", got %v", parts)
	}
	if len(parts) > 1 {
		t.Errorf("expected exactly one pdfaid:part element, found %d", len(parts))
	}

	conformances := findDescendant(nodes, pdfaidNamespace, "conformance")
	if len(conformances) == 0 || conformances[0] != "A" {
		t.Errorf("expected pdfaid:conformance = \"A\", got %v", conformances)
	}
	if len(conformances) > 1 {
		t.Errorf("expected exactly one pdfaid:conformance element, found %d", len(conformances))
	}

	controlData := findDescendant(nodes, ftxFormsNamespace, "ControlData")
	if len(controlData) == 0 {
		t.Error("expected ftx:ControlData element to be present in merged XMP")
	}

	titelCount := findDescendant(nodes, ftxControlNs, "Anzahl_Zeichen_Titel")
	if len(titelCount) == 0 || titelCount[0] != "0" {
		t.Errorf("expected control:Anzahl_Zeichen_Titel = \"0\", got %v", titelCount)
	}
}
