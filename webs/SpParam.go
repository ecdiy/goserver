package webs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"fmt"
	"strconv"
	"github.com/cihub/seelog"
	"time"
)

type Param struct {
	Auth, Param map[string]interface{}
	Ua          string
	Context     *gin.Context
}

func (p *Param) String(n string) string {
	return fmt.Sprint(p.Param[n])
}
func (p *Param) Username() string {
	return fmt.Sprint(p.Auth["Username"])
}
func (p *Param) Int64(n string, def int64) int64 {
	i, e := strconv.ParseInt(fmt.Sprint(p.Param[n]), 10, 0)
	if e != nil {
		return def
	}
	return i
}

func Parameter(c *gin.Context) (*Param, error) {
	ua, _ := c.Cookie("ua")
	if ua == "web" {
		ua = "web"
	} else {
		ua = "h5"
	}
	p, e := c.Get("param")
	if e {
		return p.(*Param), nil
	}
	row, b := c.GetRawData()
	if b == nil {
		var data map[string]interface{}
		je := json.Unmarshal(row, &data)
		if je != nil {
			seelog.Error("RawData JSON error:", row, je)
			return &Param{Ua: ua, Context: c}, je
		} else {
			px := &Param{Param: data, Ua: ua, Context: c}
			c.Set("param", px)
			return px, je
		}
	}
	return &Param{Ua: ua, Context: c}, b
}

func Post(Gin *gin.Engine, url string, fun func(param *Param, res map[string]interface{})) {
	Gin.POST(url, func(c *gin.Context) {
		param, e := Parameter(c)
		if e == nil {
			res := make(map[string]interface{})
			fun(param, res)
			res["now"] = time.Now().Format("2006-01-02T15:04:05Z")
			c.JSON(200, res)
		}
	})
}

func Auth(Gin *gin.Engine, url string, fun func(userId int64, param *Param, res map[string]interface{}), verify func(c *gin.Context) (bool, int64)) {
	Gin.POST(url, func(c *gin.Context) {
		auth, userId := verify(c)
		if auth {
			param, e := Parameter(c)
			if e == nil {
				param.Auth = c.Keys
				res := make(map[string]interface{})
				fun(userId, param, res)
				res["now"] = time.Now().Format("2006-01-02T15:04:05Z")
				c.JSON(200, res)
			}
		} else {
			c.AbortWithStatus(401)
		}
	})
}

//func GinGetAcao(Gin *gin.Engine, relativePath string, hand HandlerResult) {
//	Gin.GET(relativePath, func(c *gin.Context) {
//		res := hand(c)
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Header("Content-Type", "application/json; charset=utf-8")
//		c.JSON(http.StatusOK, res)
//	})
//}
