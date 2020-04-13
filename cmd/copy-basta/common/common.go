package common

type CommandFlag struct {
	Ref     *string
	Name    string
	Default *string
	Usage   string
}

type File struct {
	Path     string
	Template bool
	Content  []byte
}

type InputVariables map[string]interface{}
