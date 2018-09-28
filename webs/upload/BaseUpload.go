package upload

import (
	"strings"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"os"
	"github.com/gpmgo/gopm/modules/log"
	"io"
	"love/ws"
	"github.com/cihub/seelog"
	"google.golang.org/grpc"
	"context"
	"utils/webs"
	"utils"
)

func UpInit(RpcUserHost, DirUpload, url, TokenName string, doFile func(ctx *gin.Context, userId int64, out map[string]interface{}), ImgMaxWidth ... int) {
	tmpDir := DirUpload + "/temp/"
	os.MkdirAll(tmpDir, 0777)
	seelog.Info("tmpDir=", tmpDir)
	ws.WebGin.POST(url, func(ctx *gin.Context) {
		conn, err := grpc.DialContext(context.Background(), RpcUserHost, grpc.WithInsecure())
		if err == nil {
			defer conn.Close()
			client := webs.NewRpcUserClient(conn)
			wb := webs.WebBaseNew(ctx)
			ub, _ := client.Verify(context.Background(), &webs.Token{Token: wb.String(TokenName), Ua: wb.Ua})
			if ub.Result {
				out := make(map[string]interface{})
				mf, _ := ctx.MultipartForm()
				for k, _ := range mf.File {
					res := doUploadFileMd5(DirUpload, ctx, k, ImgMaxWidth...)
					if res != nil {
						break
					}
				}
				doFile(ctx, ub.UserId, out)
				ctx.JSON(200, out)
			} else {
				seelog.Error("upload file auth fail.")
				ctx.Status(401)
			}
		}
	})
}

func doUploadFileMd5(DirUpload string, c *gin.Context, key string, ImgMaxWidth ... int) error {
	tmp := strconv.FormatInt(time.Now().UnixNano(), 16)
	file, header, err := c.Request.FormFile(key)

	filename := header.Filename
	ext := ".png"
	index := strings.Index(filename, ".")
	if index > 0 {
		ext = filename[index:]
	}
	tmpFileName := DirUpload + "/temp/" + tmp + ext
	out, err := os.Create(tmpFileName)
	if err != nil {
		log.Error("上传文件创建临时文件夹失败!", tmpFileName, err)
		return err
	}
	_, err2 := io.Copy(out, file)
	out.Close()
	file.Close()
	if err2 != nil {
		log.Error("写入上传文件流时失败!", tmpFileName, err)
		return err2
	}
	md5Name, e := Md5File(tmpFileName)
	if e == nil {
		pre, uri := utils.FmtImgDir(DirUpload+"/", md5Name)
		path := pre + ext
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Rename(tmpFileName, path)
		} else {
			os.Remove(tmpFileName)
		}
		c.Set(key+".pre", pre)
		if strings.ToLower(ext) != ".gif" {
			for idx, w := range ImgMaxWidth {
				ext8 := "_" + strconv.Itoa(w) + ext
				if _, err := os.Stat(pre + ext8); os.IsNotExist(err) {
					ImgResize(path, pre+ext8, w)
				}
				uri = uri + ext8
				c.Set(key+".uri."+strconv.Itoa(idx), uri)
			}
		} else {
			uri = uri + ext
			c.Set(key+".uri", uri)
		}
		return nil
	}
	return nil
}
