package specification

type SpecData struct {
	Ignore      []string       `yaml:"ignore,omitempty"`
	PassThrough []string       `yaml:"pass-through,omitempty"`
	Variables   []VariableData `yaml:"variables"`
}

type VariableData struct {
	Name        string      `yaml:"name,omitempty"`
	DType       *string     `yaml:"type,omitempty"`
	DefaultVal  interface{} `yaml:"default,omitempty"`
	Description *string     `yaml:"description,omitempty"`
}
