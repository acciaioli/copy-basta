package common

import "os"

type CommandFlag struct {
	Ref     *string
	Name    string
	Default *string
	Usage   string
}

type File struct {
	Path     string
	Mode     os.FileMode
	Template bool
	Content  []byte
}

type InputVariables map[string]interface{}
