package gpa

import (
	"database/sql"
	"github.com/cihub/seelog"
)

func (dao *Gpa) ListArrayString(sqlString string, param ...interface{}) ([][]string, error) {
	rows, err := dao.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			var result [][]string
			for rows.Next() {
				result = append(result, rowToStringArray(rows, cols))
			}
			return result, nil
		}
	} else {
		seelog.Error("数据库查询出错:", dao.dsn, err)
		return nil, err
	}
	return nil, nil
}
func rowToStringArray(rows *sql.Rows, cols []string) []string {
	arr, _ := scan(rows, cols)
	res := make([]string, len(cols))
	for i := 0; i < len(cols); i++ {
		v := arr[i].(*sql.NullString)
		if v.Valid {
			res[i] = arr[i].(*sql.NullString).String
		} else {
			res[i] = ""
		}
	}
	return res
}
