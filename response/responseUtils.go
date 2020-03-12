package responseUitls

import "github.com/alphayan/iris"

//错误返回信息类old
type ResMsg struct {
	ErrCode int64
	ErrMsg  string
	Data    interface{}
}

func NewResMsg(resData interface{}, err error) ResponseEntity {
	var res ResponseEntity
	if err != nil {
		res.ErrCode = iris.StatusBadRequest
		res.ErrMsg = err.Error()
		return res
	}
	res.ErrCode = iris.StatusOK
	res.ErrMsg = "请求成功"
	res.Data = resData
	return res
}

//错误返回信息类
type ResponseEntity struct {
	ErrCode int64
	ErrMsg  string
	Data    interface{}
}

func NewResponseEntity(resData interface{}, err error) ResponseEntity {
	var res ResponseEntity
	if err != nil {
		res.ErrCode = iris.StatusBadRequest
		res.ErrMsg = err.Error()
		return res
	}
	res.ErrCode = iris.StatusOK
	res.ErrMsg = "请求成功"
	res.Data = resData
	return res
}
