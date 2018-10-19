package main

import (
	"utils/webs"
	"google.golang.org/grpc"
	"utils/xml"
	"utils"
	"utils/gpa"
	"github.com/cihub/seelog"
)

type Module struct {
}

func (app *Module) Include(ele *xml.Element) {
	f := getFile(ele.Value)
	seelog.Info("include file:", f)
	dom, err := xml.LoadByFile(f)
	if err == nil {
		InvokeByXml(dom)
	} else {
		seelog.Error(err)
		panic("配置文件出错:" + f)
	}
}

func (app *Module) Parameter(ele *xml.Element) {
	ps := ele.AllNodes()
	for _, p := range ps {
		utils.EnvParamSet(p.Name(), p.Value)
	}
}

func (app *Module) Gpa(ele *xml.Element) {
	dsn, b := ele.AttrValue("DbDsn")
	if b && len(dsn) > 0 {
		db := gpa.InitGpa(dsn)
		put(ele, db)
	}
}

func (app *Module) Rpc(ele *xml.Element) {
	Sql := ele.MustAttr("Sql")
	rpc := &webs.RpcUser{Sql: Sql, Gpa: getGpa(ele)}
	RpcHost := ele.MustAttr("RpcHost")
	putFunRun(func() {
		webs.RpcRegister(RpcHost, func(s *grpc.Server) {
			webs.RegisterRpcUserServer(s, rpc)
		})
	})
	put(ele, rpc)
}
