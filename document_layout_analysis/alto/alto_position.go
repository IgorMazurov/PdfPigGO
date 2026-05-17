package alto

// AltoPosition represents the position of a page in ALTO documents.
type AltoPosition string

func (p AltoPosition) String() string { return string(p) }

const (
	AltoPositionLeft    AltoPosition = "Left"
	AltoPositionRight   AltoPosition = "Right"
	AltoPositionFoldout AltoPosition = "Foldout"
	AltoPositionSingle  AltoPosition = "Single"
	AltoPositionCover   AltoPosition = "Cover"
)
