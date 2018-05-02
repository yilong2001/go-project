package objutils

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func GetFieldsWithNotDefaultValue(obj interface{}, intdefault string, strdefault string, skips []string) (map[string]interface{}, map[string]interface{}) {
	object := reflect.ValueOf(obj)
	myref := object.Elem()
	typeOfType := myref.Type()

	fieldIfArrs := make(map[string]interface{})
	fieldAddrIfArrs := make(map[string]interface{})

	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)

		isskip := false
		for _, skip := range skips {
			if typeOfType.Field(i).Name == skip {
				isskip = true
				break
			}
		}

		if isskip {
			continue
		}

		found := false
		str := fmt.Sprintf("%v", field.Interface())

		if field.Type().Name() == "int" && str != intdefault {
			found = true
		}

		if field.Type().Name() == "int8" && str != intdefault {
			found = true
		}

		if field.Type().Name() == "int16" && str != intdefault {
			found = true
		}

		if field.Type().Name() == "int32" && str != intdefault {
			found = true
		}

		if field.Type().Name() == "int64" && str != intdefault {
			found = true
		}

		if field.Type().Name() == "string" && str != strdefault {
			found = true
		}

		if found {
			nm, err := CamelToUnderLine(typeOfType.Field(i).Name)
			if err == nil {
				fieldIfArrs[nm] = field.Interface()
				fieldAddrIfArrs[nm] = field.Addr().Interface()
			}
		}
	}

	log.Println(fieldIfArrs)
	log.Println(fieldAddrIfArrs)

	return fieldIfArrs, fieldAddrIfArrs
}

func GetWholeFields(obj interface{}) (map[string]interface{}, map[string]interface{}) {
	object := reflect.ValueOf(obj)
	myref := object.Elem()
	typeOfType := myref.Type()

	fieldIfArrs := make(map[string]interface{})
	fieldAddrIfArrs := make(map[string]interface{})

	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)

		// isskip := false
		// for _, skip := range skips {
		// 	if strings.ToLower(typeOfType.Field(i).Name) == strings.ToLower(skip) {
		// 		isskip = true
		// 		break
		// 	}
		// }

		// if isskip {
		// 	continue
		// }

		found := false

		if field.Type().Name() == "int" {
			found = true
		}

		if field.Type().Name() == "int8" {
			found = true
		}

		if field.Type().Name() == "int64" {
			found = true
		}

		if field.Type().Name() == "float32" {
			found = true
		}

		if field.Type().Name() == "float64" {
			found = true
		}

		if field.Type().Name() == "string" {
			found = true
		}

		if found {
			nm, err := CamelToUnderLine(typeOfType.Field(i).Name)
			if err == nil {
				fieldIfArrs[nm] = field.Interface()
				//log.Println(typeOfType.Field(i).Name, field.Interface())
				fieldAddrIfArrs[nm] = field.Addr().Interface()
			}
		}
	}

	return fieldIfArrs, fieldAddrIfArrs
}

func GetFieldsWithSkip(obj interface{}, skips []string) (map[string]interface{}, map[string]interface{}) {
	object := reflect.ValueOf(obj)
	myref := object.Elem()
	typeOfType := myref.Type()

	fieldIfArrs := make(map[string]interface{})
	fieldAddrIfArrs := make(map[string]interface{})

	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)

		isskip := false
		for _, skip := range skips {
			if strings.ToLower(typeOfType.Field(i).Name) == strings.ToLower(skip) {
				isskip = true
				break
			}
		}

		if isskip {
			continue
		}

		found := false

		if field.Type().Name() == "int" {
			found = true
		}

		if field.Type().Name() == "int8" {
			found = true
		}

		if field.Type().Name() == "int64" {
			found = true
		}

		if field.Type().Name() == "float32" {
			found = true
		}

		if field.Type().Name() == "float64" {
			found = true
		}

		if field.Type().Name() == "string" {
			found = true
		}

		if found {
			nm, err := CamelToUnderLine(typeOfType.Field(i).Name)
			if err == nil {
				fieldIfArrs[nm] = field.Interface()
				fieldAddrIfArrs[nm] = field.Addr().Interface()
			}
		}
	}

	return fieldIfArrs, fieldAddrIfArrs
}

func GetFieldsWithSpecs(obj interface{}, specs []string) (map[string]interface{}, map[string]interface{}) {
	object := reflect.ValueOf(obj)
	myref := object.Elem()
	typeOfType := myref.Type()

	fieldIfArrs := make(map[string]interface{})
	fieldAddrIfArrs := make(map[string]interface{})

	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)

		isfind := false
		for _, spec := range specs {
			if strings.ToLower(typeOfType.Field(i).Name) == strings.ToLower(spec) {
				isfind = true
				break
			}
		}

		if !isfind {
			continue
		}

		found := false

		if field.Type().Name() == "int" {
			found = true
		}

		if field.Type().Name() == "int8" {
			found = true
		}

		if field.Type().Name() == "int64" {
			found = true
		}

		if field.Type().Name() == "float32" {
			found = true
		}

		if field.Type().Name() == "float64" {
			found = true
		}

		if field.Type().Name() == "string" {
			found = true
		}

		if found {
			nm, err := CamelToUnderLine(typeOfType.Field(i).Name)
			if err == nil {
				fieldIfArrs[nm] = field.Interface()
				fieldAddrIfArrs[nm] = field.Addr().Interface()
			}
		}
	}

	return fieldIfArrs, fieldAddrIfArrs
}
