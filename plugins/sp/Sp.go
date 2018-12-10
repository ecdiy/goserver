package sp

import (
	"github.com/cihub/seelog"
	"github.com/ecdiy/goserver/gpa"
	"database/sql"
	"github.com/ecdiy/goserver/utils"
)

//--sp
type Sp struct {
	Sql, Name, SessionName string
	Params                 []*SpParam
	Result                 []*SpResult
	//Auth                   func(c *Param) *UserBase
}

type SpResult struct {
	Name, Type string //[ [total,object],[list,list] ]
}
type SpParam struct {
	ParamName, DefaultVal string
	ValFunc               ParamValFunc
}

type ParamValFunc func(ctx *utils.Param, p *SpParam) (interface{}, int)

func (sp *Sp) Run(data map[string]interface{}, Conn *sql.DB, params ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("调用存储过程出错了.", err)
		}
	}()
	rows, err := Conn.Query(sp.Sql, params...)
	defer rows.Close()
	if err != nil {
		seelog.Error(err, "\n\t", FmtSql(sp.Sql, params...))
		return err
	}
	for node := 0; node < len(sp.Result); node++ {
		cols, err := rows.Columns()
		if err != nil {
			seelog.Error("获取结果集columns出错:", err)
			return err
		}
		r := sp.Result[node]
		dbRes, resultExt := getResult(r.Type, cols, rows)
		if resultExt {
			data[sp.Result[node].Name] = dbRes
		} else {
			//bf := plugins.GetRefByName(r.Type).(plugins.BaseFun)
			//bf(param)
		}
		if !rows.NextResultSet() {
			break
		}
	}
	return nil
}
func (sp *Sp) GetInt64(Conn *sql.DB, params ...interface{}) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("sp query int64 error;", sp.Sql, err)
		}
	}()
	rows, err := Conn.Query(sp.Sql, params...)
	defer rows.Close()
	if err != nil {
		seelog.Error("调用存储过程出错了.", sp.Sql, params, "\n\t", err)
		return 0, err
	}
	if rows.Next() {
		var r sql.NullInt64
		rows.Scan(&r)
		return r.Int64, nil
	}
	return 0, nil
}

func getResult(t string, cols []string, rows *sql.Rows) (interface{}, bool) {
	if t == "list" || t == "l" {
		var list []map[string]string
		for rows.Next() {
			list = append(list, gpa.RowToMap(rows, cols))
		}
		return list, true
	}
	if t == "object" || t == "o" || t == "map" {
		if rows.Next() {
			return gpa.RowToMap(rows, cols), true
		} else {
			return make(map[string]string), true
		}
	}
	if t == "int" {
		if rows.Next() {
			var r sql.NullInt64
			rows.Scan(&r)
			return r.Int64, true
		} else {
			return 0, true
		}
	}
	if t == "string" {
		if rows.Next() {
			var r sql.NullString
			rows.Scan(&r)
			return r.String, true
		} else {
			return 0, true
		}
	}
	seelog.Error("未知类型：", t)
	return nil, false
}
