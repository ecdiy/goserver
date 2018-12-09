package upload

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/core"
	"github.com/gin-gonic/gin"
	"strings"
	"fmt"
	"github.com/ecdiy/itgeek/gk/upload"
	"github.com/ecdiy/goserver/plugins"
)

//TODO

func (app *Module) Upload(ele *utils.Element) {
	sp, spExt := ele.AttrValue("SpRef")
	var nameFun func(c *webs.Param, tmpFileName string) (string, error)

	NameRule, nrExt := ele.AttrValue("NameRule")
	if nrExt {
		if strings.Index(NameRule, "UserId") == 0 {
			nameFun = func(c *webs.Param, tmpFileName string) (string, error) {
				UserId, _ := c.Context.Get("UserId")
				xId := fmt.Sprint(UserId)
				return xId, nil
			}
		}
		if strings.Index(NameRule, "Md5") == 0 {
			nameFun = func(c *webs.Param, tmpFileName string) (string, error) {
				return upload.Md5File(tmpFileName)
			}
		}
	} else {
		nameFun = func(c *webs.Param, tmpFileName string) (string, error) {
			return upload.Md5File(tmpFileName)
		}
	}

	if spExt {
		upload.Upload(nameFun, plugins.GetGin(ele), plugins.Data[sp].(*webs.WebSp), getVerify(ele), ele)
	} else {
		upload.Upload(nameFun, plugins.GetGin(ele), nil, getVerify(ele), ele)
	}
}
func init() {
	core.WebPlugins["File"] = WebFilePlugin
}

func WebFilePlugin(xml *utils.Element) func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}

//
//
//func (app *Module) Upload(ele *utils.Element) {
//	sp, spExt := ele.AttrValue("SpRef")
//	var nameFun func(c *webs.Param, tmpFileName string) (string, error)
//
//	NameRule, nrExt := ele.AttrValue("NameRule")
//	if nrExt {
//		if strings.Index(NameRule, "UserId") == 0 {
//			nameFun = func(c *webs.Param, tmpFileName string) (string, error) {
//				UserId, _ := c.Context.Get("UserId")
//				xId := fmt.Sprint(UserId)
//				return xId, nil
//			}
//		}
//		if strings.Index(NameRule, "Md5") == 0 {
//			nameFun = func(c *webs.Param, tmpFileName string) (string, error) {
//				return upload.Md5File(tmpFileName)
//			}
//		}
//	} else {
//		nameFun = func(c *webs.Param, tmpFileName string) (string, error) {
//			return upload.Md5File(tmpFileName)
//		}
//	}
//
//	if spExt {
//		upload.Upload(nameFun, getGin(ele), data[sp].(*webs.WebSp), getVerify(ele), ele)
//	} else {
//		upload.Upload(nameFun, getGin(ele), nil, getVerify(ele), ele)
//	}
//}
