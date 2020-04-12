package common

type File struct {
	Path     string
	Template bool
	Content  []byte
}

type InputVariables map[string]interface{}
