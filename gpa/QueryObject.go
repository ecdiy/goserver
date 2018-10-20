package gpa

import (
	"database/sql"
	"github.com/cihub/seelog"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/**
这个方法可能不是最有效的办法， 如果用   v.FieldByName(tf.Name).Addr().Interface()
会有nil无法转换的问题
用sql.NullString==> value方式能解决nil问题
*/

func (dao *Gpa) QueryObjectArray(runSql string, resultType reflect.Type, param ...interface{}) []reflect.Value {

	rows, _ := dao.Conn.Query(runSql, param...)
	defer rows.Close()
	var res [][]interface{}
	cols, _ := rows.Columns()
	numC := len(cols)
	for rows.Next() {
		oneRow := make([]interface{}, numC)
		for i := 0; i < numC; i++ {
			var inf sql.NullString
			oneRow[i] = &inf
		}
		err := rows.Scan(oneRow...)
		if err == nil {
			res = append(res, oneRow)
		}
	}
	slice := reflect.MakeSlice(resultType, len(res), len(res))
	ele := resultType.Elem()
	for idx, v := range res {
		slice.Index(idx).Set(objectConvert(ele, v, cols))
	}
	return []reflect.Value{slice, nilVf}
}

func (dao *Gpa) QueryObject(runSql string, resultType reflect.Type, param ...interface{}) []reflect.Value {
	rows, err := dao.Conn.Query(runSql, param...)
	defer rows.Close()
	rv := reflect.New(resultType).Elem()
	if err == nil {
		if rows.Next() {
			cols, _ := rows.Columns()
			numF := resultType.NumField()
			v := make([]interface{}, len(cols))
			for i := 0; i < numF; i++ {
				var inf sql.NullString
				v[i] = &inf
			}
			err := rows.Scan(v...)
			if err != nil {
				seelog.Error("对象转换出错:", "\n\t", err)
				return []reflect.Value{rv, reflect.ValueOf(false), reflect.ValueOf(err)}
			}
			return []reflect.Value{objectConvert(resultType, v, cols), reflect.ValueOf(true), nilVf}
		}
	}
	return []reflect.Value{rv, reflect.ValueOf(false), nilVf}
}

func objectConvert(ele reflect.Type, v []interface{}, cols []string) reflect.Value {
	vNew := reflect.New(ele).Elem()
	numF := ele.NumField()
	for i, item := range v {
		ns := item.(*sql.NullString)
		if !ns.Valid {
			continue
		}
		cLow := strings.ToLower(cols[i])
		for i2 := 0; i2 < numF; i2++ {
			tf := ele.Field(i2)
			if cLow == strings.ToLower(strings.ToLower(tf.Name)) {
				conV := stringToInterface(tf.Type, ns.String)
				if conV != nilVf {
					vNew.Field(i2).Set(conV)
				}
				break
			}
		}
	}
	return vNew
}

func stringToInterface(sf reflect.Type, str string) reflect.Value {
	switch sf.Name() {
	case "int64":
		v, _ := strconv.ParseInt(str, 10, 0)
		return reflect.ValueOf(v)
	case "int":
		v, _ := strconv.Atoi(str)
		return reflect.ValueOf(v)
	case "int32":
		v, _ := strconv.Atoi(str)
		return reflect.ValueOf(int32(v))
	case "string":
		return reflect.ValueOf(str)
	case "Time":
		d, de := time.Parse("2006-01-02T15:04:05Z", str)
		if de == nil {
			return reflect.ValueOf(d)
		} else {
			seelog.Error("日期格式转化错误:", str)
		}
	}
	seelog.Error("--TODO--un impl type----", sf.Name(), str)
	return nilVf
}
