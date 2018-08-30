package webs

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strings"
)

var router = make(map[string]*MultiMerge)

type MultiMerge struct {
	//Fun, Verify func(wdb *Web)
	ReqType     string
}

/**
多个请求合并
*/
func WebNoRouterMultiRequestMerge(ctx *gin.Context) {
	url := ctx.Request.URL.Path
	if strings.Index(url, ",") >= 0 {
		pp := strings.Split(url, ",")
		var os []*MultiMerge
		for i := 1; i < len(pp); i++ {
			obj, f := router[pp[i]]
			if !f {
				ctx.AbortWithStatus(404)
				seelog.Error("not find url:", pp[i])
				return
			} else {
				os = append(os, obj)
			}
		}
		//web := &Web{}
		//web.WebBase = webs.WebBaseNew(ctx)
		//web.initParam()
		//for _, mm := range os {
		//	if mm.Verify == nil {
		//		mm.Fun(web)
		//	} else {
		//		if !web.Auth {
		//			mm.Verify(web)
		//		}
		//		if web.Auth {
		//			mm.Fun(web)
		//		} else {
		//			ctx.AbortWithStatus(401)
		//			return
		//		}
		//	}
		//}
	} else {
		seelog.Warn("no router.", url)
	}
}
