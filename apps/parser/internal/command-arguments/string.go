package command_arguments

type String struct {
	value    string
	Name     string
	Optional bool
	OneOf    []string
	Hint     string
}

var _ Arg = String{}

func (String) isCommandArg() {}

func (c String) Int() int {
	return 0
}

func (c String) String() string {
	return c.value
}
func (c String) GetName() string {
	return c.Name
}

func (c String) GetHint() string {
	if c.Hint == "" {
		return c.Name
	}

	return c.Hint
}

func (c String) IsOptional() bool {
	return c.Optional
}
