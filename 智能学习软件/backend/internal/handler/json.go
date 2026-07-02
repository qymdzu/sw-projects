package handler

import "encoding/json"

// jsonUnmarshal 是对 json.Unmarshal 的薄封装，便于在 handler 包中统一替换或 mock。
func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// jsonMarshal 是对 json.Marshal 的薄封装。
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}