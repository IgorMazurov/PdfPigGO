package page

// PageXmlRegionItem is the common interface for all PAGE XML region types,
// enabling polymorphic collections of child regions. In the C# source this
// corresponds to the abstract base class PageXmlRegion whose derived types
// are listed via XmlInclude attributes:
//   - AdvertRegion
//   - ChartRegion
//   - ChemRegion
//   - CustomRegion
//   - GraphicRegion
//   - ImageRegion
//   - LineDrawingRegion
//   - MapRegion
//   - MathsRegion
//   - MusicRegion
//   - NoiseRegion
//   - SeparatorRegion
//   - TableRegion
//   - TextRegion
//   - UnknownRegion
type PageXmlRegionItem interface {
	regionItem()
}

// All concrete region types share the following common fields (defined in each
// struct rather than via embedding, following Go conventions for XML serialization):
//
//	AlternativeImages []AlternativeImage  // Alternative region images (e.g., black-and-white)
//	Coords            *PageXmlCoords      // Polygon outline of the region
//	UserDefined       []PageXmlUserAttribute // Structured custom data
//	Labels            []PageXmlLabels     // Semantic labels/tags
//	Roles             *PageXmlRoles       // Roles in context of a parent region
//	Items             []PageXmlRegionItem // Child regions nested within this region
//	Id                string              // Unique identifier for the region
//	Custom            string              // Free-form custom attribute
//	Comments          string              // Optional annotation text
//	Continuation      bool                // Whether this is a continuation of another region
