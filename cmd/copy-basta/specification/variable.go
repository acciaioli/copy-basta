package specification

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"copy-basta/cmd/copy-basta/common/log"

	"copy-basta/cmd/copy-basta/common"
)

/*
Spec variables use a subset of the OpenApi data models

https://swagger.io/docs/specification/data-models
*/

type openAPIType string

const (
	openAPIString  = openAPIType("string")
	openAPINumber  = openAPIType("number")
	openAPIInteger = openAPIType("integer")
	openAPIBoolean = openAPIType("boolean")
	openAPIArray   = openAPIType("array")
	openAPIObject  = openAPIType("object")
)

type Variable struct {
	Name        string       `yaml:"name"`
	Type        *openAPIType `yaml:"type"`
	Default     interface{}  `yaml:"default"`
	Description *string      `yaml:"description"`
}

func (v *Variable) validate() error {
	// name checks
	if v.Name == "" {
		return errors.New("variable error [name]: is required")
	}

	// type checks
	if v.Type != nil {
		if ok := func(actual openAPIType) bool {
			for _, candidate := range []openAPIType{
				openAPIString, openAPINumber, openAPIInteger, openAPIBoolean, openAPIArray, openAPIObject,
			} {
				if actual == candidate {
					return true
				}
			}
			return false
		}(*v.Type); !ok {
			return fmt.Errorf(`variable error [type]: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.Type)
		}
	} else {
		log.L.WarnWithData("spec variable without type, defaulting to any", log.Data{"name": v.Name})
	}

	// default checks
	if v.Default != nil {
		if err := v.valueOk(v.Default); err != nil {
			return fmt.Errorf("variable error [default]: %s", err.Error())
		}
	}

	return nil
}

func (v *Variable) valueOk(value interface{}) error {
	if v.Type == nil {
		return nil
	}

	actualKind := reflect.TypeOf(value).Kind()

	isOneOF := func(actual reflect.Kind, accepted []reflect.Kind) error {
		for _, acceptedKind := range accepted {
			if actual == acceptedKind {
				return nil
			}
		}
		format := "value error: decoded to type %v, expected one of %v. variable type is %s"
		return fmt.Errorf(format, actual, accepted, string(*v.Type))
	}

	var acceptedKinds []reflect.Kind

	switch *v.Type {
	case openAPIString:
		acceptedKinds = []reflect.Kind{reflect.String}
	case openAPINumber:
		acceptedKinds = []reflect.Kind{reflect.Int, reflect.Float64}
	case openAPIInteger:
		acceptedKinds = []reflect.Kind{reflect.Int}
	case openAPIBoolean:
		acceptedKinds = []reflect.Kind{reflect.Bool}
	case openAPIArray:
		acceptedKinds = []reflect.Kind{reflect.Slice}
	case openAPIObject:
		acceptedKinds = []reflect.Kind{reflect.Map}
	default:
		log.L.DebugWithData("default case should not run", log.Data{"name": v.Name, "type": *v.Type, "value": value})
		return fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.Type)
	}

	return isOneOF(actualKind, acceptedKinds)
}

func (v *Variable) prompt() string {
	sBuilder := strings.Builder{}
	qMark := common.ColoredFormat(common.ColorOrange, common.TextFormatBold, common.BGColorNone, "?")
	coloredName := common.ColoredFormat(common.ColorGreen, common.TextFormatBold, common.BGColorNone, v.Name)
	vType := func() string {
		if v.Type != nil {
			return string(*v.Type)
		}
		return "any"
	}()
	coloredType := common.ColoredFormat(common.ColorCyan, common.TextFormatBold, common.BGColorNone, vType)

	if v.Description != nil {
		coloredDescription := common.ColoredFormat(
			common.ColorGreen, common.TextFormatNormal, common.BGColorNone, *v.Description,
		)
		sBuilder.WriteString(fmt.Sprintf("%s [%s] ", coloredDescription, coloredType))
	} else {
		sBuilder.WriteString(fmt.Sprintf("[%s]", coloredType))
	}

	sBuilder.WriteString("\n")

	if v.Default != nil {
		defaultS, err := v.toString(v.Default)
		if err != nil {
			log.L.DebugWithData(
				`could not convert default value to string. 
This is probably due to a limitation of the current implementation.`,
				log.Data{"name": v.Name, "type": *v.Type, "value": v.Default},
			)
			defaultS = fmt.Sprintf("%v", v.Default)
		}
		coloredDefault := common.ColoredFormat(
			common.ColorOrange, common.TextFormatNormal, common.BGColorNone, defaultS,
		)
		sBuilder.WriteString(fmt.Sprintf("%s %s [%s]    ", qMark, coloredName, coloredDefault))
	} else {
		sBuilder.WriteString(fmt.Sprintf("%s %s    ", qMark, coloredName))
	}

	return sBuilder.String()
}

func (v *Variable) fromString(s string) (interface{}, error) {
	if v.Type == nil {
		return s, nil
	}
	var value interface{}
	var err error

	switch *v.Type {
	case openAPIString:
		value = s
	case openAPINumber:
		value, err = strconv.ParseFloat(s, 64)
	case openAPIInteger:
		value, err = strconv.Atoi(s)
	case openAPIBoolean:
		value, err = strconv.ParseBool(s)
	case openAPIArray:
		value = strings.Split(s, ",")
	case openAPIObject:
		valueMap := map[string]string{}
		for _, kvS := range strings.Split(s, ",") {
			kv := strings.SplitN(kvS, "=", 2)
			if len(kv) != 2 {
				err = fmt.Errorf("map error")
				break
			}
			valueMap[kv[0]] = kv[1]
		}
		value = valueMap
	default:
		log.L.DebugWithData("default case should not run", log.Data{"name": v.Name, "type": *v.Type, "string-value": s})
		return nil, fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.Type)
	}

	if err != nil {
		log.L.DebugWithData("external error", log.Data{"type": *v.Type, "string-value": s, "error": err.Error()})
		err = fmt.Errorf("variable value error: failed to parse string-value `%s` into open-api type `%s`", s, *v.Type)
	}
	return value, err
}

func (v *Variable) toString(value interface{}) (string, error) {
	if v.Type == nil {
		log.L.DebugWithData("toString on any defaults yo internal representation", log.Data{"value": value})
		return fmt.Sprintf("%v", value), nil
	}
	switch *v.Type {
	case openAPIString:
		fallthrough
	case openAPINumber:
		fallthrough
	case openAPIInteger:
		fallthrough
	case openAPIBoolean:
		return fmt.Sprintf("%v", value), nil
	case openAPIArray:
		vArray, ok := value.([]interface{})
		if !ok {
			log.L.DebugWithData("toString value type consistency error", log.Data{"type": *v.Type, "value": value})
			return "", fmt.Errorf("variable error: provided value inconsistent with variable type")
		}
		var items []string
		for _, item := range vArray {
			items = append(items, fmt.Sprintf("%v", item))
		}
		return strings.Join(items, ","), nil
	case openAPIObject:
		vMap, ok := value.(map[interface{}]interface{})
		if !ok {
			log.L.DebugWithData("toString value type consistency error", log.Data{"type": *v.Type, "value": value})
			return "", fmt.Errorf("variable error: provided value inconsistent with variable type")
		}
		var items []string
		for key, item := range vMap {
			items = append(items, fmt.Sprintf("%v=%v", key, item))
		}
		return strings.Join(items, ","), nil
	default:
		log.L.DebugWithData("default case should not run", log.Data{"name": v.Name})
		return "", fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.Type)
	}
}

func (v *Variable) Help() string {
	if v.Type == nil {
		return "input is not type, anything will do"
	}
	switch *v.Type {
	case openAPIString:
		return "input must be a string. example: `pizza`"
	case openAPINumber:
		return "input must be a number. example: `3.14`"
	case openAPIInteger:
		return "input must be an integer. example: `3`"
	case openAPIBoolean:
		return "input must be a boolean. example: `true`"
	case openAPIArray:
		return "input myst be an array of strings, example: `pizza,pasta,risotto`"
	case openAPIObject:
		return "input myst be an array of strings, example: `pizza=margherita,pasta=bolognese,risotto=mushroom`"
	default:
		log.L.DebugWithData("default case should not run", log.Data{"name": v.Name, "type": *v.Type})
		return ""
	}
}
