package binder

import (
	"fmt"
	"reflect"
)

func (d *RequestDecoder) BindParams(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("out must be a pointer to a struct")
	}
	val = val.Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldVal := val.Field(i)
		tag := field.Tag.Get(structTag)
		if tag == "" || !fieldVal.CanSet() {
			continue
		}

		fieldName := parseTag(tag)
		paramValue := d.Request.PathValue(fieldName)
		if paramValue == "" {
			continue
		}

		if err := setFieldValue(fieldVal, paramValue, structTagDefaultValueDelimiter); err != nil {
			return fmt.Errorf("error setting field %s: %v", field.Name, err)
		}
	}

	return nil
}

func (d *RequestEncoder) BindParams(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("out must be a pointer to a struct")
	}
	val = val.Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldVal := val.Field(i)
		tag := field.Tag.Get(structTag)
		if tag == "" || !fieldVal.CanSet() {
			continue
		}

		fieldName := parseTag(tag)

		val := fmt.Sprintf("%v", fieldVal.Interface())
		d.Request.SetPathValue(fieldName, val)
	}

	return nil
}
