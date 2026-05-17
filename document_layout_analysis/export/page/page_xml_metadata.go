package page

import "time"

// PageXmlMetadata represents the Metadata element in PAGE XML documents,
// containing creator information, timestamps, comments, user-defined attributes,
// and metadata items.
type PageXmlMetadata struct {
	// Creator identifies the software or person that created the document.
	Creator string `xml:"creator,omitempty"`

	// Created is the timestamp when the document was created (UTC).
	Created time.Time `xml:"created,omitempty"`

	// LastChange is the timestamp of the last modification (UTC).
	LastChange time.Time `xml:"lastChange,omitempty"`

	// Comments contains free-form annotation text.
	Comments string `xml:"comments,omitempty"`

	// UserDefined holds structured custom data defined by name, type and value.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// MetadataItems is a collection of metadata entries (e.g., author, image properties).
	MetadataItems []PageXmlMetadataItem `xml:"MetadataItem"`

	// ExternalRef is an external reference of any kind.
	ExternalRef string `xml:"externalRef,attr,omitempty"`
}
