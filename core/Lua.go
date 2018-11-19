package core

import (
	"github.com/ecdiy/goserver/utils"
	"io/ioutil"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/yuin/gopher-lua"
	"github.com/ecdiy/goserver/gpa"
	"github.com/ecdiy/goserver/webs"
	"encoding/json"
	"github.com/ecdiy/goserver/utils/luar"
)

func (app *Module) Lua(ele *utils.Element) {
	sl := &SpLua{uri: ele.MustAttr("Url"), Gpa: getGpa(ele), IsDev: utils.EnvIsDev,
		ContentType: ele.Attr("ContentType", "application/json; charset=utf-8"),
		LuaScriptMap: make(map[string]string), LuaScriptDir: ele.MustAttr("Dir")}
	getGin(ele).NoRoute(sl.DoWebContext)
}

type SpLua struct {
	IsDev                          bool
	LuaScriptMap                   map[string]string
	uri, LuaScriptDir, ContentType string
	Gpa                            *gpa.Gpa
}

func (sl *SpLua) GetLuaScript(uri string) (string, bool) {
	path := uri
	script, b := sl.LuaScriptMap[uri]
	if !b || sl.IsDev {
		lf := sl.LuaScriptDir + path + ".lua"
		bs, e := ioutil.ReadFile(lf)
		if e != nil || len(bs) == 0 {
			seelog.Error(lf)
			return "", false
		}
		script = string(bs)
		if len(script) > 0 {
			sl.LuaScriptMap[uri] = script
		}
	}
	return script, true
}

func (sl *SpLua) DoWebContext(ctx *gin.Context) {
	uri := ctx.Request.URL.Path
	if strings.Index(uri, sl.uri) == 0 {
		path := uri[len(sl.uri):]
		ctx.Header("Content-Type", sl.ContentType)
		script, ext := sl.GetLuaScript(path)
		if ext {
			p := webs.NewParam(ctx)
			L := lua.NewState()
			defer L.Close()
			sl.LuaFun(L, p)
			err := L.DoString(script)
			if err == nil && L.GetTop() == 1 {
				tb := L.CheckTable(1)
				m := make(map[string]interface{})
				tb.ForEach(func(value lua.LValue, value2 lua.LValue) {
					m[value.String()] = luaConvertLValueToInterface(value2)
				})
				js1, _ := json.Marshal(m)
				ctx.String(200, string(js1))
			} else {
				seelog.Warn("no return.", L.GetTop(), "\t", err)
			}
		}
	}
}

func luaConvertLValueToInterface(value2 lua.LValue) interface{} {
	if value2.Type().String() == "userdata" {
		ud := value2.(*lua.LUserData)
		return ud.Value
	} else if value2.Type().String() == "table" {
		return luaConvertTableToMap(value2)
	} else {
		return value2
	}
}

func luaConvertTableToMap(lv lua.LValue) map[string]interface{} {
	m := make(map[string]interface{})
	ut := lv.(*lua.LTable)
	ut.ForEach(func(value lua.LValue, value2 lua.LValue) {
		if value2.Type().String() == "userdata" {
			ud := value2.(*lua.LUserData)
			m[value.String()] = ud.Value
		} else if value2.Type().String() == "table" {
			m[value.String()] = luaConvertTableToMap(value2)
		} else {
			m[value.String()] = value2
		}
	})
	return m
}

func (sl *SpLua) LuaFun(L *lua.LState, param *webs.Param) {

	L.SetGlobal("db", luar.New(L, sl.Gpa))

	L.SetGlobal("param", L.NewFunction(func(state *lua.LState) int {
		name := state.ToString(1)
		if state.GetTop() == 1 {
			state.Push(luar.New(state, param.String(name)))
			return 1
		}
		return 0
	}))

	//L.SetGlobal("find", L.NewFunction(func(state *lua.LState) int {
	//	t := state.GetTop()
	//	bd := state.ToString(1)
	//	if t == 2 { // body findReg ==> [][]string
	//		reg := state.ToString(2)
	//		mc := regexp.MustCompile(reg)
	//		state.Push(luar.New(L, mc.FindAllStringSubmatch(bd, -1)))
	//		return 1
	//	}
	//	if t == 3 { // body findReg group ==> string
	//		reg := state.ToString(2)
	//		mc := regexp.MustCompile(reg)
	//		res := mc.FindAllStringSubmatch(bd, -1)
	//		group := state.ToInt(3)
	//		if len(res) >= 1 && len(res[0]) > group {
	//			state.Push(lua.LString(res[0][group]))
	//		} else {
	//			state.Push(lua.LString(""))
	//		}
	//		return 1
	//	}
	//	if t == 5 { //body,[1,-1], beginFlag ,[1,-1], endFlag, reg ==>[][]string
	//		beginFlag := L.ToString(2)
	//		endFlag := L.ToString(3)
	//
	//		flg := L.ToInt(5)
	//		if flg&1 == 1 { //1   去掉\r\n
	//			bd = trimLineRegexp.ReplaceAllString(bd, "")
	//		}
	//		bd = lang.StringFindBegin(bd, beginFlag, flg&2 == 2)
	//		if bd == "" {
	//			state.Push(lua.LNil)
	//			return 1
	//		}
	//		bd = lang.StringFindEnd(bd, endFlag, flg&4 == 4)
	//		if bd == "" {
	//			state.Push(lua.LNil)
	//			return 1
	//		}
	//		mc := regexp.MustCompile(L.ToString(4))
	//		state.Push(luar.New(L, mc.FindAllStringSubmatch(bd, -1)))
	//		return 1
	//	}
	//	state.Push(lua.LString(""))
	//	seelog.Info(`参数个数错误[2,3,5]`, t)
	//	return 1
	//}))

	//L.SetGlobal("indexOf", L.NewFunction(func(ns *lua.LState) int {
	//	ns.Push(luar.New(ns, strings.Index(ns.ToString(1), ns.ToString(2))))
	//	return 1
	//}))
	//L.SetGlobal("split", L.NewFunction(func(ns *lua.LState) int {
	//	ns.Push(luar.New(ns, strings.Split(ns.ToString(1), ns.ToString(2))))
	//	return 1
	//}))
	//
	//L.SetGlobal("bitAnd", L.NewFunction(func(state *lua.LState) int { //与
	//	state.Push(luar.New(state, state.ToInt64(1)&state.ToInt64(2)))
	//	return 1
	//}))
	//
	//L.SetGlobal("bitOr", L.NewFunction(func(state *lua.LState) int { //或
	//	state.Push(luar.New(state, state.ToInt64(1)|state.ToInt64(2)))
	//	return 1
	//}))
	//
	//L.SetGlobal("json", L.NewFunction(func(state *lua.LState) int {
	//	str := state.ToString(1)
	//	mj, ee := gjson.Parse(str).Value().(map[string]interface{})
	//	state.Push(luar.New(state, mj))
	//	state.Push(luar.New(state, ee))
	//	return 2
	//}))
	//
	//L.SetGlobal("toString", L.NewFunction(func(state *lua.LState) int {
	//	v := luaConvertLValueToInterface(state.CheckAny(1))
	//	js1, err := json.Marshal(v)
	//	state.Push(luar.New(state, string(js1)))
	//	state.Push(luar.New(state, err))
	//	return 2
	//}))

}
