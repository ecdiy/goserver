package gpa

import (
	"database/sql"
	"github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

const (
	PrimaryId     = "@Id"
	AutoIncrement = "AutoIncrement"
)

var nilVf = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

type Gpa struct {
	driver, dsn string
	Conn        *sql.DB
}


//func InitGpa(dbName string, models ... interface{}) *Gpa {
//	return Init(base.Param(base.KeyDbDriverName),
//		base.Param(base.KeyDbUri)+"/"+dbName+"?timeout=30s&charset=utf8mb4&parseTime=true", models ...)
//}

func getSqlByMethod(ft reflect.StructField) string {
	name := ft.Name
	if strings.Index(name, "FindBy") == 0 {
		ty := ft.Type.String()
		lk := strings.LastIndex(ty, "(")
		if lk > 0 {
			ty = ty[lk+1:]
		}
		d := strings.Index(ty, ".")
		x := strings.Index(ty, ",")
		if d > 0 && x > d {
			tb := ty[d+1 : x]
			rep := strings.Replace(name[6:], "And", "=? And ", -1)
			return "select * from " + tb + " where " + rep + "=?"
		} else {
			seelog.Error("错误的命令格式:" + name + "," + ty)
		}
	}
	return ""
}

func (impl *Gpa) setMethodImpl(di interface{}) {
	toe := reflect.TypeOf(di).Elem()
	voe := reflect.ValueOf(di).Elem()
	implVO := reflect.ValueOf(&impl).Elem()
	for i := 0; i < voe.NumField(); i++ {
		ft := toe.Field(i)
		runSql := strings.TrimSpace(string(ft.Tag))
		if len(runSql) < 1 {
			runSql = getSqlByMethod(ft)
			if len(runSql) < 1 {
				seelog.Error("方法定义错误,没有设置RunSql:", ft.Name, ";", ft.Type.String())
				continue
			}
		}
		if strings.Index(ft.Type.String(), "(") < 0 {
			seelog.Error("不是函数" + ft.Name)
			continue
		}
		rt := ft.Type.String()
		rt = strings.TrimSpace(rt[strings.Index(rt, ")")+1:])

		fv := voe.Field(i)
		methodName := ""
		if len(rt) < 2 {
			seelog.Flush()
			panic("没有定义返回值:" + ft.Name + ";" + runSql)
		}
		rts := strings.Replace(rt[1:len(rt)-1], "[]", "array_", -1)
		rts = strings.Replace(rts, ", ", "_", -1)
		rts = strings.Replace(rts, "_bool", "", -1)
		rts = strings.Replace(rts, "_error", "", -1)
		rts = strings.Replace(rts, " {}", "", -1)
		rts = strings.Replace(rts, "[", "_", -1)
		rts = strings.Replace(rts, "]", "_", -1)
		rtArray := strings.Split(rts, "_")
		obj := false

		for _, r := range rtArray {
			if len(r) > 1 {
				if strings.Index(r, ".") > 0 {
					methodName += "Object"
					obj = true
				} else {
					methodName += strings.ToUpper(r[0:1]) + r[1:]
				}
			}
		}
		if obj {
			if strings.Index(ft.Type.String(), "[]") < 0 {
				if strings.Index(runSql, "select * ") == 0 {
					runSql = fmtSelectAllSql(runSql, ft.Type.Out(0))
				}
				fv.Set(reflect.MakeFunc(fv.Type(), func(in []reflect.Value) []reflect.Value {
					defer func() {
						if err := recover(); err != nil {
							seelog.Error("query object fail.methodName=", methodName,
								";\n\trunSql=", runSql,
								"\n\t", err)
							seelog.Flush()
						}
					}()
					v := vti(in)
					return impl.QueryObject(runSql, ft.Type.Out(0), v...)
				}))
			} else {
				if strings.Index(runSql, "select * ") == 0 {
					runSql = fmtSelectAllSql(runSql, ft.Type.Out(0).Elem())
				}
				fv.Set(reflect.MakeFunc(fv.Type(), func(in []reflect.Value) []reflect.Value {
					defer func() {
						if err := recover(); err != nil {
							seelog.Error("query object fail.methodName=", methodName,
								";\n\trunSql=", runSql,
								"\n\t", err)
							seelog.Flush()
						}
					}()
					v := vti(in)
					return impl.QueryObjectArray(runSql, ft.Type.Out(0), v...)
				}))
			}
		} else {
			lowSql := strings.ToLower(runSql)[0:6]
			if lowSql == "insert" {
				methodName = "Insert"
			} else if lowSql == "update" || lowSql == "delete" || lowSql == "replac" {
				methodName = "Exec"
			} else {
				methodName = "Query" + methodName
				methodName = strings.Replace(methodName, "QueryArray", "List", -1)
				//methodName = strings.Replace(methodName, "Bool", "", -1)
			}
			implM, b := implVO.Type().MethodByName(methodName)
			if b {
				fv.Set(reflect.MakeFunc(fv.Type(), func(in []reflect.Value) []reflect.Value {
					params := make([]reflect.Value, len(in)+1)
					defer func() {
						if err := recover(); err != nil {
							seelog.Error(impl.Conn == nil, ";", impl.dsn, ";\n\tmethodName=", methodName,
								";runSql=", runSql, ",", vti(in),
								"\n\t", err)
							seelog.Flush()
						}
					}()
					params[0] = reflect.ValueOf(runSql)
					for idx, pin := range in {
						params[idx+1] = reflect.ValueOf(pin.Interface())
					}
					return implVO.Method(implM.Index).Call(params)
				}))
			} else {
				msg := "方法没有现实:\nfunc (impl *Impl) " + methodName + "(rows sql.Rows, cols []string) " + rt + "{\n\t\n}\n;" + ft.Name + ";sql=" + runSql
				seelog.Error(msg)
				panic(msg)
			}
		}
	}
}
