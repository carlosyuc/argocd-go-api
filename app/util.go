package main

// import (
// 	"os"
// 	"reflect"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// const (
// 	Tag        = "env"
// 	AliasTag   = "env-alias"
// 	DefaultTag = "env-default"
// 	Separator  = ","
// )

// type (
// 	lookupFunc = func(string) (string, bool)

// 	environmentVariable struct {
// 		Field        reflect.Value
// 		Name         string
// 		Alias        string
// 		HasAlias     bool
// 		Required     bool
// 		DefaultValue string
// 	}
// )

// func Load(config interface{}) any {
// 	lookup := os.LookupEnv
// 	return load(config, lookup)
// }

// func load(config interface{}, lookup lookupFunc) any {
// 	variables, err := parseFields(config)
// 	if err != nil {
// 		return err
// 	}

// 	for _, variable := range variables {
// 		value, set := lookup(variable.Name)
// 		if !set && variable.HasAlias {
// 			value, set = lookup(variable.Alias)
// 		}
// 		if !set {
// 			// 	if variable.Required {
// 			// 		return newRequiredEnvironmentVariableNotSetError(variable.Name)
// 			// 	}
// 			value = variable.DefaultValue
// 		}

// 		if err = setValue(variable.Field, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func setValue(field reflect.Value, value string) {
// 	const timePath = "time"
// 	switch field.Type().Kind() {
// 	case reflect.String:
// 		field.SetString(value)

// 	case reflect.Bool:
// 		return setBoolean(field, value)

// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		if field.Kind() == reflect.Int64 && field.Type().PkgPath() == timePath && field.Type().Name() == "Duration" {
// 			return setDuration(field, value)
// 		}
// 		return setInteger(field, value)

// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 		return setUnsignedInteger(field, value)

// 	case reflect.Float32, reflect.Float64:
// 		return setFloat(field, value)

// 	case reflect.Slice:
// 		return setSlice(field, value)

// 	case reflect.Map:
// 		return setMap(field, value)

// 	case reflect.Struct:
// 		if field.Type().PkgPath() == timePath && field.Type().Name() == "Time" {
// 			return setTime(field, value)
// 		}

// 	default:
// 		return newUnsupportedFieldKindError(field.Type().Kind())
// 	}
// 	return nil
// }

// func parseFields(config interface{}) ([]environmentVariable, any) {
// 	variables := make([]environmentVariable, 0)

// 	layers := []interface{}{config}
// 	for layerIndex := 0; layerIndex < len(layers); layerIndex++ {
// 		layer := reflect.ValueOf(layers[layerIndex])
// 		if layer.Kind() == reflect.Ptr {
// 			layer = layer.Elem()
// 		}
// 		// if layer.Kind() != reflect.Struct {
// 		// 	return nil, newConfigTypeIsNotAStructError()
// 		// }

// 		for fieldIndex := 0; fieldIndex < layer.NumField(); fieldIndex++ {
// 			field := layer.Field(fieldIndex)
// 			if field.Kind() == reflect.Struct && field.Type() != reflect.TypeOf(time.Time{}) {
// 				layers = append(layers, field.Addr().Interface())
// 				continue
// 			}

// 			if !field.CanSet() {
// 				continue
// 			}

// 			tag := layer.Type().Field(fieldIndex).Tag
// 			envVariable, hasEnvTag := tag.Lookup(Tag)
// 			if !hasEnvTag {
// 				continue
// 			}
// 			aliasVariable, hasAliasTag := tag.Lookup(AliasTag)
// 			defaultValue, hasDefaultTag := tag.Lookup(DefaultTag)

// 			variable := environmentVariable{
// 				Field:        field,
// 				Name:         envVariable,
// 				Alias:        aliasVariable,
// 				HasAlias:     hasAliasTag,
// 				Required:     !hasDefaultTag,
// 				DefaultValue: defaultValue,
// 			}

// 			variables = append(variables, variable)
// 		}
// 	}
// 	return variables, nil
// }

// func setBoolean(field reflect.Value, value string) (err errors.Error) {
// 	boolean, parseErr := strconv.ParseBool(value)
// 	if parseErr != nil {
// 		// return newMalformedConfigTagError("boolean", value)
// 	}
// 	field.SetBool(boolean)
// 	return
// }

// func setDuration(field reflect.Value, value string) (err errors.Error) {
// 	duration, parseErr := time.ParseDuration(value)
// 	if parseErr != nil {
// 		// return newMalformedConfigTagError("time.Duration", value)
// 	}
// 	field.SetInt(int64(duration))
// 	return
// }

// func setInteger(field reflect.Value, value string) (err errors.Error) {
// 	number, parseErr := strconv.ParseInt(value, 0, field.Type().Bits())
// 	if parseErr != nil {
// 		// return newMalformedConfigTagError("integer", value)
// 	}
// 	field.SetInt(number)
// 	return
// }

// func setUnsignedInteger(field reflect.Value, value string) (err errors.Error) {
// 	number, parseErr := strconv.ParseUint(value, 0, field.Type().Bits())
// 	if parseErr != nil {
// 		// return newMalformedConfigTagError("unsigned integer", value)
// 	}
// 	field.SetUint(number)
// 	return
// }

// func setFloat(field reflect.Value, value string) (err errors.Error) {
// 	number, parseErr := strconv.ParseFloat(value, field.Type().Bits())
// 	if parseErr != nil {
// 		// return newMalformedConfigTagError("float", value)
// 	}
// 	field.SetFloat(number)
// 	return
// }

// func setSlice(field reflect.Value, value string) (err errors.Error) {
// 	slice := reflect.MakeSlice(field.Type(), 0, 0)

// 	if field.Type().Elem().Kind() == reflect.Uint8 {
// 		slice = reflect.ValueOf([]byte(value))
// 	} else if len(strings.TrimSpace(value)) != 0 {
// 		elements := strings.Split(value, Separator)
// 		slice = reflect.MakeSlice(field.Type(), len(elements), len(elements))
// 		for i, element := range elements {
// 			if err = setValue(slice.Index(i), element); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	field.Set(slice)
// 	return
// }

// func setMap(field reflect.Value, value string) (err errors.Error) {
// 	const ExpectedLength = 2
// 	mapValue := reflect.MakeMap(field.Type())
// 	if len(strings.TrimSpace(value)) != 0 {
// 		entries := strings.Split(value, Separator)
// 		for _, entry := range entries {
// 			const substringCount = 2
// 			keyValuePair := strings.SplitN(entry, ":", substringCount)
// 			if len(keyValuePair) != ExpectedLength {
// 				return newMalformedConfigTagError("map entry", keyValuePair)
// 			}
// 			key := reflect.New(field.Type().Key()).Elem()
// 			err = setValue(key, keyValuePair[0])
// 			if err != nil {
// 				return
// 			}
// 			val := reflect.New(field.Type().Elem()).Elem()
// 			err = setValue(val, keyValuePair[1])
// 			if err != nil {
// 				return
// 			}
// 			mapValue.SetMapIndex(key, val)
// 		}
// 	}
// 	field.Set(mapValue)
// 	return
// }

// func setTime(field reflect.Value, value string) (err errors.Error) {
// 	parsedTime, parseErr := time.Parse(time.RFC3339, value)
// 	if parseErr != nil {
// 		return newMalformedConfigTagError("time.Time", value)
// 	}
// 	field.Set(reflect.ValueOf(parsedTime))
// 	return
// }
