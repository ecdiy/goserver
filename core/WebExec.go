package core

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/ecdiy/goserver/webs"
	"reflect"
	"github.com/cihub/seelog"
	"strings"
)

func (app *Module) WebExec(ele *utils.Element) {
	we := &WebExec{ele: ele}
	we.webExec = reflect.ValueOf(we)
	method, mb := ele.AttrValue("Method")
	if mb && strings.ToLower(method) == "post" {
		getGin(ele).POST(ele.MustAttr("Url"), we.run)
	} else {
		getGin(ele).GET(ele.MustAttr("Url"), we.run)
	}
}

type WebExec struct {
	ele     *utils.Element
	webExec reflect.Value
}

func (we *WebExec) run(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("WebExec fail;", err)
		}
	}()
	wb := webs.NewParam(ctx)
	ns := we.ele.AllNodes()
	for _, n := range ns {
		inputs := make([]reflect.Value, 2)
		inputs[0] = reflect.ValueOf(n)
		inputs[1] = reflect.ValueOf(wb)
		m := we.webExec.MethodByName(n.Name())
		if m.IsValid() {
			v := m.Call(inputs)
			if len(v) == 1 && v[0].IsNil() {
				continue
			} else {
				wb.Out["Code"] = 1
				wb.Out["Msg"] = v[0].Interface().(error).Error()
				break
			}
		} else {
			panic("没有实现的方法:" + n.Name())
		}
	}
	ctx.JSON(200, wb.Out)
}

//----
func (we *WebExec) Http(ele *utils.Element, wb *webs.Param) error {
	hc := &HttpCore{}
	return hc.DoHttp(ele, wb)
}

func (we *WebExec) ExecSp(ele *utils.Element, wb *webs.Param) error {
	sp := &webs.WebSp{Gpa: getGpa(ele)}
	sp.Init()
	spName := ele.MustAttr("SpName")
	sp.SpExec(spName, wb)
	return nil
}
