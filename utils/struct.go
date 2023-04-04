package utils

// ReturnCodeStruct 通用返回信息
type ReturnCodeStruct []struct {
	Msg  string                   `json:"msg"`  // 描述信息
	Code int                      `json:"code"` // 状态码
	Data []map[string]interface{} `json:"data"` // 返回数据
}
