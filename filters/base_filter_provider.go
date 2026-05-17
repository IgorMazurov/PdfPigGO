package filters

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// BaseFilterProvider is the base implementation of FilterProvider that resolves
// filter names from a dictionary of known filter instances.
type BaseFilterProvider struct {
	filterInstances map[string]Filter
}

// NewBaseFilterProvider creates a new BaseFilterProvider with the given filters.
func NewBaseFilterProvider(filterInstances map[string]Filter) *BaseFilterProvider {
	return &BaseFilterProvider{
		filterInstances: filterInstances,
	}
}

// GetFilters returns the filters corresponding to the /Filter or /F entry in the
// given dictionary. The value may be a single name token or an array of name tokens.
func (p *BaseFilterProvider) GetFilters(dictionary *tokens.DictionaryToken) ([]Filter, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}

	token := getFilterOrFallback(dictionary)
	if token == nil {
		return []Filter{}, nil
	}

	switch t := token.(type) {
	case *tokens.ArrayToken:
		result := make([]Filter, len(t.Data()))
		for i, filterToken := range t.Data() {
			nameToken, ok := filterToken.(*tokens.NameToken)
			if !ok {
				return nil, core.NewPdfDocumentFormatException(
					fmt.Sprintf("filter array element is not a name token: %v", filterToken))
			}
			filter, err := p.getFilterStrict(nameToken.Data())
			if err != nil {
				return nil, err
			}
			result[i] = filter
		}
		return result, nil
	case *tokens.NameToken:
		filter, err := p.getFilterStrict(t.Data())
		if err != nil {
			return nil, err
		}
		return []Filter{filter}, nil
	default:
		return nil, core.NewPdfDocumentFormatException(
			fmt.Sprintf("the filter for the stream was not a valid object. Expected name or array, instead got: %v", token))
	}
}

// GetNamedFilters returns the filters corresponding to the given filter names.
func (p *BaseFilterProvider) GetNamedFilters(names []*tokens.NameToken) ([]Filter, error) {
	result := make([]Filter, 0, len(names))

	for _, name := range names {
		filter, err := p.getFilterStrict(name.Data())
		if err != nil {
			return nil, err
		}
		result = append(result, filter)
	}

	return result, nil
}

// GetAllFilters returns all registered filters in the provider.
func (p *BaseFilterProvider) GetAllFilters() []Filter {
	seen := make(map[Filter]bool)
	var result []Filter

	for _, f := range p.filterInstances {
		if !seen[f] {
			seen[f] = true
			result = append(result, f)
		}
	}

	return result
}

func (p *BaseFilterProvider) getFilterStrict(name string) (Filter, error) {
	filter, ok := p.filterInstances[name]
	if !ok {
		return nil, fmt.Errorf("the filter with the name %s is not supported yet. Please raise an issue", name)
	}
	return filter, nil
}

func getFilterOrFallback(dict *tokens.DictionaryToken) tokens.Token {
	token, ok := dict.TryGet(tokens.Filter)
	if ok {
		return token
	}
	token, ok = dict.TryGet(tokens.F)
	if ok {
		return token
	}
	return nil
}
