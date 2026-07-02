package service_test

import "encoding/json"

// jsonStdUnmarshal 是对 encoding/json.Unmarshal 的薄封装，便于测试中调用。
func jsonStdUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}