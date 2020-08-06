package utils

import (
	"reflect"
)

// GetStructFields 获取结构体中的所有字段名
func GetStructFields(t interface{}) []string {
	s := reflect.ValueOf(t).Elem()
	typeOfT := s.Type()

	names := []string{}
	for i := 0; i < s.NumField(); i++ {
		names = append(names, typeOfT.Field(i).Name)
	}

	return names
}
