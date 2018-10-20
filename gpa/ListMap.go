package gpa

import "github.com/cihub/seelog"

func (dao *Gpa) ListMapStringString(sql string, param ...interface{}) ([]map[string]string, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("Query fail.\n\t", sql , param, "\n", err)
		}
	}()
	rows, err := dao.Conn.Query(sql, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			result := make([]map[string]string, 0)
			for rows.Next() {
				result = append(result, RowToMap(rows, cols))
			}
			return result, nil
		}
	} else {
		seelog.Flush()
		seelog.Error("ListMapStringString Error:", sql, err)
	}
	return nil, nil
}
func (dao *Gpa) ListMapStringInterface(sql string, param ...interface{}) ([]map[string]interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error("Query fail.\n\t", sql , param, "\n", err)
		}
	}()
	rows, err := dao.Conn.Query(sql, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			result := make([]map[string]interface{}, 0)
			for rows.Next() {
				result = append(result, rowToInterface(rows, cols))
			}
			return result, nil
		}
	} else {
		seelog.Flush()
		seelog.Error("ListMapStringString Error:", sql, err)
	}
	return nil, nil
}
