package webs

import "github.com/gin-gonic/gin"

type WebBase struct {
	Param, Out map[string]interface{} //参数，输出
	Ua         string
	Context    *gin.Context
}
