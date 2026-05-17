package page

// PageXmlRelations is a container for one-to-one relations between layout
// objects (for example: DropCap - paragraph, caption - image).
type PageXmlRelations struct {
	// Relations holds the collection of relation elements.
	Relations []*PageXmlRelation `xml:"Relation"`
}
