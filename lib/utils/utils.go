package utils

import (
	"encoding/json"
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

func JSONEncode(obj interface{}) (str string, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}

func JSONDecode(str string, t interface{}) error {
	return json.Unmarshal([]byte(str), t)
}

// ContainsInt returns true if an int is present in a iteratee.
func ContainsInt(s []int, v int) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsInt32 returns true if an int32 is present in a iteratee.
func ContainsInt32(s []int32, v int32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsInt64 returns true if an int64 is present in a iteratee.
func ContainsInt64(s []int64, v int64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsUInt returns true if an uint is present in a iteratee.
func ContainsUInt(s []uint, v uint) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsUInt32 returns true if an uint32 is present in a iteratee.
func ContainsUInt32(s []uint32, v uint32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsUInt64 returns true if an uint64 is present in a iteratee.
func ContainsUInt64(s []uint64, v uint64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsString returns true if a string is present in a iteratee.
func ContainsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsFloat32 returns true if a float32 is present in a iteratee.
func ContainsFloat32(s []float32, v float32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainsFloat64 returns true if a float64 is present in a iteratee.
func ContainsFloat64(s []float64, v float64) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

// SumInt32 sums a int32 iteratee and returns the sum of all elements
func SumInt32(s []int32) (sum int32) {
	for _, v := range s {
		sum += v
	}
	return
}

// SumInt64 sums a int64 iteratee and returns the sum of all elements
func SumInt64(s []int64) (sum int64) {
	for _, v := range s {
		sum += v
	}
	return
}
