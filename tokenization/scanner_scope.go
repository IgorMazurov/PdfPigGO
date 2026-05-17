package tokenization

// ScannerScope represents the current scope of the token scanner.
type ScannerScope int

const (
	// ScannerScopeNone indicates reading normally, outside any container.
	ScannerScopeNone ScannerScope = 0

	// ScannerScopeArray indicates reading inside an array.
	ScannerScopeArray ScannerScope = 1

	// ScannerScopeDictionary indicates reading inside a dictionary.
	ScannerScopeDictionary ScannerScope = 2
)
