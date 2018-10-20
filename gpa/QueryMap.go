package gpa

import (
	"database/sql"
	"github.com/cihub/seelog"
)

func (dao *Gpa) QueryMapStringString(sqlString string, param ...interface{}) (map[string]string, bool, error) {
	rows, err := dao.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			if rows.Next() {
				arr, _ := scan(rows, cols)
				res := make(map[string]string)
				for i := 0; i < len(cols); i++ {
					res[cols[i]] = arr[i].(*sql.NullString).String
				}
				return res, true, nil
			}
		}
		seelog.Warn("not match data?", param)
	} else {
		seelog.Error("QueryMapStringString Fail.", err)
	}
	return nil, false, nil
}
