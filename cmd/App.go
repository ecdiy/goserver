package main

import (
	_ "github.com/ecdiy/goserver/plugins/sp"                  //存储过程
	_ "github.com/ecdiy/goserver/plugins/verify"              //验证
	_ "github.com/ecdiy/goserver/plugins/web"                 //web
	_ "github.com/ecdiy/goserver/plugins/web/file"            //文件
	_ "github.com/ecdiy/goserver/plugins/web/image/captcha"   //验证码
	_ "github.com/ecdiy/goserver/plugins/web/image/ImgResize"    //缩放
	_ "github.com/ecdiy/goserver/plugins/web/image/ImgBase64" //缩放
	_ "github.com/ecdiy/goserver/plugins/web/image/QrCode"    //二维码
	_ "github.com/ecdiy/goserver/plugins/web/sp"              //URL映射存储过程
	_ "github.com/ecdiy/goserver/plugins/web/upload"          //上传
	_ "github.com/ecdiy/goserver/plugins/web/ws"              //WebSocket
	_ "github.com/ecdiy/goserver/plugins/web/webexec"         //WebExec
	"github.com/ecdiy/goserver/plugins"
)

func main() {
	plugins.StartCore()
}
