package specification

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"copy-basta/cmd/common/log"

	"copy-basta/cmd/common"
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

type SpecVariable struct {
	Name        string       `yaml:"name"`
	Type        *openAPIType `yaml:"type"`
	Default     interface{}  `yaml:"default"`
	Description *string      `yaml:"description"`
}

func (v *SpecVariable) validate() error {
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

func (v *SpecVariable) valueOk(value interface{}) error {
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

func (v *SpecVariable) prompt() string {
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
		coloredDefault := common.ColoredFormat(
			common.ColorOrange, common.TextFormatNormal, common.BGColorNone, fmt.Sprintf("%v", v.Default),
		)
		sBuilder.WriteString(fmt.Sprintf("%s %s [%v]    ", qMark, coloredName, coloredDefault))
	} else {
		sBuilder.WriteString(fmt.Sprintf("%s %s    ", qMark, coloredName))
	}

	return sBuilder.String()
}

func (v *SpecVariable) process(s string) (interface{}, error) {
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
		log.L.DebugWithData("default case should not run", log.Data{"name": v.Name, "type": *v.Type, "value": value})
		return nil, fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.Type)
	}

	if err != nil {
		log.L.DebugWithData("external error", log.Data{"type": *v.Type, "string-value": value, "error": err.Error()})
		err = fmt.Errorf("variable value error: failed to parse string-value %s, open-api type is %s", value, *v.Type)
	}
	return value, err
}
