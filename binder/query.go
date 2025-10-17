package binder

import (
	"fmt"
	"reflect"
)

func (d *RequestDecoder) BindQuery(v interface{}) error {
	params := d.Request.URL.Query()

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
		queryValue := params.Get(fieldName)
		if queryValue == "" {
			continue
		}

		if err := setFieldValue(fieldVal, queryValue, structTagDefaultValueDelimiter); err != nil {
			return fmt.Errorf("error setting field %s: %v", field.Name, err)
		}
	}

	return nil
}

func (d *RequestEncoder) BindQuery(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("out must be a pointer to a struct")
	}
	val = val.Elem()

	query := d.Request.URL.Query()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldVal := val.Field(i)
		tag := field.Tag.Get(structTag)
		if tag == "" || !fieldVal.CanSet() {
			continue
		}

		fieldName := parseTag(tag)
		query.Add(fieldName, fmt.Sprintf("%v", fieldVal.Interface()))
	}

	d.Request.URL.RawQuery = query.Encode()
	return nil
}
