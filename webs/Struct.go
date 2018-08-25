package webs

import "github.com/gin-gonic/gin"

type WebBase struct {
	Param, Out map[string]interface{} //参数，输出
	Ua         string
	Context    *gin.Context
}
type Sp struct {
	Sql               string
	Name, SessionName string
	Params            []SpParam
	Info              map[string]interface{}
}
type SpParam struct {
	ParamName, ParamType string
	Length               int
	ValFunc              ParamValFunc
	Param                []string
	//FuncType             string
}

type ParamValFunc func(ctx *gin.Context, param map[string]interface{}) (string, error)
