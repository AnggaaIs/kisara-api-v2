package utils

import "reflect"

// CreateNew creates a new instance of the same type as the input interface
func CreateNew(i interface{}) interface{} {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
