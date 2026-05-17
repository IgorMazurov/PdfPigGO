package acroforms

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/acroforms/fields"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Catalog represents the minimal catalog interface needed by AcroFormFactory.
type Catalog interface {
	CatalogDictionary() *tokens.DictionaryToken
	GetPageNumberByReference(ref core.IndirectReference) *int
}

var inheritableFields = map[string]bool{
	tokens.Ft.Data():   true,
	tokens.Ff.Data():   true,
	tokens.V.Data():    true,
	tokens.Dv.Data():   true,
	tokens.Aa.Data():   true,
}

// AcroFormFactory extracts the AcroForm from a PDF document if available.
type AcroFormFactory struct {
	tokenScanner  tokenization.PdfTokenScanner
	filterProvider filters.LookupFilterProvider
	objectOffsets map[core.IndirectReference]core.XrefLocation
}

// NewAcroFormFactory creates a new AcroFormFactory.
func NewAcroFormFactory(
	tokenScanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
	objectOffsets map[core.IndirectReference]core.XrefLocation,
) (*AcroFormFactory, error) {
	if tokenScanner == nil {
		return nil, fmt.Errorf("tokenScanner cannot be nil")
	}
	if filterProvider == nil {
		return nil, fmt.Errorf("filterProvider cannot be nil")
	}
	if objectOffsets == nil {
		return nil, fmt.Errorf("objectOffsets cannot be nil")
	}

	return &AcroFormFactory{
		tokenScanner:   tokenScanner,
		filterProvider: filterProvider,
		objectOffsets:  objectOffsets,
	}, nil
}

// GetAcroForm retrieves the AcroForm from the document if applicable.
func (f *AcroFormFactory) GetAcroForm(catalog Catalog) (*AcroForm, error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog cannot be nil")
	}

	catalogDict := catalog.CatalogDictionary()
	if catalogDict == nil {
		return nil, nil
	}

	acroRawToken, ok := catalogDict.TryGet(tokens.AcroForm)
	if !ok {
		return nil, nil
	}

	acroDict, resolved := parts.TryGet[*tokens.DictionaryToken](acroRawToken, f.tokenScanner)
	if !resolved || acroDict == nil {
		fieldsRefs, err := f.bruteForceScanFields()
		if err != nil || len(fieldsRefs) == 0 {
			return nil, nil
		}

		arrData := make([]tokens.Token, len(fieldsRefs))
		for i, ref := range fieldsRefs {
			arrData[i] = ref
		}
		arr := tokens.NewArrayToken(arrData)
		var syntheticDict *tokens.DictionaryToken
		syntheticDict, _ = tokens.WithMap(map[string]tokens.Token{
			tokens.Fields.Data(): arr,
		})
		acroDict = syntheticDict
	}

	dictToken := acroDict
	if dictToken == nil {
		return nil, &core.PdfDocumentFormatException{Message: "acroForm dictionary token is not a DictionaryToken"}
	}

	sigFlags, err := f.getSignatureFlags(dictToken)
	if err != nil {
		return nil, err
	}
	needAppearances, err := f.getNeedAppearances(dictToken)
	if err != nil {
		return nil, err
	}

	fieldsArray, err := f.getFieldsArray(dictToken)
	if err != nil || fieldsArray == nil {
		return nil, nil
	}

	fieldsMap := make(map[core.IndirectReference]any, fieldsArray.Length())

	for _, fieldToken := range fieldsArray.Data() {
		refToken, ok := fieldToken.(*tokens.IndirectReferenceToken)
		if !ok || refToken == nil {
			return nil, &core.PdfDocumentFormatException{
				Message: fmt.Sprintf("the fields array should only contain indirect references, instead got: %v", fieldToken),
			}
		}

		fieldDict, resolved := parts.TryGet[*tokens.DictionaryToken](fieldToken, f.tokenScanner)
		if !resolved || fieldDict == nil {
			continue
		}

		dict := fieldDict

		field := f.getAcroField(dict, catalog, []*tokens.DictionaryToken{})
		if field != nil {
			ref := refToken.Data()
			fieldsMap[ref] = field
		}
	}

	return NewAcroForm(dictToken, sigFlags, needAppearances, fieldsMap)
}

func (f *AcroFormFactory) bruteForceScanFields() ([]*tokens.IndirectReferenceToken, error) {
	var fieldsRefs []*tokens.IndirectReferenceToken

	for ref := range f.objectOffsets {
		refToken := tokens.NewIndirectReferenceToken(ref)
		dict, ok := parts.TryGet[tokens.Token](refToken, f.tokenScanner)
		if !ok || dict == nil {
			continue
		}

		dictToken, ok := dict.(*tokens.DictionaryToken)
		if !ok {
			continue
		}

		hasKids := dictToken.ContainsKey(tokens.Kids)
		hasName := dictToken.ContainsKey(tokens.T)
		if hasKids && hasName {
			fieldsRefs = append(fieldsRefs, refToken)
		}
	}

	return fieldsRefs, nil
}

func (f *AcroFormFactory) getSignatureFlags(dict *tokens.DictionaryToken) (SignatureFlags, error) {
	token, ok := dict.TryGet(tokens.SigFlags)
	if !ok {
		return 0, nil
	}

	numToken, resolved := parts.TryGet[*tokens.NumericToken](token, f.tokenScanner)
	if !resolved || numToken == nil {
		return 0, nil
	}

	return SignatureFlags(numToken.IntVal()), nil
}

func (f *AcroFormFactory) getNeedAppearances(dict *tokens.DictionaryToken) (bool, error) {
	token, ok := dict.TryGet(tokens.NeedAppearances)
	if !ok {
		return false, nil
	}

	boolToken, resolved := parts.TryGet[*tokens.BooleanToken](token, f.tokenScanner)
	if !resolved || boolToken == nil {
		return false, nil
	}

	return boolToken.Data(), nil
}

func (f *AcroFormFactory) getFieldsArray(dict *tokens.DictionaryToken) (*tokens.ArrayToken, error) {
	token, ok := dict.TryGet(tokens.Fields)
	if !ok {
		return nil, nil
	}

	arrToken, resolved := parts.TryGet[*tokens.ArrayToken](token, f.tokenScanner)
	if !resolved || arrToken == nil {
		return nil, nil
	}

	return arrToken, nil
}

func (f *AcroFormFactory) getAcroField(
	fieldDict *tokens.DictionaryToken,
	catalog Catalog,
	parentDictionaries []*tokens.DictionaryToken,
) any {
	if fieldDict == nil {
		return nil
	}

	combinedDict, inheritsValue := createInheritedDictionary(fieldDict, parentDictionaries)

	fieldTypeRaw, _ := combinedDict.TryGet(tokens.Ft)
	fieldFlagsRaw, _ := combinedDict.TryGet(tokens.Ff)

	var fieldTypeStr string
	if nameTok, ok := fieldTypeRaw.(*tokens.NameToken); ok && nameTok != nil {
		fieldTypeStr = nameTok.Data()
	}

	var fieldFlags uint32 = 0
	if numToken, ok := fieldFlagsRaw.(*tokens.NumericToken); ok && numToken != nil {
		fieldFlags = uint32(numToken.LongVal())
	}

	kids := make([]struct {
		hasParent bool
		dict      *tokens.DictionaryToken
	}, 0)

	kidsTokenRaw, hasKids := combinedDict.TryGet(tokens.Kids)
	if hasKids {
		kidsArr, ok := kidsTokenRaw.(*tokens.ArrayToken)
		if !ok || kidsArr == nil {
			resolved, _ := parts.TryGet[*tokens.ArrayToken](kidsTokenRaw, f.tokenScanner)
			if resolved != nil {
				kidsArr = resolved
			}
		}
		if kidsArr != nil {
			for _, kid := range kidsArr.Data() {
				refToken, ok := kid.(*tokens.IndirectReferenceToken)
				if !ok || refToken == nil {
					continue
				}

				kidDict, resolvedKid := parts.TryGet[*tokens.DictionaryToken](kid, f.tokenScanner)
				if !resolvedKid || kidDict == nil {
					continue
				}

				hasParent := kidDict.ContainsKey(tokens.Parent)
				kids = append(kids, struct {
					hasParent bool
					dict      *tokens.DictionaryToken
				}{hasParent: hasParent, dict: kidDict})
			}
		}
	}

	partialFieldName := ""
	tToken, _ := combinedDict.TryGet(tokens.T)
	if strTok, r := parts.TryGet[*tokens.StringToken](tToken, f.tokenScanner); r && strTok != nil {
		partialFieldName = strTok.Data()
	} else if hexTok, r := parts.TryGet[*tokens.HexToken](tToken, f.tokenScanner); r && hexTok != nil {
		partialFieldName = hexTok.Data()
	}

	alternateFieldName := ""
	tuToken, _ := combinedDict.TryGet(tokens.Tu)
	if strTok, r := parts.TryGet[*tokens.StringToken](tuToken, f.tokenScanner); r && strTok != nil {
		alternateFieldName = strTok.Data()
	}

	mappingName := ""
	tmToken, _ := combinedDict.TryGet(tokens.Tm)
	if strTok, r := parts.TryGet[*tokens.StringToken](tmToken, f.tokenScanner); r && strTok != nil {
		mappingName = strTok.Data()
	}

	var parentRef *core.IndirectReference
	parentToken, _ := combinedDict.TryGet(tokens.Parent)
	if refTok, ok := parentToken.(*tokens.IndirectReferenceToken); ok && refTok != nil {
		ref := refTok.Data()
		parentRef = &ref
	}

	information := fields.NewAcroFieldCommonInformation(parentRef, partialFieldName, alternateFieldName, mappingName)

	var pageNumber *int
	pageRefRaw, _ := combinedDict.TryGet(tokens.P)
	if refTok, ok := pageRefRaw.(*tokens.IndirectReferenceToken); ok && refTok != nil {
		if catalog != nil {
			pageNumber = catalog.GetPageNumberByReference(refTok.Data())
		}
	}

	var bounds *core.PdfRectangle
	rectToken, _ := combinedDict.TryGet(tokens.Rect)
	if rectArr, ok := rectToken.(*tokens.ArrayToken); ok && rectArr != nil && rectArr.Length() == 4 {
		bounds = arrayToRectangle(rectArr)
	}

	newParentDictionaries := make([]*tokens.DictionaryToken, len(parentDictionaries)+1)
	copy(newParentDictionaries, parentDictionaries)
	newParentDictionaries[len(parentDictionaries)] = combinedDict

	children := make([]any, 0)
	for _, kid := range kids {
		if !kid.hasParent {
			continue
		}
		childAny := f.getAcroField(kid.dict, catalog, newParentDictionaries)
		if childAny != nil {
			children = append(children, childAny)
		}
	}

	return f.createField(fieldTypeStr, combinedDict, fieldFlags, information, pageNumber, bounds, children, inheritsValue)
}

func (f *AcroFormFactory) createField(
	fieldTypeStr string,
	dict *tokens.DictionaryToken,
	fieldFlags uint32,
	information *fields.AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
	children []any,
	inheritsValue bool,
) any {
	if fieldTypeStr == "" {
		result, _ := fields.NewAcroNonTerminalField(dict, "Non-Terminal Field", fieldFlags, information, fields.AcroTypeUnknown, children)
		return result
	}

	switch fieldTypeStr {
	case tokens.Btn.Data():
		return f.createButtonField(dict, fieldTypeStr, fieldFlags, information, pageNumber, bounds, children, inheritsValue)
	case tokens.Tx.Data():
		return f.getTextField(dict, fieldTypeStr, fieldFlags, information, pageNumber, bounds)
	case tokens.Ch.Data():
		return f.getChoiceField(dict, fieldTypeStr, fieldFlags, information, pageNumber, bounds)
	case tokens.Sig.Data():
		result, _ := fields.NewAcroSignatureField(dict, fieldTypeStr, fieldFlags, information, pageNumber, bounds)
		return result
	default:
		result, _ := fields.NewAcroNonTerminalField(dict, fieldTypeStr, fieldFlags, information, fields.AcroTypeUnknown, children)
		return result
	}
}

func toAcroFieldBase(field interface{}) *fields.AcroFieldBase {
	if field == nil {
		return nil
	}

	switch v := field.(type) {
	case *fields.AcroFieldBase:
		if v == nil {
			return nil
		}
		return v
	case *fields.AcroNonTerminalField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroTextField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroCheckboxField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroRadioButtonField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroPushButtonField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroSignatureField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroRadioButtonsField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroCheckboxesField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroComboBoxField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroListBoxField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	}
	return nil
}

func (f *AcroFormFactory) createButtonField(
	dict *tokens.DictionaryToken,
	fieldTypeStr string,
	fieldFlags uint32,
	information *fields.AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
	children []any,
	inheritsValue bool,
) any {
	btnFlags := fields.AcroButtonFieldFlags(fieldFlags)

	if btnFlags&fields.Radio != 0 {
		if len(children) > 0 {
			result, _ := fields.NewAcroRadioButtonsField(dict, fieldTypeStr, btnFlags, information, children)
			return result
		}

		isChecked, valueToken := f.getCheckedState(dict, inheritsValue)
		result, _ := fields.NewAcroRadioButtonField(dict, fieldTypeStr, btnFlags, information, pageNumber, bounds, valueToken, isChecked)
		return result
	}

	if btnFlags&fields.PushButton != 0 {
		result, _ := fields.NewAcroPushButtonField(dict, fieldTypeStr, btnFlags, information, pageNumber, bounds)
		return result
	}

	if len(children) > 0 {
		result, _ := fields.NewAcroCheckboxesField(dict, fieldTypeStr, btnFlags, information, children)
		return result
	}

	isChecked, valueToken := f.getCheckedState(dict, inheritsValue)
	result, _ := fields.NewAcroCheckboxField(dict, fieldTypeStr, btnFlags, information, valueToken, isChecked, pageNumber, bounds)
	return result
}

func (f *AcroFormFactory) getTextField(
	dict *tokens.DictionaryToken,
	fieldTypeStr string,
	fieldFlags uint32,
	information *fields.AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
) any {
	textValue := f.extractTextValue(dict)
	maxLen := f.extractMaxLen(dict)

	result, _ := fields.NewAcroTextField(dict, fieldTypeStr, fields.AcroTextFieldFlags(fieldFlags), information, &textValue, maxLen, pageNumber, bounds)
	return result
}

func (f *AcroFormFactory) extractTextValue(dict *tokens.DictionaryToken) string {
	vToken, ok := dict.TryGet(tokens.V)
	if !ok {
		return ""
	}

	if strTok, resolved := parts.TryGet[*tokens.StringToken](vToken, f.tokenScanner); resolved && strTok != nil {
		return strTok.Data()
	}
	if hexTok, resolved := parts.TryGet[*tokens.HexToken](vToken, f.tokenScanner); resolved && hexTok != nil {
		return hexTok.Data()
	}

	streamTok, resolved := parts.TryGet[*tokens.StreamToken](vToken, f.tokenScanner)
	if !resolved || streamTok == nil {
		return ""
	}

	data, err := decodeStream(streamTok, f.filterProvider, f.tokenScanner)
	if err != nil {
		return ""
	}
	return core.BytesAsLatin1String(data)
}

// decodeStream decodes a StreamToken through its registered filter chain.
func decodeStream(stream *tokens.StreamToken, provider filters.LookupFilterProvider, scanner tokenization.PdfTokenScanner) ([]byte, error) {
	filterList, err := provider.GetFiltersWithScanner(stream.StreamDictionary, scanner)
	if err != nil || len(filterList) == 0 {
		return stream.Data(), nil
	}

	data := stream.Data()
	for i, flt := range filterList {
		decoded, err := flt.Decode(data, stream.StreamDictionary, provider, i)
		if err != nil {
			return data, err
		}
		data = decoded
	}
	return data, nil
}

func (f *AcroFormFactory) extractMaxLen(dict *tokens.DictionaryToken) *int {
	token, ok := dict.TryGet(tokens.MaxLen)
	if !ok {
		return nil
	}

	numTok, resolved := parts.TryGet[*tokens.NumericToken](token, f.tokenScanner)
	if !resolved || numTok == nil {
		return nil
	}

	v := numTok.IntVal()
	return &v
}

func (f *AcroFormFactory) getChoiceField(
	dict *tokens.DictionaryToken,
	fieldTypeStr string,
	fieldFlags uint32,
	information *fields.AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
) any {
	selectedOptions, selectedIndices := f.extractSelectedOptions(dict)
	options := f.extractChoiceOptions(dict, selectedOptions, selectedIndices)

	choiceFlags := fields.AcroChoiceFieldFlags(fieldFlags)

	if choiceFlags&fields.AcroChoiceCombo != 0 {
		result, _ := fields.NewAcroComboBoxField(dict, fieldTypeStr, choiceFlags, information, options, selectedOptions, selectedIndices, pageNumber, bounds)
		return result
	}

	topIndex := f.extractTopIndex(dict)
	result, _ := fields.NewAcroListBoxField(dict, fieldTypeStr, choiceFlags, information, options, selectedOptions, selectedIndices, topIndex, pageNumber, bounds)
	return result
}

func (f *AcroFormFactory) extractSelectedOptions(dict *tokens.DictionaryToken) ([]string, []int) {
	selectedOptions := []string{}

	vToken, ok := dict.TryGet(tokens.V)
	if !ok {
		return selectedOptions, nil
	}

	if strTok, resolved := parts.TryGet[*tokens.StringToken](vToken, f.tokenScanner); resolved && strTok != nil {
		selectedOptions = []string{strTok.Data()}
	} else if hexTok, resolved := parts.TryGet[*tokens.HexToken](vToken, f.tokenScanner); resolved && hexTok != nil {
		selectedOptions = []string{hexTok.Data()}
	} else if arrTok, resolved := parts.TryGet[*tokens.ArrayToken](vToken, f.tokenScanner); resolved && arrTok != nil {
		for _, elem := range arrTok.Data() {
			if sTok, r := parts.TryGet[*tokens.StringToken](elem, f.tokenScanner); r && sTok != nil {
				selectedOptions = append(selectedOptions, sTok.Data())
			} else if hTok, r := parts.TryGet[*tokens.HexToken](elem, f.tokenScanner); r && hTok != nil {
				selectedOptions = append(selectedOptions, hTok.Data())
			}
		}
	}

	selectedIndices := extractSelectedIndices(dict, f.tokenScanner)
	return selectedOptions, selectedIndices
}

func extractSelectedIndices(
	dict *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
) []int {
	iToken, ok := dict.TryGet(tokens.I)
	if !ok {
		return nil
	}

	arrTok, ok := iToken.(*tokens.ArrayToken)
	if !ok || arrTok == nil {
		resolved, _ := parts.TryGet[*tokens.ArrayToken](iToken, scanner)
		if resolved == nil {
			return nil
		}
		arrTok = resolved
	}
	if arrTok == nil {
		return nil
	}

	data := arrTok.Data()
	selectedIndices := make([]int, len(data))
	for i := 0; i < len(data); i++ {
		elem := data[i]
		numTok, resolved := parts.TryGet[*tokens.NumericToken](elem, scanner)
		if !resolved || numTok == nil {
			continue
		}
		selectedIndices[i] = numTok.IntVal()
	}

	return selectedIndices
}

func (f *AcroFormFactory) extractChoiceOptions(
	dict *tokens.DictionaryToken,
	selectedOptions []string,
	selectedIndices []int,
) []*fields.AcroChoiceOption {
	options := make([]*fields.AcroChoiceOption, 0)

	optToken, ok := dict.TryGet(tokens.Opt)
	if !ok {
		return options
	}

	arrTok, ok := optToken.(*tokens.ArrayToken)
	if !ok || arrTok == nil {
		return options
	}

	for i, optionToken := range arrTok.Data() {
		if strTok, r := parts.TryGet[*tokens.StringToken](optionToken, f.tokenScanner); r && strTok != nil {
			isSelected := isChoiceSelected(selectedOptions, selectedIndices, i, strTok.Data())
			options = append(options, fields.NewAcroChoiceOption(i, isSelected, strTok.Data(), nil))
		} else if hexTok, r := parts.TryGet[*tokens.HexToken](optionToken, f.tokenScanner); r && hexTok != nil {
			isSelected := isChoiceSelected(selectedOptions, selectedIndices, i, hexTok.Data())
			options = append(options, fields.NewAcroChoiceOption(i, isSelected, hexTok.Data(), nil))
		} else if innerArr, r := parts.TryGet[*tokens.ArrayToken](optionToken, f.tokenScanner); r && innerArr != nil {
			if innerArr.Length() != 2 {
				continue
			}

			exportValue, name := f.resolveOptionArray(innerArr)
			if exportValue == "" || name == "" {
				continue
			}

			isSelected := isChoiceSelected(selectedOptions, selectedIndices, i, name)
			options = append(options, fields.NewAcroChoiceOption(i, isSelected, name, &exportValue))
		}
	}

	return options
}

func (f *AcroFormFactory) resolveOptionArray(arr *tokens.ArrayToken) (string, string) {
	data := arr.Data()
	if len(data) < 2 {
		return "", ""
	}

	var exportValue string
	if sTok, r := parts.TryGet[*tokens.StringToken](data[0], f.tokenScanner); r && sTok != nil {
		exportValue = sTok.Data()
	} else if hTok, r := parts.TryGet[*tokens.HexToken](data[0], f.tokenScanner); r && hTok != nil {
		exportValue = hTok.Data()
	}

	var name string
	if sTok, r := parts.TryGet[*tokens.StringToken](data[1], f.tokenScanner); r && sTok != nil {
		name = sTok.Data()
	} else if hTok, r := parts.TryGet[*tokens.HexToken](data[1], f.tokenScanner); r && hTok != nil {
		name = hTok.Data()
	}

	return exportValue, name
}

func (f *AcroFormFactory) extractTopIndex(dict *tokens.DictionaryToken) *int {
	token, ok := dict.TryGet(tokens.Ti)
	if !ok {
		return nil
	}

	numTok, resolved := parts.TryGet[*tokens.NumericToken](token, f.tokenScanner)
	if !resolved || numTok == nil {
		return nil
	}

	v := numTok.IntVal()
	return &v
}

func (f *AcroFormFactory) getCheckedState(dict *tokens.DictionaryToken, inheritsValue bool) (bool, *tokens.NameToken) {
	vRaw, hasV := dict.TryGet(tokens.V)

	if !hasV {
		asRaw, _ := dict.TryGet(tokens.As)
		_, hasAp := dict.TryGet(tokens.Ap)

		nameTok, resolvedAs := parts.TryGet[*tokens.NameToken](asRaw, f.tokenScanner)
		if resolvedAs && nameTok != nil && hasAp {
			isChecked := !strings.EqualFold(nameTok.Data(), tokens.OffAcroform.Data())
			return isChecked, nameTok
		}

		return false, tokens.OffAcroform
	}

	nameVal, resolvedV := parts.TryGet[*tokens.NameToken](vRaw, f.tokenScanner)
	if !resolvedV || nameVal == nil {
		return false, tokens.OffAcroform
	}

	if inheritsValue {
		asRaw, _ := dict.TryGet(tokens.As)
		asName, resolvedAs := parts.TryGet[*tokens.NameToken](asRaw, f.tokenScanner)
		if resolvedAs && asName != nil {
			isChecked := asName.Data() == nameVal.Data()
			return isChecked, asName
		}
	}

	isChecked := !strings.EqualFold(nameVal.Data(), tokens.OffAcroform.Data())
	return isChecked, nameVal
}

func arrayToRectangle(arr *tokens.ArrayToken) *core.PdfRectangle {
	data := arr.Data()
	if len(data) < 4 {
		return nil
	}

	vals := make([]float64, 0, 4)
	for _, tok := range data[:4] {
		if numTok, ok := tok.(*tokens.NumericToken); ok && numTok != nil {
			vals = append(vals, numTok.DoubleVal())
		}
	}

	if len(vals) < 4 {
		return nil
	}

	rect := core.NewPdfRectangleFloat(vals[0], vals[1], vals[2], vals[3])
	return &rect
}

func createInheritedDictionary(
	fieldDict *tokens.DictionaryToken,
	parents []*tokens.DictionaryToken,
) (*tokens.DictionaryToken, bool) {
	if len(parents) == 0 || fieldDict == nil {
		return fieldDict, false
	}

	inheritsValue := false
	inheritedMap := make(map[string]tokens.Token)

	for _, parent := range parents {
		if parent == nil {
			continue
		}
		for k, v := range parent.Data() {
			if inheritableFields[k] {
				inheritedMap[k] = v
				if k == tokens.V.Data() {
					inheritsValue = true
				}
			}
		}
	}

	for k, v := range fieldDict.Data() {
		inheritedMap[k] = v
		if k == tokens.V.Data() {
			inheritsValue = false
		}
	}

	result, _ := tokens.WithMap(inheritedMap)
	return result, inheritsValue
}

func isChoiceSelected(selectedOptionNames []string, selectedOptionIndices []int, index int, name string) bool {
	if len(selectedOptionNames) == 0 {
		return false
	}

	for _, optionName := range selectedOptionNames {
		if optionName != name {
			continue
		}

		if selectedOptionIndices == nil {
			return true
		}

		for _, idx := range selectedOptionIndices {
			if idx == index {
				return true
			}
		}

		return false
	}

	return false
}
