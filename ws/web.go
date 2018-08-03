package ws

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
	"strconv"
)

type Web struct {
	param, Out map[string]interface{}
	Ua         string
	Context    *gin.Context
	//---
	Auth          bool
	UserId, Score int64
	Username      string
}

func (p *Web) initParam() {
	if p.param == nil {
		row, b := p.Context.GetRawData()
		if b == nil {
			var data map[string]interface{}
			je := json.Unmarshal(row, &data)
			if je == nil {
				p.param = data
			}
		}
	}
}
func (p *Web) String(n string) string {
	return fmt.Sprint(p.param[ n])
}
func (p *Web) Int64(n string) int64 {
	return p.Int64Default(n, 0)
}
func (p *Web) Int64Default(n string, def int64) int64 {
	i, e := strconv.ParseInt(fmt.Sprint(p.param[ n]), 10, 0)
	if e != nil {
		return def
	}
	return i
}
func (p *Web) Start() int64 {
	return p.StartPageSize(PageSize)
}
func (p *Web) StartPageSize(ps int64) int64 {
	page := p.Int64Default("page", 1)
	if page < 1 {
		return int64(0)
	}
	return (page - 1) * ps
}
func (p *Web) Result(result ... interface{}) {
	if result != nil {
		p.Out["result"] = result
	}
}
func (p *Web) ST(st *ST, result ... interface{}) {
	p.Out["code"] = st.Code
	p.Out["msg"] = st.Msg
	p.Result(result ...)
}

func GetUa(ctx *gin.Context) string {
	ua := ctx.Request.UserAgent()
	if len(ua) == 0 {
		return "web"
	}
	if UaH5.MatchString(ua) {
		return "h5"
	}
	if UaSeo.MatchString(ua) {
		return "seo"
	}
	return "web"
}

func byAuthFun(url string, fun func(wdb *Web), auth func(c *gin.Context) *Web) {
	WebGin.POST(url, func(c *gin.Context) {
		web := auth(c)
		if web.Auth {
			web.initParam()
			fun(web)
			if len(web.Out) > 0 {
				c.JSON(200, web.Out)
			}
		} else {
			c.AbortWithStatus(401)
		}
	})
}

func AdminAuth(url string, fun func(wdb *Web)) {
	byAuthFun(url, fun, VerifyAdmin)
}

func WebAuth(url string, fun func(wdb *Web)) {
	byAuthFun(url, fun, VerifyRpc)
}

func WebPost(url string, fun func(wdb *Web)) {
	WebGin.POST(url, func(c *gin.Context) {
		web := VerifyRpc(c)
		web.initParam()
		fun(web)
		c.JSON(200, web.Out)
	})
}

func WebBase(url string, fun func(wdb *Web)) {
	WebGin.POST(url, func(c *gin.Context) {
		web := &Web{}
		web.Ua = GetUa(c)
		web.Context = c
		web.Out = make(map[string]interface{})
		web.initParam()
		fun(web)
		c.JSON(200, web.Out)
	})
}
