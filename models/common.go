package models

import "gisa/backend/logic"

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code logic.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Obj  interface{}          `json:"obj"`
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
