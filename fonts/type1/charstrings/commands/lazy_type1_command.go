package commands

import "fmt"

// LazyType1Command represents the deferred execution of a Type 1 charstring command.
type LazyType1Command struct {
	name       string
	runCommand func(*Type1BuildCharContext)
}

// NewLazyType1Command creates a new LazyType1Command.
// name is the command name per the Type 1 charstring specification.
// runCommand is the action to execute when evaluating the command, modifying the context.
func NewLazyType1Command(name string, runCommand func(*Type1BuildCharContext)) *LazyType1Command {
	if name == "" {
		panic("name cannot be empty")
	}
	if runCommand == nil {
		panic("runCommand cannot be nil")
	}
	return &LazyType1Command{
		name:       name,
		runCommand: runCommand,
	}
}

// Name returns the command name.
func (c *LazyType1Command) Name() string {
	return c.name
}

// Run evaluates the command against the given context.
func (c *LazyType1Command) Run(context *Type1BuildCharContext) {
	if context == nil {
		panic("context cannot be nil")
	}
	c.runCommand(context)
}

// String returns the command name.
func (c *LazyType1Command) String() string {
	return c.name
}

// _ ensures LazyType1Command implements fmt.Stringer.
var _ fmt.Stringer = (*LazyType1Command)(nil)
