package upload

import (
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins/sp"
	"github.com/ecdiy/goserver/plugins"
	"strings"
	"fmt"
	"github.com/ecdiy/itgeek/gk/upload"
	"github.com/ecdiy/goserver/plugins/verify"
	"github.com/cihub/seelog"
	"strconv"
	"time"
	"os"
	"github.com/gpmgo/gopm/modules/log"
	"io"
	"github.com/ecdiy/goserver/plugins/web/image"
	"github.com/gin-gonic/gin"
)

//文件上传
/**

  window 下， os.Rename 需要在同一磁盘下，否则会不成功。
TmpDir,MainWidth,ImgWidth（多个用,分隔） 可选

	两种方式:
    1.上传完成后不调用存储过程，输出 文件名对应的参数Url
    <Upload WebRef="Web" TmpDir="./upload/temp/" DirUpload="./upload/" ImgWidth="800" MainWidth="800" UrlPrefix="/upload"/>

	2.上传完成后调用存储过程， 存储过程返回值决定输出内容
    <Upload SpRef="Sp" WebRef="Web" TmpDir="./upload/temp/" DirUpload="./upload/" ImgWidth="800" Sp="Upload" MainWidth="800" UrlPrefix="/upload"/>

 */
type WebUpload struct {
	nameFun                                     func(c *utils.Param, tmpFileName string) (string, error)
	MainWidth                                   int
	ImgWidth                                    []int
	tmpDir, spName, UrlPrefix, DirUpload, ToExt string
	webSp                                       *sp.WebSp
	Cover                                       bool
	baseFun                                     plugins.BaseFun
}

func (wu *WebUpload) setNameFun(ele *utils.Element) {
	NameRule, nrExt := ele.AttrValue("NameRule")
	if nrExt {
		if strings.Index(NameRule, "UserId") == 0 {
			wu.nameFun = func(c *utils.Param, tmpFileName string) (string, error) {
				UserId, _ := c.Context.Get("UserId")
				xId := fmt.Sprint(UserId)
				return xId, nil
			}
		}
		if strings.Index(NameRule, "Md5") == 0 {
			wu.nameFun = func(c *utils.Param, tmpFileName string) (string, error) {
				return upload.Md5File(tmpFileName)
			}
		}
	} else {
		wu.nameFun = func(c *utils.Param, tmpFileName string) (string, error) {
			return upload.Md5File(tmpFileName)
		}
	}
}

func (wu *WebUpload) Upload(c *gin.Context) {
	wb := utils.NewParam(c)
	ub := wu.baseFun(wb).(*verify.UserBase)

	if ub != nil && ub.Result {
		mf, _ := wb.Context.MultipartForm()
		for k, _ := range mf.File {
			wu.doUploadFileMd5(wb, k)
		}
		if wu.webSp != nil {
			wu.webSp.SpExec(wu.spName, wb)
			wb.Context.JSON(200, wb.Out)
		} else {
			wb.Context.JSON(200, wb.Param)
		}
	} else {
		seelog.Error("upload file auth fail.")
		wb.Context.Status(401)
	}
}

func (wu *WebUpload) doUploadFileMd5(c *utils.Param, key string) {
	tmp := strconv.FormatInt(time.Now().UnixNano(), 16)
	file, header, err := c.Context.Request.FormFile(key)

	filename := header.Filename
	ext := ".png"
	index := strings.Index(filename, ".")
	if index > 0 {
		ext = filename[index:]
	}
	tmpFileName := wu.tmpDir + tmp + ext
	out, err := os.Create(tmpFileName)
	if err != nil {
		log.Error("上传文件创建临时文件夹失败!", tmpFileName, err)
		return
	}
	_, err2 := io.Copy(out, file)
	out.Close()
	file.Close()
	if err2 != nil {
		log.Error("写入上传文件流时失败!", tmpFileName, err)
		return
	}
	md5Name, e := wu.nameFun(c, tmpFileName)
	if e == nil {
		pre, uri := utils.FmtImgDir(wu.DirUpload, md5Name)
		path := pre + ext
		if _, err := os.Stat(path); wu.Cover || os.IsNotExist(err) {
			os.Rename(tmpFileName, path)
		} else {
			os.Remove(tmpFileName)
		}
		if strings.ToLower(ext) != ".gif" && wu.ImgWidth != nil {
			for _, w := range wu.ImgWidth {
				ext8 := "_" + strconv.Itoa(w) + wu.ToExt
				if wu.Cover {
					os.Remove(pre + ext8)
				}
				if _, err := os.Stat(pre + ext8); os.IsNotExist(err) {
					image.ImgResize(path, pre+ext8, w)
				}
				if wu.MainWidth == 0 || wu.MainWidth == w || len(wu.ImgWidth) == 1 {
					ut := uri + ext8
					c.Param[key+"Url"] = wu.UrlPrefix + ut
					c.Param["location"] = wu.UrlPrefix + ut //for TinyMce
				}
			}
		} else {
			c.Param[key+"Url"] = wu.UrlPrefix + uri + ext
			c.Param["location"] = wu.UrlPrefix + uri + ext //for TinyMce
		}
	}
}
