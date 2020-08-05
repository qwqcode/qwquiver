package utils

import "encoding/json"

// JSONFormat JSON 字符串生成
func JSONFormat(obj interface{}) (str string, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}

// JSONParse JSON 解析
func JSONParse(str string, t interface{}) error {
	return json.Unmarshal([]byte(str), t)
}
