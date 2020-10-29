package specification

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"copy-basta/internal/common/log"

	"copy-basta/internal/common"

	"gopkg.in/yaml.v2"
)

/*
Spec variables use a subset of the OpenApi data models

https://swagger.io/docs/specification/data-models
*/

const (
	openAPIString  = "string"
	openAPINumber  = "number"
	openAPIInteger = "integer"
	openAPIBoolean = "boolean"
	openAPIArray   = "array"
	openAPIObject  = "object"
)

var openAPITypes = []string{
	openAPIString,
	openAPINumber,
	openAPIInteger,
	openAPIBoolean,
	openAPIArray,
	openAPIObject,
}

type Variables []Variable

type Variable struct {
	name        string
	dtype       *string
	defaultVal  interface{}
	description *string
}

func NewVariables(varData []VariableData) (Variables, error) {
	vars := Variables{}
	for _, vd := range varData {
		v := Variable{
			name:        vd.Name,
			dtype:       vd.DType,
			defaultVal:  vd.DefaultVal,
			description: vd.Description,
		}
		if err := v.validate(); err != nil {
			return nil, err
		}
		vars = append(vars, v)
	}
	return vars, nil
}

func (vars Variables) InputFromFile(inputYAML string) (common.InputVariables, error) {
	yamlFile, err := ioutil.ReadFile(inputYAML)
	if err != nil {
		return nil, err
	}

	input := common.InputVariables{}
	err = yaml.Unmarshal(yamlFile, &input)
	if err != nil {
		return nil, err
	}

	for _, v := range vars {
		value, ok := input[v.name]
		if !ok {
			if v.defaultVal != nil {
				return nil, fmt.Errorf("no value nor default for %s", v.name)
			}
			input[v.name] = v.defaultVal
		}
		if err := v.valueOk(value); err != nil {
			return nil, err
		}
	}

	return input, nil
}

func (vars Variables) InputFromStdIn() (common.InputVariables, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("\n")
	inputVars := common.InputVariables{}
	for _, v := range vars {
		for retry := 3; retry > 0; retry-- {
			userInput, err := v.promptLoop(r)
			if err != nil {
				return nil, err
			}

			if userInput != nil {
				value, err := v.fromString(*userInput)
				if err != nil {
					if retry > 1 {
						fmt.Println(v.Help())
						continue
					}
					return nil, err
				}
				inputVars[v.name] = value
			} else {
				inputVars[v.name] = v.defaultVal
			}
			break
		}

	}
	return inputVars, nil
}

func (v *Variable) validate() error {
	// name checks
	if v.name == "" {
		return errors.New("variable error [name]: is required")
	}

	// type checks
	if v.dtype != nil {
		if ok := func(actualType string) bool {
			for _, candidateType := range openAPITypes {
				if actualType == candidateType {
					return true
				}
			}
			return false
		}(*v.dtype); !ok {
			return fmt.Errorf(`variable error [type]: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.dtype)
		}
	} else {
		log.L.WarnWithData("spec variable without type, defaulting to any", log.Data{"name": v.name})
	}

	// default checks
	if v.defaultVal != nil {
		if err := v.valueOk(v.defaultVal); err != nil {
			return fmt.Errorf("variable error [default]: %s", err.Error())
		}
	}

	return nil
}

func (v *Variable) valueOk(value interface{}) error {
	if v.dtype == nil {
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
		return fmt.Errorf(format, actual, accepted, *v.dtype)
	}

	var acceptedKinds []reflect.Kind

	switch *v.dtype {
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
		log.L.DebugWithData("default case should not run", log.Data{"name": v.name, "type": *v.dtype, "value": value})
		return fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.dtype)
	}

	return isOneOF(actualKind, acceptedKinds)
}

func (v *Variable) promptLoop(r *bufio.Reader) (*string, error) {
	for {
		fmt.Print(v.prompt())
		userInput, err := r.ReadString('\n')
		fmt.Print("\n")
		if err != nil {
			return nil, err
		}
		userInput = strings.TrimSuffix(userInput, "\n")

		if userInput != "" {
			return &userInput, nil
		}

		if v.defaultVal != nil {
			return nil, nil
		}
	}
}

func (v *Variable) prompt() string {
	sBuilder := strings.Builder{}
	qMark := common.ColoredFormat(common.ColorOrange, common.TextFormatBold, common.BGColorNone, "?")
	coloredName := common.ColoredFormat(common.ColorGreen, common.TextFormatBold, common.BGColorNone, v.name)
	vType := func() string {
		if v.dtype != nil {
			return *v.dtype
		}
		return "any"
	}()
	coloredType := common.ColoredFormat(common.ColorCyan, common.TextFormatBold, common.BGColorNone, vType)

	if v.description != nil {
		coloredDescription := common.ColoredFormat(
			common.ColorGreen, common.TextFormatNormal, common.BGColorNone, *v.description,
		)
		sBuilder.WriteString(fmt.Sprintf("%s [%s] ", coloredDescription, coloredType))
	} else {
		sBuilder.WriteString(fmt.Sprintf("[%s]", coloredType))
	}

	sBuilder.WriteString("\n")

	if v.defaultVal != nil {
		defaultS, err := v.toString(v.defaultVal)
		if err != nil {
			log.L.DebugWithData(
				`could not convert default value to string. 
This is probably due to a limitation of the current implementation.`,
				log.Data{"name": v.name, "type": *v.dtype, "value": v.defaultVal},
			)
			defaultS = fmt.Sprintf("%v", v.defaultVal)
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
	if v.dtype == nil {
		return s, nil
	}
	var value interface{}
	var err error

	switch *v.dtype {
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
		log.L.DebugWithData("default case should not run", log.Data{"name": v.name, "type": *v.dtype, "string-value": s})
		return nil, fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.dtype)
	}

	if err != nil {
		log.L.DebugWithData("external error", log.Data{"type": *v.dtype, "string-value": s, "error": err.Error()})
		err = fmt.Errorf("variable value error: failed to parse string-value `%s` into open-api type `%s`", s, *v.dtype)
	}
	return value, err
}

func (v *Variable) toString(value interface{}) (string, error) {
	if v.dtype == nil {
		log.L.DebugWithData("toString on any defaults yo internal representation", log.Data{"value": value})
		return fmt.Sprintf("%v", value), nil
	}
	switch *v.dtype {
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
			log.L.DebugWithData("toString value type consistency error", log.Data{"type": *v.dtype, "value": value})
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
			log.L.DebugWithData("toString value type consistency error", log.Data{"type": *v.dtype, "value": value})
			return "", fmt.Errorf("variable error: provided value inconsistent with variable type")
		}
		var items []string
		for key, item := range vMap {
			items = append(items, fmt.Sprintf("%v=%v", key, item))
		}
		return strings.Join(items, ","), nil
	default:
		log.L.DebugWithData("default case should not run", log.Data{"name": v.name})
		return "", fmt.Errorf(`variable error: %s is not a valid type. 
only open-api types are supported (https://swagger.io/docs/specification/data-models/data-types)`, *v.dtype)
	}
}

func (v *Variable) Help() string {
	if v.dtype == nil {
		return "input is not type, anything will do"
	}
	switch *v.dtype {
	case openAPIString:
		return "input must be a string. example: `pizza`"
	case openAPINumber:
		return "input must be a number. example: `3.14`"
	case openAPIInteger:
		return "input must be an integer. example: `3`"
	case openAPIBoolean:
		return "input must be a boolean. example: `true`"
	case openAPIArray:
		return "input must be an array of strings, example: `pizza,pasta,risotto`"
	case openAPIObject:
		return "input must be an string to string map , example: `pizza=margherita,pasta=bolognese,risotto=mushroom`"
	default:
		log.L.DebugWithData("default case should not run", log.Data{"name": v.name, "type": *v.dtype})
		return ""
	}
}
