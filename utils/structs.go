package utils

import (
	"errors"
	"reflect"
)

var (
	errTagIsEmpty = errors.New("The field doesn't have Tag")
)

func ParseStructToMap(m interface{}) (generatedMap map[string]interface{}, err error) {
	var (
		value         = reflect.ValueOf(m)
		typeOf        = reflect.TypeOf(m)
		tagName       = "collumn"
		tagIgnoreName = "ignoreInsertUpdate"
	)

	generatedMap = make(map[string]interface{})

	for i := 0; i < typeOf.Elem().NumField(); i++ {
		fieldName := typeOf.Elem().Field(i).Tag.Get(tagName)
		ignore := typeOf.Elem().Field(i).Tag.Get(tagIgnoreName)

		if fieldName == "" {
			err = errTagIsEmpty
			return
		} else if fieldName == "-" || ignore != "" {
			continue
		}

		generatedMap[fieldName] = value.Elem().Field(i).Interface()
	}

	return
}
