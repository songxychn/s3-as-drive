package types

type ResultEnum int

const (
	SUCCESS ResultEnum = 2000
	ERROR   ResultEnum = 5000
)

type Result struct {
	Code ResultEnum  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) Result {
	return Result{Code: SUCCESS, Msg: "成功", Data: data}
}

func SuccessEmpty() Result {
	return Result{Code: SUCCESS, Msg: "成功"}
}

func Error(msg string) Result {
	return Result{Code: ERROR, Msg: msg}
}
