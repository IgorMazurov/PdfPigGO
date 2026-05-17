package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/crossreference"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CreateDocumentInformation parses the document information dictionary from a PDF file trailer.
func CreateDocumentInformation(
	scanner tokenization.PdfTokenScanner,
	trailer *crossreference.TrailerDictionary,
	isLenientParsing bool,
) (*content.DocumentInformation, error) {
	token := trailer.Info()

	if ref, ok := token.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(ref.Data())
		if obj == nil {
			return content.DefaultDocumentInformation, nil
		}
		token = obj.Data()
	}

	if token == nil {
		return content.DefaultDocumentInformation, nil
	}

	if infoParsed, ok := token.(*tokens.DictionaryToken); ok {
		title := getEntryOrDefault(infoParsed, tokens.Title, scanner)
		author := getEntryOrDefault(infoParsed, tokens.Author, scanner)
		subject := getEntryOrDefault(infoParsed, tokens.Subject, scanner)
		keywords := getEntryOrDefault(infoParsed, tokens.Keywords, scanner)
		creator := getEntryOrDefault(infoParsed, tokens.Creator, scanner)
		producer := getEntryOrDefault(infoParsed, tokens.Producer, scanner)
		creationDate := getEntryOrDefault(infoParsed, tokens.CreationDate, scanner)
		modifiedDate := getEntryOrDefault(infoParsed, tokens.ModDate, scanner)

		return content.NewDocumentInformation(infoParsed, title, author, subject,
			keywords, creator, producer, creationDate, modifiedDate), nil
	}

	if streamToken, ok := token.(*tokens.StreamToken); ok {
		streamDictionary := streamToken.StreamDictionary
		typeNameTok, foundType := streamDictionary.TryGet(tokens.Type)
		if !foundType {
			return nil, core.NewPdfDocumentFormatException("Unknown document metadata type was found")
		}
		nameTok, ok := typeNameTok.(*tokens.NameToken)
		if !ok || !nameTok.Equals(tokens.Metadata) {
			return nil, core.NewPdfDocumentFormatException("Unknown document metadata type was found")
		}

		subtypeTok, foundSubtype := streamDictionary.TryGet(tokens.Subtype)
		if !foundSubtype {
			return nil, core.NewPdfDocumentFormatException("Unknown document metadata subtype was found")
		}
		nameTok2, ok := subtypeTok.(*tokens.NameToken)
		if !ok || !nameTok2.Equals(tokens.Xml) {
			return nil, core.NewPdfDocumentFormatException("Unknown document metadata subtype was found")
		}

		return content.DefaultDocumentInformation, nil
	}

	if isLenientParsing {
		return content.DefaultDocumentInformation, nil
	}

	tokenTypeName := "nil"
	if token != nil {
		tokenTypeName = fmt.Sprintf("%T", token)
	}
	return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("Unknown document information token was found %s", tokenTypeName))
}

// getEntryOrDefault extracts a string value from the info dictionary for the given key.
func getEntryOrDefault(infoDictionary *tokens.DictionaryToken, key *tokens.NameToken, scanner tokenization.PdfTokenScanner) *string {
	if infoDictionary == nil {
		return nil
	}

	value, found := infoDictionary.TryGet(key)
	if !found {
		return nil
	}

	if idr, ok := value.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(idr.Data())
		if obj == nil {
			return nil
		}

		data := obj.Data()
		if strI, ok := data.(*tokens.StringToken); ok {
			s := strI.Data()
			return &s
		}
		if hexI, ok := data.(*tokens.HexToken); ok {
			s := hexI.Data()
			return &s
		}

		return nil
	}

	if str, ok := value.(*tokens.StringToken); ok {
		s := str.Data()
		return &s
	}

	if hex, ok := value.(*tokens.HexToken); ok {
		s := hex.Data()
		return &s
	}

	return nil
}
