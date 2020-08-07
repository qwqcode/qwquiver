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

// Sum 求和
func Sum(input ...float64) float64 {
	var sum float64 = 0
	for i := range input {
		sum += input[i]
	}
	return sum
}
