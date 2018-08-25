package webs

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
	"utils"
)



func (p *WebBase) String(n string) string {
	v, vb := p.Param[n]
	if vb {
		return fmt.Sprint(v)
	} else {
		v2, _ := p.Context.GetQuery(n)
		return v2
	}
}
func (p *WebBase) Int64(n string) int64 {
	return p.Int64Default(n, 0)
}
func (p *WebBase) Int64Default(n string, def int64) int64 {
	i, e := strconv.ParseInt(p.String(n), 10, 0)
	if e != nil {
		return def
	}
	return i
}
func (p *WebBase) Start() int64 {
	return p.StartPageSize(DefaultPageSize)
}
func (p *WebBase) StartPageSize(ps int64) int64 {
	page := p.Int64Default("page", 1)
	if page < 1 {
		return int64(0)
	}
	return (page - 1) * ps
}
func (p *WebBase) Result(result ...interface{}) {
	if result != nil {
		p.Out["result"] = result
	}
}
func (p *WebBase) ST(st *utils.ST, result ...interface{}) {
	p.Out["code"] = st.Code
	p.Out["msg"] = st.Msg
	p.Result(result...)
}

func WebBaseNew(c *gin.Context) *WebBase {
	web := &WebBase{}
	web.Ua = GetUa(c)
	web.Context = c
	web.Out = make(map[string]interface{})
	return web
}
