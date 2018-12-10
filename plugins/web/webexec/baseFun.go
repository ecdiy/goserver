package webexec

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins/http"
	"github.com/ecdiy/goserver/plugins/sp"
	"github.com/ecdiy/goserver/plugins"
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/gpa"
	"fmt"
	"github.com/ecdiy/goserver/plugins/web/ws"
	"encoding/json"
)

//----
func Http(we *WebExec, ele *utils.Element, wb *utils.Param) error {
	hc := &http.HCore{}
	return hc.DoHttp(ele, wb)
}

func Sp(we *WebExec, ele *utils.Element, wb *utils.Param) error {
	spName := ele.MustAttr("SpName")
	webSp, spExt := we.cache[spName]
	var spn *sp.WebSp
	if !spExt {
		spn = &sp.WebSp{Gpa: plugins.GetRef(ele, "Gpa").(*gpa.Gpa)}
		spn.Init(ele)
		we.cache[spName] = spn
	} else {
		spn = webSp.(*sp.WebSp)
	}
	code := spn.SpExec(spName, wb)
	if code != 200 {
		seelog.Error("", spName)
	}
	return nil
}

func Param(we *WebExec, ele *utils.Element, wb *utils.Param) error {
	for _, attr := range ele.Attrs {
		wb.Param[attr.Name()] = attr.Value
	}
	return nil
}

func Sql(we *WebExec, ele *utils.Element, wb *utils.Param) error {
	dao := plugins.GetRef(ele, "Gpa").(*gpa.Gpa)
	dao.Exec(ele.Value)
	return nil
}

func WebSocket(we *WebExec, ele *utils.Element, wb *utils.Param) error {
	UserIdName := ele.Attr("SocketIdName", "UserId")
	userIdFc, idExt := wb.Out[UserIdName]
	if idExt {
		userId := fmt.Sprint(userIdFc)
		bs, err := json.Marshal(wb.Out)
		if err == nil {
			ws.WrMsg(userId, bs)
		}
	} else {
		seelog.Error("session value not find key=", UserIdName)
	}
	return nil
}
