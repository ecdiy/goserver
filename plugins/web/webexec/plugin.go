package webexec

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins/web"
)

var fs = map[string]func(we *WebExec, ele *utils.Element, wb *utils.Param) error{
	"Http": Http, "Sp": Sp, "WebSocket": WebSocket, "Sql": Sql, "Param": Param,
}

func init() {
	web.RegisterWebPlugin("WebExec", func(ele *utils.Element) func(c *gin.Context) {
		we := &WebExec{Ele: ele, cache: make(map[string]interface{})}
		return we.run
	})
}

type WebExec struct {
	Ele *utils.Element

	cache map[string]interface{}
}

func (we *WebExec) run(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("WebExec fail;", err)
		}
	}()
	wb := utils.NewParam(ctx)
	ns := we.Ele.AllNodes()
	for _, n := range ns {
		f := fs[n.Name()]
		f(we, n, wb)
	}
	ctx.JSON(200, wb.Out)
}
