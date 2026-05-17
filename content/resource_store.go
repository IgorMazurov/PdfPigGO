package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/util"
)

// ResourceStore manages the resources associated with a PDF content stream,
// including fonts, XObjects, color spaces, patterns, and shadings.
type ResourceStore interface {
	// LoadResourceDictionary loads the given resource dictionary into the store.
	LoadResourceDictionary(resourceDictionary *tokens.DictionaryToken)

	// UnloadResourceDictionary removes any named resources and associated state
	// for the last resource dictionary loaded. Does not affect cached resources,
	// just the labels associated with them.
	UnloadResourceDictionary()

	// GetFont returns the font corresponding to the given name.
	GetFont(name *tokens.NameToken) fonts.Font

	// TryGetXObject attempts to get the XObject stream corresponding to the name.
	// Returns true and the stream if found, false and nil otherwise.
	TryGetXObject(name *tokens.NameToken) (*tokens.StreamToken, bool)

	// GetExtendedGraphicsStateDictionary returns the extended graphics state
	// dictionary corresponding to the name, or nil if not found.
	GetExtendedGraphicsStateDictionary(name *tokens.NameToken) *tokens.DictionaryToken

	// GetFontDirectly returns the font from the indirect reference token.
	GetFontDirectly(fontReferenceToken *tokens.IndirectReferenceToken) fonts.Font

	// TryGetNamedColorSpace attempts to get the named color space by its name.
	// Returns true and the color space if found, false and a zero value otherwise.
	TryGetNamedColorSpace(name *tokens.NameToken) (colors.ResourceColorSpace, bool)

	// GetColorSpaceDetails returns the color space details corresponding to the
	// name and dictionary.
	GetColorSpaceDetails(name *tokens.NameToken, dictionary *tokens.DictionaryToken) colors.ColorSpaceDetails

	// GetMarkedContentPropertiesDictionary returns the marked content properties
	// dictionary corresponding to the name, or nil if not found.
	GetMarkedContentPropertiesDictionary(name *tokens.NameToken) *tokens.DictionaryToken

	// GetPatterns returns all pattern colors as a map keyed by their names.
	GetPatterns() map[*tokens.NameToken]colors.PatternColor

	// GetShading returns the shading corresponding to the name.
	GetShading(name *tokens.NameToken) colors.Shading
}

// stackDict is a simple stack-based dictionary that supports Push/Pop
// for scoped resource management. Each push creates a new scope layer; set
// writes to the topmost layer; tryGet reads through all layers from top down.
type stackDict[K comparable, V any] struct {
	layers []map[K]V
}

func newStackDict[K comparable, V any]() *stackDict[K, V] {
	return &stackDict[K, V]{
		layers: []map[K]V{{}},
	}
}

func (s *stackDict[K, V]) push() {
	s.layers = append(s.layers, make(map[K]V))
}

func (s *stackDict[K, V]) pop() {
	if len(s.layers) > 1 {
		s.layers = s.layers[:len(s.layers)-1]
	}
}

func (s *stackDict[K, V]) set(key K, value V) {
	if len(s.layers) > 0 {
		s.layers[len(s.layers)-1][key] = value
	}
}

func (s *stackDict[K, V]) tryGet(key K) (V, bool) {
	for i := len(s.layers) - 1; i >= 0; i-- {
		if val, ok := s.layers[i][key]; ok {
			return val, true
		}
	}
	var zero V
	return zero, false
}

// resourceStoreImpl is the concrete implementation of ResourceStore.
type resourceStoreImpl struct {
	scanner        tokenization.PdfTokenScanner
	fontFactory    pdffonts.FontFactory
	filterProvider util.PatternParserProvider
	parsingOptions *ParsingOptions

	loadedFonts                  map[string]fonts.Font
	loadedDirectFonts            map[*tokens.NameToken]fonts.Font
	currentFontState             *stackDict[*tokens.NameToken, core.IndirectReference]
	currentXObjectState          *stackDict[*tokens.NameToken, core.IndirectReference]
	extendedGraphicsStates       *stackDict[*tokens.NameToken, *tokens.DictionaryToken]
	namedColorSpaces             *stackDict[*tokens.NameToken, colors.ResourceColorSpace]
	loadedNamedColorSpaceDetails map[*tokens.NameToken]colors.ColorSpaceDetails

	markedContentProperties map[*tokens.NameToken]*tokens.DictionaryToken
	shadingsProperties      map[*tokens.NameToken]colors.Shading
	patternsProperties     map[*tokens.NameToken]colors.PatternColor

	lastLoadedFontName *tokens.NameToken
	lastLoadedFont     fonts.Font
}

// NewResourceStore creates a new ResourceStore with the given dependencies.
func NewResourceStore(
	scanner tokenization.PdfTokenScanner,
	fontFactory pdffonts.FontFactory,
	filterProvider util.PatternParserProvider,
	parsingOptions *ParsingOptions,
) ResourceStore {
	return &resourceStoreImpl{
		scanner:        scanner,
		fontFactory:    fontFactory,
		filterProvider: filterProvider,
		parsingOptions: parsingOptions,

		loadedFonts:                  make(map[string]fonts.Font),
		loadedDirectFonts:            make(map[*tokens.NameToken]fonts.Font),
		currentFontState:             newStackDict[*tokens.NameToken, core.IndirectReference](),
		currentXObjectState:          newStackDict[*tokens.NameToken, core.IndirectReference](),
		extendedGraphicsStates:       newStackDict[*tokens.NameToken, *tokens.DictionaryToken](),
		namedColorSpaces:             newStackDict[*tokens.NameToken, colors.ResourceColorSpace](),
		loadedNamedColorSpaceDetails: make(map[*tokens.NameToken]colors.ColorSpaceDetails),

		markedContentProperties: make(map[*tokens.NameToken]*tokens.DictionaryToken),
		shadingsProperties:      make(map[*tokens.NameToken]colors.Shading),
		patternsProperties:     make(map[*tokens.NameToken]colors.PatternColor),
	}
}

// LoadResourceDictionary loads the given resource dictionary into the store.
func (r *resourceStoreImpl) LoadResourceDictionary(resourceDictionary *tokens.DictionaryToken) {
	r.lastLoadedFontName = nil
	r.lastLoadedFont = nil
	r.loadedNamedColorSpaceDetails = make(map[*tokens.NameToken]colors.ColorSpaceDetails)

	r.namedColorSpaces.push()
	r.currentFontState.push()
	r.currentXObjectState.push()
	r.extendedGraphicsStates.push()

	// Load fonts
	if fontBase, found := resourceDictionary.TryGet(tokens.Font); found {
		fontDict := resolveToDictionary(fontBase, r.scanner)
		if fontDict != nil {
			r.loadFontDictionary(fontDict)
		}
	}

	// Load XObjects
	if xobjectBase, found := resourceDictionary.TryGet(tokens.Xobject); found {
		xobjectDict := resolveToDictionary(xobjectBase, r.scanner)
		if xobjectDict != nil {
			for keyStr, val := range xobjectDict.Data() {
				if _, isNull := val.(*tokens.NullToken); isNull {
					continue
				}
				refTok, ok := val.(*tokens.IndirectReferenceToken)
				if !ok {
					r.parsingOptions.Logger.Error(fmt.Sprintf("Expected the XObject dictionary value for key /%s to be an indirect reference, instead got: %v.", keyStr, val))
					continue
				}
				r.currentXObjectState.set(tokens.Create(keyStr), refTok.Data())
			}
		}
	}

	// Load Extended Graphics States
	if extGStateBase, found := resourceDictionary.TryGet(tokens.ExtGState); found {
		extGStateDict := resolveToDictionary(extGStateBase, r.scanner)
		if extGStateDict != nil {
			for keyStr, val := range extGStateDict.Data() {
				name := tokens.Create(keyStr)
				state := resolveToDictionary(val, r.scanner)
				if state != nil {
					r.extendedGraphicsStates.set(name, state)
				}
			}
		}
	}

	// Load Color Spaces
	if csBase, found := resourceDictionary.TryGet(tokens.ColorSpace); found {
		csDict := resolveToDictionary(csBase, r.scanner)
		if csDict != nil {
			for keyStr, val := range csDict.Data() {
				name := tokens.Create(keyStr)

				if csName, ok := val.(*tokens.NameToken); ok {
					r.namedColorSpaces.set(name, colors.ResourceColorSpace{Name: csName})
				} else if csArray, ok := parts.TryGet[*tokens.ArrayToken](val, r.scanner); ok {
					data := csArray.Data()
					if len(data) == 0 {
						r.parsingOptions.Logger.Error(fmt.Sprintf("Empty ColorSpace array encountered in page resource dictionary: %v.", resourceDictionary))
						continue
					}
					firstName, ok := data[0].(*tokens.NameToken)
					if !ok {
						r.parsingOptions.Logger.Error(fmt.Sprintf("Invalid ColorSpace array encountered in page resource dictionary: %v.", csArray))
						continue
					}
					r.namedColorSpaces.set(name, colors.ResourceColorSpace{Name: firstName, Data: csArray})
				} else if r.parsingOptions.UseLenientParsing {
					if dict := resolveToDictionary(val, r.scanner); dict != nil {
						if csNameTok, found := dict.TryGet(tokens.ColorSpace); found {
							if csName, ok := csNameTok.(*tokens.NameToken); ok {
								r.namedColorSpaces.set(name, colors.ResourceColorSpace{Name: csName})
							}
						}
					}
				} else {
					r.parsingOptions.Logger.Error(fmt.Sprintf("Invalid ColorSpace token encountered in page resource dictionary: %v.", val))
				}
			}
		}
	}

	// Load Patterns
	if patternBase, found := resourceDictionary.TryGet(tokens.Pattern); found {
		patternDict := resolveToDictionary(patternBase, r.scanner)
		if patternDict != nil {
			for keyStr, val := range patternDict.Data() {
				name := tokens.Create(keyStr)
				pattern, err := util.CreatePattern(val, r.scanner, r, r.filterProvider)
				if err != nil {
					r.parsingOptions.Logger.Error(fmt.Sprintf("Failed to load pattern /%s: %v", keyStr, err))
					continue
				}
				r.patternsProperties[name] = pattern
			}
		}
	}

	// Load Marked Content Properties
	if propsBase, found := resourceDictionary.TryGet(tokens.Properties); found {
		propsList := resolveToDictionary(propsBase, r.scanner)
		if propsList != nil {
			for keyStr, val := range propsList.Data() {
				keyName := tokens.Create(keyStr)
				namedProps := resolveToDictionary(val, r.scanner)
				if namedProps != nil {
					r.markedContentProperties[keyName] = namedProps
				}
			}
		}
	}

	// Load Shadings
	if shadingBase, found := resourceDictionary.TryGet(tokens.Shading); found {
		shadingDict := resolveToDictionary(shadingBase, r.scanner)
		if shadingDict != nil {
			for keyStr, val := range shadingDict.Data() {
				keyName := tokens.Create(keyStr)
				shading, err := util.CreateShading(val, r.scanner, r, r.filterProvider)
				if err != nil {
					r.parsingOptions.Logger.Error(fmt.Sprintf("Failed to load shading /%s: %v", keyStr, err))
					continue
				}
				r.shadingsProperties[keyName] = shading
			}
		}
	}
}

// UnloadResourceDictionary removes any named resources and associated state
// for the last resource dictionary loaded.
func (r *resourceStoreImpl) UnloadResourceDictionary() {
	r.lastLoadedFontName = nil
	r.lastLoadedFont = nil
	r.loadedNamedColorSpaceDetails = make(map[*tokens.NameToken]colors.ColorSpaceDetails)
	r.currentFontState.pop()
	r.currentXObjectState.pop()
	r.namedColorSpaces.pop()
	r.extendedGraphicsStates.pop()
}

func (r *resourceStoreImpl) loadFontDictionary(fontDict *tokens.DictionaryToken) {
	r.lastLoadedFontName = nil
	r.lastLoadedFont = nil

	for keyStr, val := range fontDict.Data() {
		if indRefTok, ok := val.(*tokens.IndirectReferenceToken); ok {
			ref := indRefTok.Data()
			name := tokens.Create(keyStr)
			r.currentFontState.set(name, ref)

			refKey := ref.String()
			if _, loaded := r.loadedFonts[refKey]; loaded {
				continue
			}

			fontObj := resolveToDictionary(indRefTok, r.scanner)
			if fontObj == nil {
				continue
			}

			var font fonts.Font
			var err error
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						err = fmt.Errorf("font parsing panicked: %v", rec)
					}
				}()
				font, err = r.fontFactory.Get(fontObj)
			}()
	if err != nil {
			if !r.parsingOptions.SkipMissingFonts {
				panic(err)
			}
			r.parsingOptions.Logger.Error(fmt.Sprintf("Failed to load font: %v", err))
			continue
		}
		r.loadedFonts[refKey] = font
		} else if dict, ok := val.(*tokens.DictionaryToken); ok {
			var font fonts.Font
			var err error
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						err = fmt.Errorf("font parsing panicked: %v", rec)
					}
				}()
				font, err = r.fontFactory.Get(dict)
			}()
			if err != nil {
				if !r.parsingOptions.SkipMissingFonts {
					panic(err)
				}
				r.parsingOptions.Logger.Error(fmt.Sprintf("Failed to load direct font: %v", err))
				continue
			}
			r.loadedDirectFonts[tokens.Create(keyStr)] = font
		}
	}
}

// GetFont returns the font corresponding to the given name.
func (r *resourceStoreImpl) GetFont(name *tokens.NameToken) fonts.Font {
	if r.lastLoadedFontName != nil && r.lastLoadedFontName.Data() == name.Data() {
		return r.lastLoadedFont
	}

	var font fonts.Font

	ref, found := r.currentFontState.tryGet(name)
	if found {
		refKey := ref.String()
		if f, ok := r.loadedFonts[refKey]; ok {
			font = f
		}
	} else if f, ok := r.loadedDirectFonts[name]; ok {
		font = f
	}

	r.lastLoadedFontName = name
	r.lastLoadedFont = font
	return font
}

// GetFontDirectly returns the font from the indirect reference token.
func (r *resourceStoreImpl) GetFontDirectly(fontReferenceToken *tokens.IndirectReferenceToken) fonts.Font {
	r.lastLoadedFontName = nil
	r.lastLoadedFont = nil

	fontDict := resolveToDictionary(fontReferenceToken, r.scanner)
	if fontDict == nil {
		r.parsingOptions.Logger.Error(fmt.Sprintf("The requested font reference token %v wasn't a font.", fontReferenceToken))
		return nil
	}

	font, err := r.fontFactory.Get(fontDict)
	if err != nil {
		r.parsingOptions.Logger.Error(fmt.Sprintf("Failed to get font directly: %v", err))
		return nil
	}
	return font
}

// TryGetNamedColorSpace attempts to get the named color space by its name.
func (r *resourceStoreImpl) TryGetNamedColorSpace(name *tokens.NameToken) (colors.ResourceColorSpace, bool) {
	if name == nil {
		return colors.ResourceColorSpace{}, false
	}
	cs, found := r.namedColorSpaces.tryGet(name)
	return cs, found
}

// GetColorSpaceDetails returns the color space details for the given name and dictionary.
func (r *resourceStoreImpl) GetColorSpaceDetails(name *tokens.NameToken, dictionary *tokens.DictionaryToken) colors.ColorSpaceDetails {
	if dictionary == nil {
		dictionary, _ = tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
	}

	// Null color space for images
	if name == nil {
		return util.GetColorSpaceDetails(nil, dictionary, r.scanner, r, r.filterProvider, false)
	}

	if cs, ok := colors.TryMapToColorSpace(name); ok {
		csPtr := &cs
		return util.GetColorSpaceDetails(csPtr, dictionary, r.scanner, r, r.filterProvider, false)
	}

	// Named color spaces — check cache first
	if cached, found := r.loadedNamedColorSpaceDetails[name]; found {
		return cached
	}

	if namedCS, found := r.TryGetNamedColorSpace(name); found {
		if mapped, ok := colors.TryMapToColorSpace(namedCS.Name); ok {
			mappedPtr := &mapped
			if namedCS.Data == nil {
				csd := util.GetColorSpaceDetails(mappedPtr, dictionary, r.scanner, r, r.filterProvider, false)
				r.loadedNamedColorSpaceDetails[name] = csd
				return csd
			}
			if arr, ok := namedCS.Data.(*tokens.ArrayToken); ok {
				csd := util.GetColorSpaceDetails(mappedPtr, dictionary.With(tokens.ColorSpace, arr), r.scanner, r, r.filterProvider, false)
				r.loadedNamedColorSpaceDetails[name] = csd
				return csd
			}
		}
	}

	r.parsingOptions.Logger.Error(fmt.Sprintf("Could not find color space for token '%v'.", name))
	return colors.UnsupportedColorSpaceDetails
}

// TryGetXObject attempts to get the XObject stream corresponding to the name.
func (r *resourceStoreImpl) TryGetXObject(name *tokens.NameToken) (*tokens.StreamToken, bool) {
	ref, found := r.currentXObjectState.tryGet(name)
	if !found {
		return nil, false
	}

	obj := r.scanner.Get(ref)
	if obj == nil {
		return nil, false
	}

	if stream, ok := obj.Data().(*tokens.StreamToken); ok {
		return stream, true
	}
	return nil, false
}

// GetExtendedGraphicsStateDictionary returns the extended graphics state dictionary.
func (r *resourceStoreImpl) GetExtendedGraphicsStateDictionary(name *tokens.NameToken) *tokens.DictionaryToken {
	dict, found := r.extendedGraphicsStates.tryGet(name)
	if found {
		return dict
	}

	r.parsingOptions.Logger.Error(fmt.Sprintf("The graphic state dictionary does not contain the key '%v'.", name))
	return nil
}

// GetMarkedContentPropertiesDictionary returns the marked content properties dictionary.
func (r *resourceStoreImpl) GetMarkedContentPropertiesDictionary(name *tokens.NameToken) *tokens.DictionaryToken {
	if dict, found := r.markedContentProperties[name]; found {
		return dict
	}
	return nil
}

// GetShading returns the shading corresponding to the name.
func (r *resourceStoreImpl) GetShading(name *tokens.NameToken) colors.Shading {
	if s, found := r.shadingsProperties[name]; found {
		return s
	}
	return nil
}

// GetPatterns returns all pattern colors as a map keyed by their names.
func (r *resourceStoreImpl) GetPatterns() map[*tokens.NameToken]colors.PatternColor {
	result := make(map[*tokens.NameToken]colors.PatternColor, len(r.patternsProperties))
	for k, v := range r.patternsProperties {
		result[k] = v
	}
	return result
}


