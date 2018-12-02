package main

import (
	"github.com/ecdiy/goserver/core"

	_ "github.com/ecdiy/goserver/plugins/web/file"
)

func main() {
	core.StartCore()
}
