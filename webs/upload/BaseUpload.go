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
)

func UpInit(RpcUserHost, DirUpload, url, TokenName string, ImgMaxWidth int, doFile func(userId int64, uri, pre string, out map[string]interface{})) {
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
					item := make(map[string]interface{})
					res := doUploadFileMd5(ub.UserId, DirUpload, ImgMaxWidth, ctx, k, item, doFile)
					out[k] = item
					if res != nil {
						break
					}
				}
				ctx.JSON(200, out)
			} else {
				seelog.Error("upload file auth fail.")
				ctx.Status(401)
			}
		}
	})
}

func doUploadFileMd5(userId int64, DirUpload string, ImgMaxWidth int, c *gin.Context,
	key string, m map[string]interface{},
	doFile func(userId int64, uri, pre string, out map[string]interface{})) error {
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
		pre, uri := FmtImgDir(DirUpload+"/", md5Name)
		path := pre + ext
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Rename(tmpFileName, path)
		} else {
			os.Remove(tmpFileName)
		}
		if strings.ToLower(ext) != ".gif" {
			if ImgMaxWidth > 0 {
				ext8 := "_" + strconv.Itoa(ImgMaxWidth) + ext
				if _, err := os.Stat(pre + ext8); os.IsNotExist(err) {
					ImgResize(path, pre+ext8, ImgMaxWidth)
				}
				uri = uri + ext8
			} else {
				uri = uri + ext
			}
		} else {
			uri = uri + ext
		}
		doFile(userId, uri, pre, m)
		return nil
	}
	return nil
}
