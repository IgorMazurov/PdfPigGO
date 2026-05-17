package document_layout_analysis

// DlaOptions stores options that configure the operation of methods of the
// document layout analysis algorithm.
type DlaOptions interface {
	// MaxDegreeOfParallelism returns the maximum number of concurrent tasks enabled.
	// A positive value limits the number of concurrent operations to the set value.
	// If it is -1, there is no limit on the number of concurrently running operations.
	MaxDegreeOfParallelism() int

	// SetMaxDegreeOfParallelism sets the maximum number of concurrent tasks enabled.
	SetMaxDegreeOfParallelism(value int)
}
