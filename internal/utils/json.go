package utils

import "encoding/json"

func Json(obj interface{}, indent ...bool) string {
	if len(indent) > 0 {
		bytes, err := json.MarshalIndent(obj, "", "  ") // indent with two spaces
		if err != nil {
			return ""
		}
		return string(bytes)
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bytes)
}
