package webexec

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"reflect"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/plugins/web"
)

func init() {
	web.RegisterWebPlugin("WebExec", func(ele *utils.Element) func(c *gin.Context) {
		we := &WebExec{Ele: ele}
		we.WebExec = reflect.ValueOf(we)
		return we.run
	})
}

type WebExec struct {
	Ele     *utils.Element
	WebExec reflect.Value
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
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(n)
		inputs[1] = reflect.ValueOf(wb)
		m := we.WebExec.MethodByName(n.Name())
		if m.IsValid() {
			v := m.Call(inputs)
			if len(v) == 1 && v[0].IsNil() {
				continue
			} else {
				wb.Out["Code"] = 1
				if len(v) >= 1 {
					wb.Out["Msg"] = v[0].Interface().(error).Error()
				}
				seelog.Error("error on.", n.Name())
				return
			}
		} else {
			seelog.Error("没有实现的方法:" + n.Name())
			return
		}
	}
	ctx.JSON(200, wb.Out)
}
