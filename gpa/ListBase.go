package gpa

import (
	"database/sql"
	"github.com/cihub/seelog"
)

func (me *Gpa) ListInt64(sqlString string, param ...interface{}) ([]int64, error) {
	rows, err := me.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			colLen := len(cols)
			if colLen == 1 {
				var list []int64
				for rows.Next() {
					var r sql.NullInt64
					rows.Scan(&r)
					if r.Valid {
						list = append(list, r.Int64)
					} else {
						list = append(list, 0)
					}
				}
				return list, nil
			} else {
				for rows.Next() {
					arr := make([]interface{}, colLen)
					for i := 0; i < colLen; i++ {
						var ix sql.NullInt64
						arr[i] = &ix
					}
					es := rows.Scan(arr...)
					if es != nil {
						seelog.Info("scan int64 fail. ", es)
						return nil, es
					} else {
						var res = make([]int64, colLen)
						for i := 0; i < colLen; i++ {
							ix := arr[i].(*sql.NullInt64)
							if ix.Valid {
								res[i] = ix.Int64
							}
						}
						return res, nil
					}
				}
			}
		}
	}
	return nil, nil
}

func (me *Gpa) ListString(sqlString string, param ...interface{}) ([]string, error) {
	rows, err := me.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			colLen := len(cols)
			if colLen == 1 {
				var list []string
				for rows.Next() {
					var r sql.NullString
					rows.Scan(&r)
					list = append(list, r.String)
				}
				return list, nil
			} else {
				if rows.Next() {
					arr, _ := scan(rows, cols)
					resArray := make([]string, colLen)
					for i := 0; i < len(cols); i++ {
						resArray[i] = arr[i].(*sql.NullString).String
					}
					return resArray, nil
				}
			}
		}
	} else {
		seelog.Error("ListString error,", sqlString, param, err)
	}
	return nil, nil
}
