package ImgBase64

import (
	"github.com/ecdiy/goserver/plugins/web/image"
	"github.com/ecdiy/goserver/plugins/web"
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins"
	"strings"
	"encoding/base64"
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"io/ioutil"
)

func init() { /**
base 64 保存成图片.
 */
	web.RegisterWeb("ImgBase64", func(ele *utils.Element) func(c *utils.Param) {

		bi := &Base64Img{
			OutFile:  ele.MustAttr("Uri"),
			Dir:      ele.MustAttr("Dir"),
			vf:       plugins.GetRef(ele, "Verify").(plugins.BaseFun),
			Param:    strings.Split(ele.Attr("Param", ""), ","),
			GinParam: strings.Split(ele.Attr("GinParam", ""), ","),
		}
		if len(bi.Param) == 0 && len(bi.GinParam) == 0 {
			return bi.SaveByMd5
		} else {
			return bi.SaveByParamName
		}

	})
}

type Base64Img struct {
	GinParam, Param []string
	Dir, OutFile    string
	vf              plugins.BaseFun
}

func (bi *Base64Img) SaveByMd5(wb *utils.Param) {
	bi.vf(wb)
	data := wb.String("data")
	md5 := image.Md5Byte([]byte(data))
	path, uri := image.FmtImgDir(bi.Dir, md5)
	wb.Out["Uri"] = uri + ".jpg"
	ddd, _ := base64.StdEncoding.DecodeString(data) //成图片文件并把文件写入到buffer
	ioutil.WriteFile(path+".jpg", ddd, 0666)
}

func (bi *Base64Img) SaveByParamName(wb *utils.Param) {
	bi.vf(wb)
	pp := bi.OutFile
	err := false
	for _, p := range bi.GinParam {
		v, vb := wb.Context.Get(p)
		if vb {
			pp = strings.Replace(pp, "${"+p+"}", fmt.Sprint(v), -1)
		} else {
			seelog.Error("Gin参数错误:", p)
			wb.Out["ParaErr"] = p
			err = true
			break
		}
	}
	if !err {
		for _, p := range bi.Param {
			vs := wb.String(p)
			if len(vs) > 0 {
				pp = strings.Replace(pp, "${"+p+"}", vs, -1)
			} else {
				seelog.Error("参数错误:", p)
				wb.Out["ParaErr"] = p
				err = true
				break
			}
		}
		if !err {
			wb.Out["Uri"] = pp + ".jpg"
			ix := strings.LastIndex(pp, "/")
			if ix > 0 {
				pDir := bi.Dir + pp[0:ix]
				os.MkdirAll(pDir, 0644)
			}
			data := wb.String("data")
			flg := "data:image/png;base64,"
			if data[0:len(flg)] == flg {
				data = data[len(flg):]
			}
			ddd, _ := base64.StdEncoding.DecodeString(data) //成图片文件并把文件写入到buffer
			ioutil.WriteFile(bi.Dir+pp+".jpg", ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
		}
	}
}

//func   UploadBase64Img(c *gin.Context) {
//	ds := c.PostForm(u.ParamName)
//	ext := ".png"
//	index := strings.Index(ds, ",")
//	if index > 0 {
//		ext := ds[0:index]
//		if strings.Index(ext, "jpeg") > 0 {
//			ext = ".jpg"
//		}
//		ds = ds[index+1:]
//	}
//	bs, _ := base64.StdEncoding.DecodeString(ds)
//	path := UploadDir + FmtImgDir(Md5Byte(bs))
//	souFile := path + ext
//	if f, err := os.Stat(souFile); err == nil && f.Size() == int64(len(bs)) {
//		log.Info("上传文件已经存在:", souFile, ";", f.Size(), ";", len(bs))
//	} else {
//		log.Info("上传文件:", souFile)
//		out, err := os.Create(souFile)
//		if err != nil {
//			log.Error(err)
//		}
//		defer out.Close()
//		out.Write(bs)
//		if u.LimitWidth != nil && len(u.LimitWidth) > 0 {
//			for _, w := range u.LimitWidth {
//				ImgResize(souFile, path+"_"+strconv.Itoa(w)+ext, w)
//			}
//		}
//	}
//}

//func Upload(c *gin.Context) {
//name := c.Query("name")
//ext := ".jpg"
//idx := strings.Index(name, ".")
//if idx > 0 {
//	ext = name[idx:]
//}
//tmp := DirUpload + strconv.FormatInt(time.Now().UnixNano(), 10) + ext;
//c.SaveToFile("upfile", tmp)
//json := make(map[string]string)
//md5Uri, size := fileutil.GetUploadFile(tmp)
//md5Uri += ext
//fullUrl := UrlUpload + md5Uri
//data := make(map[string]interface{})
//c.Ctx.Input.SetParam("url", fullUrl)
//e, st := ecsp.WebExec(data, c.Ctx, sysBaseUpload, "")
//if e == nil && st.Code == util.ST_OK.Code {
//	fileutil.Rename(tmp, DirUpload + md5Uri)
//	json["state"] = "SUCCESS"
//	json["type"] = ext
//	json["original"] = name
//	json["name"] = name
//	json["url"] = fullUrl[1:]
//	json["fullUrl"] = fullUrl
//	json["size"] = strconv.FormatInt(size, 16)
//	c.Data["json"] = json
//	c.ServeJSON()
//} else {
//	os.Remove(tmp)
//	c.Ctx.ResponseWriter.WriteHeader(500)
//}
//}
