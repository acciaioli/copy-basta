package specification

type SpecData struct {
	Ignore      []string       `yaml:"ignore"`
	PassThrough []string       `yaml:"pass-through"`
	Variables   []VariableData `yaml:"variables"`
	OnOverwrite OnOverwrite    `yaml:"on-overwrite"`
}

type VariableData struct {
	Name        string      `yaml:"name"`
	DType       *string     `yaml:"type"`
	DefaultVal  interface{} `yaml:"default"`
	Description *string     `yaml:"description"`
}

type OnOverwrite struct {
	Exclude []string `yaml:"exclude"`
}
