package core

import (
	"goserver/utils"
	"github.com/gin-gonic/gin"
	"goserver/webs"
	"io/ioutil"
	"gopkg.in/russross/blackfriday.v2"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
	"strings"
	"strconv"
)

func (app *Module) Md(ele *utils.Element) {
	dir := ele.MustAttr("MdDir")
	fun := func(context *gin.Context) {
		wp := webs.NewParam(context)
		md := wp.String("Md")
		bs, e := ioutil.ReadFile(dir + md)
		if e == nil {
			context.JSON(200, Md2Html(bs))
		} else {
			context.JSON(200, "")
		}
	}
	gj, gje := ele.AttrValue("GetJsonUrl")
	if gje {
		getGin(ele).GET(gj, fun)
	}
	pj, pje := ele.AttrValue("PostJsonUrl")
	if pje {
		getGin(ele).POST(pj, fun)
	}
}

func Md2Html(bs []byte) string {
	unsafe := blackfriday.Run(bs)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	h := string(html)
	re := regexp.MustCompile("<h[1-4]>")
	ss := re.FindAllString(h, -1)
	str := ""
	for i, s := range ss {
		idx := strings.Index(h, s)
		if idx >= 0 {
			str = str + h[0:idx+len(s)] + `<a id="` + strconv.Itoa(i) + `"></a>`
			h = h[idx+len(s):]
		}
	}
	str += h
	return str
}
