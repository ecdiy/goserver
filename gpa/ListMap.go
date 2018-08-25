package gpa

import "github.com/cihub/seelog"

func (g *Gpa) ListMapStringString(sql string, param ...interface{}) ([]map[string]string, error) {
	rows, err := g.Conn.Query(sql, param...)
	defer rows.Close()
	if err == nil {
		if cols, ec := rows.Columns(); ec == nil {
			result := make([]map[string]string, 0)
			for rows.Next() {
				result = append(result, rowToMap(rows, cols))
			}
			return result, nil
		}
	} else {
		seelog.Flush()
		seelog.Error("ListMapStringString Error:", sql, err)
	}
	return nil, nil
}
func (g *Gpa) ListMapStringInterface(sql string, param ...interface{}) ([]map[string]interface{}, error) {
	rows, err := g.Conn.Query(sql, param...)
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
