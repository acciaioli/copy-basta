package specification

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func newSpecVar(t openAPIType) Variable {
	return Variable{
		Type: &t,
	}
}

func Test_SpecVariable_valueOK(t *testing.T) {
	tests := []struct {
		name    string
		specVar Variable
		value   interface{}
	}{
		{
			name:    "string to string",
			specVar: newSpecVar(openAPIString),
			value:   "a string",
		},
		{
			name:    "int to number",
			specVar: newSpecVar(openAPINumber),
			value:   10,
		},
		{
			name:    "float to number",
			specVar: newSpecVar(openAPINumber),
			value:   10.2,
		},
		{
			name:    "int to integer",
			specVar: newSpecVar(openAPIInteger),
			value:   11,
		},
		{
			name:    "bool to boolean",
			specVar: newSpecVar(openAPIBoolean),
			value:   true,
		},
		{
			name:    "slice to array",
			specVar: newSpecVar(openAPIArray),
			value:   []interface{}{"hello", 12},
		},
		{
			name:    "map to object",
			specVar: newSpecVar(openAPIObject),
			value:   map[string]interface{}{"string": "value", "integer": 13},
		},
		{
			name: "missing type",
			specVar: Variable{
				Type: nil,
			},
			value: "any value would do",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.specVar.valueOk(tt.value)
			require.Nil(t, err)
		})
	}
}

func Test_SpecVariable_valueOK_error(t *testing.T) {
	tests := []struct {
		name    string
		specVar Variable
		value   interface{}
	}{
		{
			name:    "int to string",
			specVar: newSpecVar(openAPIString),
			value:   4,
		},
		{
			name:    "string to number",
			specVar: newSpecVar(openAPINumber),
			value:   "not a number",
		},
		{
			name:    "bool to integer",
			specVar: newSpecVar(openAPIInteger),
			value:   false,
		},
		{
			name:    "float to boolean",
			specVar: newSpecVar(openAPIBoolean),
			value:   9.3,
		},
		{
			name:    "map to array",
			specVar: newSpecVar(openAPIArray),
			value:   map[string]interface{}{"bool": true},
		},
		{
			name:    "map to object",
			specVar: newSpecVar(openAPIObject),
			value:   []interface{}{"bye", 934},
		},
		{
			name:    "unknown type",
			specVar: newSpecVar("unknown"),
			value:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.specVar.valueOk(tt.value)
			require.NotNil(t, err)
		})
	}
}

func Test_SpecVariable_validate(t *testing.T) {
	tests := []struct {
		name    string
		specVar Variable
	}{
		{
			name: "simple",
			specVar: Variable{
				Name:        "simple",
				Type:        nil,
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "complete",
			specVar: Variable{
				Name:        "complete",
				Type:        func() *openAPIType { v := openAPIInteger; return &v }(),
				Default:     2289,
				Description: func() *string { s := "a legit integer"; return &s }(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.specVar.validate()
			require.Nil(t, err)
		})
	}
}

func Test_SpecVariable_validate_error(t *testing.T) {
	tests := []struct {
		name    string
		specVar Variable
	}{
		{
			name: "missing name",
			specVar: Variable{
				Type:        func() *openAPIType { v := openAPIBoolean; return &v }(),
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "invalid type",
			specVar: Variable{
				Name:        "myName",
				Type:        func() *openAPIType { v := openAPIType("notValid"); return &v }(),
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "invalid default",
			specVar: Variable{
				Name:        "myName",
				Type:        func() *openAPIType { v := openAPIBoolean; return &v }(),
				Default:     44,
				Description: func() *string { s := "a boolean, therefore not a integer"; return &s }(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.specVar.validate()
			require.NotNil(t, err)
		})
	}
}

func Test_SpecVariable_prompt(t *testing.T) {
	myType := openAPIType("myType")

	tests := []struct {
		name       string
		specVar    Variable
		expectedIn []string
	}{
		{
			name: "simple",
			specVar: Variable{
				Name:        "myVariable",
				Type:        &myType,
				Default:     nil,
				Description: nil,
			},
			expectedIn: []string{"myType", "?", "myVariable"},
		},
		{
			name: "with default",
			specVar: Variable{
				Name:        "myVariable",
				Type:        &myType,
				Default:     "myDefault",
				Description: nil,
			},
			expectedIn: []string{"myType", "?", "myVariable", "myDefault"},
		},
		{
			name: "with description",
			specVar: Variable{
				Name:        "myVariable",
				Type:        &myType,
				Default:     nil,
				Description: func() *string { s := "my template variable description 1"; return &s }(),
			},
			expectedIn: []string{"my template variable description 1", "myType", "?", "myVariable"},
		},
		{
			name: "with default and description",
			specVar: Variable{
				Name:        "myVariable",
				Type:        &myType,
				Default:     "myDefault",
				Description: func() *string { s := "my template variable description 2"; return &s }(),
			},
			expectedIn: []string{"my template variable description 2", "myType", "?", "myVariable", "myDefault"},
		},
		{
			name: "no type",
			specVar: Variable{
				Name:        "myVariable",
				Type:        nil,
				Default:     "myDefault",
				Description: func() *string { s := "my template variable description 3"; return &s }(),
			},
			expectedIn: []string{"my template variable description 3", "?", "myVariable", "myDefault"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := tt.specVar.prompt()
			for _, s := range tt.expectedIn {
				require.Contains(t, text, s)
			}

		})
	}
}

func Test_SpecVariable_process(t *testing.T) {
	tests := []struct {
		name          string
		specVar       Variable
		text          string
		expectedValue interface{}
	}{
		{
			name:          "string",
			specVar:       newSpecVar(openAPIString),
			text:          "a string",
			expectedValue: "a string",
		},
		{
			name:          "number",
			specVar:       newSpecVar(openAPINumber),
			text:          "42.1",
			expectedValue: 42.1,
		},
		{
			name:          "integer",
			specVar:       newSpecVar(openAPIInteger),
			text:          "73",
			expectedValue: 73,
		},
		{
			name:          "boolean",
			specVar:       newSpecVar(openAPIBoolean),
			text:          "true",
			expectedValue: true,
		},
		{
			name:          "slice",
			specVar:       newSpecVar(openAPIArray),
			text:          "eleven,12",
			expectedValue: []string{"eleven", "12"},
		},
		{
			name:          "map",
			specVar:       newSpecVar(openAPIObject),
			text:          "key1=value1,key2=22,key3=false",
			expectedValue: map[string]string{"key1": "value1", "key2": "22", "key3": "false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.specVar.fromString(tt.text)
			require.Nil(t, err)
			require.Equal(t, tt.expectedValue, value)
		})
	}
}
