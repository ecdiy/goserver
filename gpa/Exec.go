package gpa

import (
	"github.com/cihub/seelog"
	"reflect"
	"strings"
)

/**
返回影响记录的行数
*/

func (me *Gpa) Exec(runSql string, p ...interface{}) (int64, error) {
	row, er := me.conn.Exec(runSql, p...)
	if er == nil {
		ra, _ := row.RowsAffected()
		return ra, nil
	} else {
		seelog.Error("SQL执行失败:\n\tsql=", runSql, "\n\tparam=", p, "\n\terror=", er)
		return -1, er
	}
}

/**
返回 自增ID
*/
func (me *Gpa) Insert(s string, param ...interface{}) (int64, error) {
	row, err := me.conn.Exec(s, param...)
	if err == nil {
		return row.LastInsertId()
	} else {
		seelog.Error("insert对象失败,", s, param, err)
		return -1, err
	}
}
func (me *Gpa) Save(model interface{}) (int64, error) {
	toe := reflect.TypeOf(model).Elem()
	voe := reflect.ValueOf(model).Elem()
	e, err := me.exist(toe, voe)
	if err != nil {
		return -1, err
	} else {
		if e == 1 {
			return me.update(toe, voe)
		} else {
			return me.insert(toe, voe)
		}
	}
}

func (me *Gpa) InsertObject(model interface{}) (int64, error) {
	toe := reflect.TypeOf(model).Elem()
	voe := reflect.ValueOf(model).Elem()
	return me.insert(toe, voe)
}

func (me *Gpa) InsertMap(table string, val map[string]string) (int64, error) {
	sql, fix := "insert into "+table+"(", ""
	var va []interface{}
	for k, v := range val {
		sql += k + ","
		fix += "?,"
		va = append(va, v)
	}
	sql = sql[0:len(sql)-1] + ")values(" + fix[0:len(fix)-1] + ")"
	return me.Insert(sql, va...)
}

func (me *Gpa) exist(toe reflect.Type, voe reflect.Value) (int64, error) {
	n := toe.NumField()
	s := "select 1 from " + toe.Name() + " where "
	var param []interface{}
	for i := 0; i < n; i++ {
		fx := toe.Field(i)
		tag := string(fx.Tag)
		col := strings.ToLower(fx.Name[0:1]) + fx.Name[1:]
		vF := voe.FieldByName(fx.Name).String()
		if len(vF) >= 1 {
			if strings.Index(tag, PrimaryId) >= 0 {
				s += col + "=? and "
				param = append(param, vF)
			}
		}
	}
	s = s[0 : len(s)-4]
	//fmt.Println(s)
	rows, err := me.conn.Query(s, param...)
	defer rows.Close()
	if err != nil {
		seelog.Error("SQL出错", s, ":", err)
		return -1, err
	}
	if rows.Next() {
		var res int64
		rows.Scan(&res)
		return res, nil
	} else {
		return 0, nil
	}
}

func (me *Gpa) insert(toe reflect.Type, voe reflect.Value) (int64, error) {
	n := toe.NumField()
	s := "insert into " + toe.Name() + "("
	cols := ""
	var param []interface{}
	auto := ""
	for i := 0; i < n; i++ {
		fx := toe.Field(i)
		tag := string(fx.Tag)
		vF := voe.FieldByName(fx.Name).String()
		if len(vF) == 0 {
			continue
		}
		c := strings.ToLower(fx.Name[0:1]) + fx.Name[1:]
		if strings.Index(tag, AutoIncrement) < 0 {
			s += c + ","
			cols += "?,"
			param = append(param, voe.FieldByName(fx.Name).Interface())
		} else {
			auto = fx.Name
		}
	}

	s = s[0:len(s)-1] + ")values(" + cols[0:len(cols)-1] + ")"
	row, err := me.conn.Exec(s, param...)
	if err == nil {
		rii, _ := row.LastInsertId()
		if rii > 0 && len(auto) > 0 {
			voe.FieldByName(auto).SetInt(rii)
		} else {
			rii, _ = row.RowsAffected()
		}
		return rii, err
	} else {
		seelog.Error("insert对象失败,", s, param, err)
		return -1, err
	}
}

func (me *Gpa) update(toe reflect.Type, voe reflect.Value) (int64, error) {
	n := toe.NumField()
	s := "update " + toe.Name() + " set "
	pri := ""
	var param, priParam []interface{}
	for i := 0; i < n; i++ {
		fx := toe.Field(i)
		tag := string(fx.Tag)
		col := strings.ToLower(fx.Name[0:1]) + fx.Name[1:]
		vF := voe.FieldByName(fx.Name).String()
		if len(vF) >= 1 {
			if strings.Index(tag, PrimaryId) >= 0 || strings.Index(tag, AutoIncrement) >= 0 {
				pri += col + "=? and "
				priParam = append(priParam, vF)
			} else {
				s += col + "=?,"
				param = append(param, vF)
			}
		}
	}
	s = s[0:len(s)-1] + " where " + pri[0:len(pri)-4]
	//fmt.Println(s)
	for _, v := range priParam {
		param = append(param, v)
	}
	row, er := me.conn.Exec(s, param...)
	if er == nil {
		raf, _ := row.RowsAffected()
		return raf, er
	} else {
		seelog.Error("update对象失败,", er)
		return -1, er
	}
}
