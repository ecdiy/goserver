package webs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
	"io/ioutil"
	"github.com/cihub/seelog"
	"strings"
	"goserver/utils"
)

type Param struct {
	//Auth,
	Param, Out      map[string]interface{}
	Ua, ContentType string
	Context         *gin.Context
}

//func (p *Param) Username() string {
//	return fmt.Sprint(p.Auth["Username"])
//}

func (p *Param) Print() {
	for k, v := range p.Context.Request.PostForm {
		seelog.Info("post from:", k, v)
	}
	for k, v := range p.Context.Request.Form {
		seelog.Info(" from:", k, v)
	}

	data, _ := ioutil.ReadAll(p.Context.Request.Body)
	seelog.Info("ctx.Request.body: %v", string(data))
}

func (p *Param) AllParameter() string {
	m := make(map[string]interface{})
	m["param"] = p.Param
	m["form"] = p.Context.Request.Form
	b, _ := json.Marshal(m)
	return string(b)
}

func (p *Param) String(n string) string {
	v, vb := p.Param[n]
	if vb {
		return fmt.Sprint(v)
	}
	vx := p.Context.GetHeader(n)
	if vx != "" {
		return vx
	}
	sut, e := p.Context.Cookie(n)
	if e == nil && len(sut) > 0 {
		return sut
	}
	v2 := p.Context.Param(n)
	if v2 == "" {
		v2 = p.Context.PostForm(n)
		if v2 == "" {
			v2 = p.Context.Query(n)
			if v2 == "" {
				v2, _ = p.Context.GetQuery(n)
			}
		}
	}
	return v2
}

func (p *Param) Int64(n string) int64 {
	return p.Int64Default(n, 0)
}
func (p *Param) Int64Default(n string, def int64) int64 {
	vi := p.String(n)
	if vi == "" {
		return def
	}
	i, e := strconv.ParseInt(vi, 10, 0)
	if e != nil {
		return def
	}
	return i
}
func (p *Param) Start() int64 {
	return p.StartPageSize(DefaultPageSize)
}
func (p *Param) StartPageSize(ps int64) int64 {
	page := p.Int64Default("page", 1)
	if page < 1 {
		return int64(0)
	}
	return (page - 1) * ps
}
func (p *Param) Result(result ...interface{}) {
	if result != nil {
		p.Out["result"] = result
	}
}
func (p *Param) ST(st *utils.ST, result ...interface{}) {
	p.Out["code"] = st.Code
	p.Out["msg"] = st.Msg
	p.Result(result...)
}
func (p *Param) OK(result ...interface{}) {
	p.ST(utils.OK, result ...)
}
func NewParam(c *gin.Context) *Param {
	web := &Param{Context: c}
	web.ContentType = c.ContentType()
	if strings.Index(web.ContentType, "json") > 0 {
		row, b := web.Context.GetRawData()
		if b == nil && len(row) > 0 {
			je := json.Unmarshal(row, &web.Param)
			if je != nil {
				seelog.Error("param error", je, ";\n\t", string(row))
				//	web.Context.Set("param", web.Param)
			}
		}
	}
	web.Ua = web.String("Ua")
	if web.Ua == "" {
		web.Ua = GetUa(c)
	}
	web.Out = make(map[string]interface{})
	return web
}
