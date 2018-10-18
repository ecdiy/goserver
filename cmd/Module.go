package main

import (
	"utils/webs"
	"google.golang.org/grpc"
	"github.com/cihub/seelog"
	"utils/xml"
	"utils"
	"utils/gpa"
)

type Module struct {
}


// <RpcSp Sql="select UserId from Token where Token=? and Ua=?" BindHost="" TokenName="token" Url="" Ext="Ajax"/>
func (app *Module) RpcSp(ele *xml.Element)   {

	Sql, b0 := ele.AttrValue("RpcHost")
	if !b0 {
		seelog.Warn("没有设置SQL,default=select UserId from Token where Token=? and Ua=?")
		Sql = "select UserId from Token where Token=? and Ua=?"
	}
	rpc := &webs.RpcUser{Sql: Sql}

	RpcHost, b0 := ele.AttrValue("RpcHost")
	if !b0 {
		panic("[RpcSp]缺少参数:RpcHost [127.0.0.1:9200]")
	}
	go func() {
		webs.RpcRegister(RpcHost, func(s *grpc.Server) {
			webs.RegisterRpcUserServer(s, rpc)
		})
	}()

	//ext, b0 := ele.AttrValue("Ext")
	//if !b0 {
	//	panic("[RpcSp]没有设置:Ext [Ajax]")
	//}
	//tn, b1 := ele.AttrValue("TokenName")
	//if !b1 {
	//	panic("[RpcSp]缺少参数:TokenName [token]")
	//}
	//url, b2 := ele.AttrValue("Url")
	//if !b2 {
	//	panic("[RpcSp]没有设置Ext,[/sp/:sp]")
	//}
	//webs.RegisterSpAjax(webs.GetAuthByRpc(rpc, tn), url, ext)

}

func (app *Module) RpcHost(ele *xml.Element)   {
	//RpcHost, b0 := ele.AttrValue("RpcHost")
	//if !b0 {
	//	panic("[RpcHost]缺少参数:RpcHost [127.0.0.1:9200]")
	//}
	//tn, b1 := ele.AttrValue("TokenName")
	//if !b1 {
	//	panic("[RpcHost]缺少参数:TokenName [token]")
	//}
	//ext, b0 := ele.AttrValue("Ext")
	//if !b0 {
	//	panic("[RpcHost]没有设置Ext,[Ajax]")
	//}
	//url, b2 := ele.AttrValue("Url")
	//if !b2 {
	//	panic("[RpcHost]没有设置Url,[/sp/:sp]")
	//}
	//webs.RegisterSpAjax(webs.GetAuthFunByHost(RpcHost, tn), url, ext)

}


func (app *Module) Parameter(ele *xml.Element) {
	ps := ele.AllNodes()
	for _, p := range ps {
		utils.EnvParamSet(p.Name(), p.Value)
	}
	//	for _, attr := range ecXml.Attrs {
	//		f := ecd.FieldByName(attr.Name())
	//		v := attr.Value
	//		v = funcs.GetVal(v)
	//		if f.CanSet() {
	//			switch f.Type().Name() {
	//			case "string":
	//				f.SetString(v)
	//				break
	//			case "bool":
	//				if "1" == v || "true" == v {
	//					f.SetBool(true)
	//				} else {
	//					f.SetBool(false)
	//				}
	//				break
	//			default:
	//				array := strings.Split(v, ",")
	//				f.Set(reflect.ValueOf(array))
	//			}
	//

}

func (app *Module) Gpa(ele *xml.Element) {
	dsn, b := ele.AttrValue("DbDsn")
	if b && len(dsn) > 0 {
		db := gpa.InitGpa(dsn)
		put(ele, db)
	}
}

func (app *Module) Rpc(ele *xml.Element) {
	Sql := ele.MustAttr("RpcHost")
	rpc := &webs.RpcUser{Sql: Sql, Gpa: getGpa(ele)}
	RpcHost := ele.MustAttr("RpcHost")
	putFunRun(func() {
		webs.RpcRegister(RpcHost, func(s *grpc.Server) {
			webs.RegisterRpcUserServer(s, rpc)
		})
	})
	put(ele, rpc)
}
