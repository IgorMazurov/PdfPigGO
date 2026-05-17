package charstrings

import (
	"log"
)

// LazyType2Command represents the deferred execution of a Type 2 charstring command.
type LazyType2Command struct {
	name                 string
	minimumStackParameters int
	runCommand           func(*Type2BuildCharContext)
}

// NewLazyType2Command creates a new LazyType2Command.
// name is the command name per the Type 2 charstring specification.
// minimumStackParameters is the minimum number of arguments required on the stack, or -1 to skip checking.
// runCommand is the action to execute when evaluating the command, modifying the context.
func NewLazyType2Command(name string, minimumStackParameters int, runCommand func(*Type2BuildCharContext)) *LazyType2Command {
	if name == "" {
		panic("name cannot be empty")
	}
	if runCommand == nil {
		panic("runCommand cannot be nil")
	}
	return &LazyType2Command{
		name:                 name,
		minimumStackParameters: minimumStackParameters,
		runCommand:           runCommand,
	}
}

// Name returns the command name. See the Type 2 charstring specification for possible command names.
func (c *LazyType2Command) Name() string {
	return c.name
}

// Run evaluates the command against the given context.
// If the stack has fewer elements than minimumStackParameters, a warning is logged,
// the stack is cleared, and the command is not executed.
func (c *LazyType2Command) Run(context *Type2BuildCharContext) {
	if context == nil {
		panic("context cannot be nil")
	}

	stackLen := context.Stack().Length()
	if c.minimumStackParameters >= 0 && stackLen < c.minimumStackParameters {
		log.Printf("Warning: CFF CharString command '%s' expected %d arguments. Got: %d. Command ignored and stack cleared.",
			c.name, c.minimumStackParameters, stackLen)
		context.Stack().Clear()
		return
	}

	c.runCommand(context)
}

// String returns the command name.
func (c *LazyType2Command) String() string {
	return c.name
}
