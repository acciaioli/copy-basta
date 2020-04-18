package specification

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SpecVariable_valueOK(t *testing.T) {
	tests := []struct {
		name    string
		specVar SpecVariable
		value   interface{}
	}{
		{
			name:    "string to string",
			specVar: SpecVariable{Type: openAPIString},
			value:   "a string",
		},
		{
			name:    "int to number",
			specVar: SpecVariable{Type: openAPINumber},
			value:   10,
		},
		{
			name:    "float to number",
			specVar: SpecVariable{Type: openAPINumber},
			value:   10.2,
		},
		{
			name:    "int to integer",
			specVar: SpecVariable{Type: openAPIInteger},
			value:   11,
		},
		{
			name:    "bool to boolean",
			specVar: SpecVariable{Type: openAPIBoolean},
			value:   true,
		},
		{
			name:    "slice to array",
			specVar: SpecVariable{Type: openAPIArray},
			value:   []interface{}{"hello", 12},
		},
		{
			name:    "map to object",
			specVar: SpecVariable{Type: openAPIObject},
			value:   map[string]interface{}{"string": "value", "integer": 13},
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
		specVar SpecVariable
		value   interface{}
	}{
		{
			name:    "int to string",
			specVar: SpecVariable{Type: openAPIString},
			value:   4,
		},
		{
			name:    "string to number",
			specVar: SpecVariable{Type: openAPINumber},
			value:   "not a number",
		},
		{
			name:    "bool to integer",
			specVar: SpecVariable{Type: openAPIInteger},
			value:   false,
		},
		{
			name:    "float to boolean",
			specVar: SpecVariable{Type: openAPIBoolean},
			value:   9.3,
		},
		{
			name:    "map to array",
			specVar: SpecVariable{Type: openAPIArray},
			value:   map[string]interface{}{"bool": true},
		},
		{
			name:    "map to object",
			specVar: SpecVariable{Type: openAPIObject},
			value:   []interface{}{"bye", 934},
		},
		{
			name:    "unknown type",
			specVar: SpecVariable{Type: openAPIType("unknown")},
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
		specVar SpecVariable
	}{
		{
			name: "simple",
			specVar: SpecVariable{
				Name:        "simple",
				Type:        openAPIString,
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "complete",
			specVar: SpecVariable{
				Name:        "complete",
				Type:        openAPIInteger,
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
		specVar SpecVariable
	}{
		{
			name: "missing name",
			specVar: SpecVariable{
				Type:        openAPIBoolean,
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "missing type",
			specVar: SpecVariable{
				Name:        "myName",
				Type:        "",
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "invalid type",
			specVar: SpecVariable{
				Name:        "myName",
				Type:        openAPIType("notValid"),
				Default:     nil,
				Description: nil,
			},
		},
		{
			name: "invalid default",
			specVar: SpecVariable{
				Name:        "myName",
				Type:        openAPIBoolean,
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
	tests := []struct {
		name       string
		specVar    SpecVariable
		expectedIn []string
	}{
		{
			name: "simple",
			specVar: SpecVariable{
				Name:        "myVariable",
				Type:        "myType",
				Default:     nil,
				Description: nil,
			},
			expectedIn: []string{"myType", "?", "myVariable"},
		},
		{
			name: "with default",
			specVar: SpecVariable{
				Name:        "myVariable",
				Type:        "myType",
				Default:     "myDefault",
				Description: nil,
			},
			expectedIn: []string{"myType", "?", "myVariable", "myDefault"},
		},
		{
			name: "with description",
			specVar: SpecVariable{
				Name:        "myVariable",
				Type:        "myType",
				Default:     nil,
				Description: func() *string { s := "my template variable description"; return &s }(),
			},
			expectedIn: []string{"my template variable description", "myType", "?", "myVariable"},
		},
		{
			name: "with default and description",
			specVar: SpecVariable{
				Name:        "myVariable",
				Type:        "myType",
				Default:     "myDefault",
				Description: func() *string { s := "my template variable description"; return &s }(),
			},
			expectedIn: []string{"my template variable description", "myType", "?", "myVariable", "myDefault"},
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
		specVar       SpecVariable
		text          string
		expectedValue interface{}
	}{
		{
			name:          "string",
			specVar:       SpecVariable{Type: openAPIString},
			text:          "a string",
			expectedValue: "a string",
		},
		{
			name:          "number",
			specVar:       SpecVariable{Type: openAPINumber},
			text:          "42.1",
			expectedValue: 42.1,
		},
		{
			name:          "integer",
			specVar:       SpecVariable{Type: openAPIInteger},
			text:          "73",
			expectedValue: 73,
		},
		{
			name:          "boolean",
			specVar:       SpecVariable{Type: openAPIBoolean},
			text:          "true",
			expectedValue: true,
		},
		{
			name:          "slice",
			specVar:       SpecVariable{Type: openAPIArray},
			text:          "eleven,12",
			expectedValue: []string{"eleven", "12"},
		},
		{
			name:          "map",
			specVar:       SpecVariable{Type: openAPIObject},
			text:          "key1=value1,key2=22,key3=false",
			expectedValue: map[string]string{"key1": "value1", "key2": "22", "key3": "false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.specVar.process(tt.text)
			require.Nil(t, err)
			require.Equal(t, tt.expectedValue, value)
		})
	}
}
