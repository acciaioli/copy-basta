package specification

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spin14/copy-basta/cmd/copy-basta/common/log"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

/*
Spec variables use a subset of the OpenApi data models

https://swagger.io/docs/specification/data-models
*/

type openAPiType string

const (
	openAPiString  = openAPiType("string")
	openAPiNumber  = openAPiType("number")
	openAPiInteger = openAPiType("integer")
	openAPiBoolean = openAPiType("boolean")
	openAPiArray   = openAPiType("array")
	openAPiObject  = openAPiType("object")
)

type SpecVariable struct {
	Name        string      `yaml:"name"`
	Type        openAPiType `yaml:"type"`
	Default     interface{} `yaml:"default"`
	Description *string     `yaml:"description"`
}

func (v *SpecVariable) validate() error {
	// name checks
	if v.Name == "" {
		return errors.New("variable error [name]: is required")
	}

	// type checks
	if v.Type == "" {
		return errors.New("variable error [type]: is required")
	}
	if ok := func(t openAPiType) bool {
		for _, t := range []openAPiType{openAPiString, openAPiNumber, openAPiInteger, openAPiBoolean, openAPiArray, openAPiObject} {
			if v.Type == t {
				return true
			}
		}
		return false
	}(v.Type); !ok {
		return fmt.Errorf(`variable error [type]: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, v.Type)
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
	actualKind := reflect.TypeOf(value).Kind()

	isOneOF := func(actual reflect.Kind, accepted []reflect.Kind) error {
		for _, acceptedKind := range accepted {
			if actual == acceptedKind {
				return nil
			}
		}
		return fmt.Errorf("value error: decoded to type %v, expected one of %v. variable type is %s", actual, accepted, string(v.Type))
	}

	var acceptedKinds []reflect.Kind

	switch v.Type {
	case openAPiString:
		acceptedKinds = []reflect.Kind{reflect.String}
	case openAPiNumber:
		acceptedKinds = []reflect.Kind{reflect.Int, reflect.Float64}
	case openAPiInteger:
		acceptedKinds = []reflect.Kind{reflect.Int}
	case openAPiBoolean:
		acceptedKinds = []reflect.Kind{reflect.Bool}
	case openAPiArray:
		acceptedKinds = []reflect.Kind{reflect.Slice}
	case openAPiObject:
		acceptedKinds = []reflect.Kind{reflect.Map}
	default:
		log.TheLogger.DebugWithData("default case should not run", log.LoggerData{"type": v.Type, "value": value})
		return fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, v.Type)
	}

	return isOneOF(actualKind, acceptedKinds)
}

func (v *SpecVariable) prompt() string {
	sBuilder := strings.Builder{}
	qMark := common.ColoredFormat(common.ColorOrange, common.TextFormatBold, common.BGColorNone, "?")
	coloredName := common.ColoredFormat(common.ColorGreen, common.TextFormatBold, common.BGColorNone, v.Name)
	coloredType := common.ColoredFormat(common.ColorCyan, common.TextFormatBold, common.BGColorNone, string(v.Type))

	if v.Description != nil {
		coloredDescription := common.ColoredFormat(common.ColorGreen, common.TextFormatNormal, common.BGColorNone, *v.Description)
		sBuilder.WriteString(fmt.Sprintf("%s [%s] ", coloredDescription, coloredType))
	} else {
		sBuilder.WriteString(fmt.Sprintf("[%s]", coloredType))
	}

	sBuilder.WriteString("\n")

	if v.Default != nil {
		coloredDefault := common.ColoredFormat(common.ColorOrange, common.TextFormatNormal, common.BGColorNone, fmt.Sprintf("%v", v.Default))
		sBuilder.WriteString(fmt.Sprintf("%s %s [%v]    ", qMark, coloredName, coloredDefault))
	} else {
		sBuilder.WriteString(fmt.Sprintf("%s %s    ", qMark, coloredName))
	}

	return sBuilder.String()
}

func (v *SpecVariable) process(s string) (interface{}, error) {
	var value interface{}
	var err error

	switch v.Type {
	case openAPiString:
		value = s
	case openAPiNumber:
		value, err = strconv.ParseFloat(s, 64)
	case openAPiInteger:
		value, err = strconv.Atoi(s)
	case openAPiBoolean:
		value, err = strconv.ParseBool(s)
	case openAPiArray:
		value = strings.Split(s, ",")
	case openAPiObject:
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
		log.TheLogger.DebugWithData("default case should not run", log.LoggerData{"type": v.Type, "value": value})
		return nil, fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, v.Type)
	}

	if err != nil {
		log.TheLogger.DebugWithData("external error", log.LoggerData{"type": v.Type, "string-value": value, "error": err.Error()})
		err = fmt.Errorf("variable value error: failed to parse string-value %s, open-api type is %s", value, v.Type)
	}
	return value, err
}
