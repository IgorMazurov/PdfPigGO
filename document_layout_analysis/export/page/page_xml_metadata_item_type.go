package page

// PageXmlMetadataItemType represents the type of a metadata item.
type PageXmlMetadataItemType byte

const (
	// MetaItemAuthor indicates the metadata describes an author.
	MetaItemAuthor PageXmlMetadataItemType = iota

	// MetaItemImageProperties indicates the metadata describes image properties.
	MetaItemImageProperties

	// MetaItemProcessingStep indicates the metadata describes a processing step.
	MetaItemProcessingStep

	// MetaItemOther indicates miscellaneous metadata not fitting other categories.
	MetaItemOther
)
