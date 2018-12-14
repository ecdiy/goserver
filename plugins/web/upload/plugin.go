package upload

import (
	"github.com/ecdiy/goserver/utils"
	"strings"
	"github.com/ecdiy/goserver/plugins"
	"github.com/ecdiy/goserver/plugins/web"
	"github.com/ecdiy/goserver/plugins/sp"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
)

func init() {
	web.RegisterWebPlugin("Upload", func(ele *utils.Element) func(c *gin.Context) {
		wu := &WebUpload{UrlPrefix: ele.MustAttr("UrlPrefix"), DirUpload: ele.MustAttr("DirUpload")}
		wu.setNameFun(ele)
		wu.baseFun = plugins.GetRef(ele, "Verify").(plugins.BaseFun)
		CoverVal, CoverExt := ele.AttrValue("Cover")
		wu.ToExt = "." + ele.Attr("ToExt", "png")
		wu.Cover = false
		if CoverExt && CoverVal == "1" {
			wu.Cover = true
		}

		spName, spExt := ele.AttrValue("Sp")
		if spExt {
			wu.spName = spName
			wu.webSp = &sp.WebSp{}
			wu.webSp.Init(ele)
		}
		os.MkdirAll(wu.DirUpload, 0777)
		edc := wu.DirUpload[len(wu.DirUpload)-1:]
		if edc != "/" && edc != "/" {
			wu.DirUpload += "/"
		}
		wu.tmpDir = ele.Attr("TmpDir", wu.DirUpload+"temp/")
		os.MkdirAll(wu.tmpDir, 0777)
		var ImgWidth []int
		iw, iwb := ele.AttrValue("ImgWidth")
		if iwb {
			iws := strings.Split(iw, ",")
			for _, w := range iws {
				wi, _ := strconv.Atoi(w)
				ImgWidth = append(ImgWidth, wi)
			}
		}
		wu.MainWidth, _ = strconv.Atoi(ele.Attr("MainWidth", "0"))
		return wu.Upload
	})
}
