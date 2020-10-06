package specification

type specData struct {
	Ignore      []string       `yaml:"ignore"`
	PassThrough []string       `yaml:"pass-through"`
	Variables   []variableData `yaml:"variables"`
}

type variableData struct {
	Name        string      `yaml:"name"`
	DType       *string     `yaml:"type"`
	DefaultVal  interface{} `yaml:"default"`
	Description *string     `yaml:"description"`
}
