package handlers

import (
	"github.com/alphayan/iris"
)

/**
 * 返回消息
 */
type Msg struct {
	ErrMsg string `json:"errMsg"`
}

/**
 * 返回json数据格式
 */
func AfterHandlersJSON(ctx iris.Context, resData interface{}, err error) {
	res := make(map[string]interface{})
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		res := Msg{ErrMsg: err.Error()}
		ctx.JSON(res)
	} else {
		res["errCode"] = 0
		if resData != nil {
			res["data"] = resData
		}
		ctx.JSON(res)
	}
}

/**
 * 返回jsonp格式数据
 */
func AfterHandlersJSONP(ctx iris.Context, resData interface{}, err error) {
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		resData = Msg{ErrMsg: err.Error()}
	}
	ctx.JSONP(resData)
}
