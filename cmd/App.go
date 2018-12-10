package main

import (
	_ "github.com/ecdiy/goserver/plugins/sp"
	_ "github.com/ecdiy/goserver/plugins/verify"
	_ "github.com/ecdiy/goserver/plugins/web"
	_ "github.com/ecdiy/goserver/plugins/web/captcha"
	_ "github.com/ecdiy/goserver/plugins/web/file"
	_ "github.com/ecdiy/goserver/plugins/web/image/resize"
	_ "github.com/ecdiy/goserver/plugins/web/sp"
	_ "github.com/ecdiy/goserver/plugins/web/upload"
	_ "github.com/ecdiy/goserver/plugins/web/websocket"
	_ "github.com/ecdiy/goserver/plugins/web/webexec"
	"github.com/ecdiy/goserver/plugins"
)

func main() {
	plugins.StartCore()
}
