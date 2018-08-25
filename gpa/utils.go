package gpa

import (
	"reflect"
	"database/sql"
)

func fmtSelectAllSql(runSql string, objCls reflect.Type) string {
	//seelog.Info("query object runSql=", objCls.Name(), runSql)
	n := objCls.NumField()
	fields := ""
	for i := 0; i < n; i++ {
		fields += objCls.Field(i).Name + ","
	}
	return "select " + fields[0:len(fields)-1] + runSql[8:]
}

func scan(rows *sql.Rows, cols []string) ([]interface{}, error) {
	arr := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		var inf sql.NullString
		arr[i] = &inf
	}
	return arr, rows.Scan(arr...)
}

func rowToInterface(rows *sql.Rows, cols []string) map[string]interface{} {
	arr, _ := scan(rows, cols)
	res := make(map[string]interface{})
	for i := 0; i < len(cols); i++ {
		v := arr[i].(*sql.NullString)
		if v.Valid {
			res[cols[i]] = arr[i].(*sql.NullString).String
		} else {
			res[cols[i]] = ""
		}
	}
	return res
}

func RowToMap(rows *sql.Rows, cols []string) map[string]string {
	arr, _ := scan(rows, cols)
	res := make(map[string]string)
	for i := 0; i < len(cols); i++ {
		v := arr[i].(*sql.NullString)
		if v.Valid {
			res[cols[i]] = arr[i].(*sql.NullString).String
		} else {
			res[cols[i]] = ""
		}
	}
	return res
}

func vti(in []reflect.Value) []interface{} {
	p := make([]interface{}, len(in))
	for idx, pin := range in {
		p[idx] = pin.Interface()
	}
	return p
}
