package adobe_font_metrics

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/fonts/encodings"
)

// AdobeFontMetricsEncoding is an encoding derived from an Adobe Font Metrics file.
type AdobeFontMetricsEncoding struct {
	*encodings.Encoding
}

// NewAdobeFontMetricsEncoding creates a new encoding from AFM metrics.
func NewAdobeFontMetricsEncoding(metrics *AdobeFontMetrics) (*AdobeFontMetricsEncoding, error) {
	if metrics == nil {
		return nil, errors.New("metrics cannot be nil")
	}

	e := &AdobeFontMetricsEncoding{
		Encoding: encodings.NewEncoding(),
	}

	for name, characterMetric := range metrics.CharacterMetrics {
		e.Add(characterMetric.CharacterCode, name)
	}

	return e, nil
}

// EncodingName returns the name of this encoding.
func (e *AdobeFontMetricsEncoding) EncodingName() string {
	return "AFM"
}
