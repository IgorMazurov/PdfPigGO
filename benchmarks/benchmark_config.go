// Package benchmarks provides benchmark utilities for PdfPigGO.
//
// The C# project uses BenchmarkDotNet with a NuGetPackageConfig that sets up
// separate benchmark jobs for "Local" and "Latest" builds on .NET 8.0.
// Go replaces this with its native testing.B framework — benchmarks are run
// via `go test -bench=.` and require no separate configuration file.
//
// Constants below are retained for reference by other benchmark files.
package benchmarks

const (
	// Local is the identifier for benchmarks against a locally built version.
	Local = "Local"

	// Latest is the identifier for benchmarks against the latest released version.
	Latest = "Latest"
)
