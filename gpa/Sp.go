package gpa

import (
	"strings"
	"github.com/cihub/seelog"
	"strconv"
)



func (g *Gpa) Sp(sp string, params ... interface{}) (map[string][]map[string]string, error) {
	data := make(map[string][]map[string]string)
	runSql := sp
	if strings.Index(sp, "(") < 0 {
		runSql = "call " + runSql
		pSql := ""
		for i := 0; i < len(params); i++ {
			if len(pSql) > 0 {
				pSql += ","
			}
			pSql += "?"
		}
		runSql += "(" + pSql + ")"
	}
	rows, err := g.conn.Query(runSql, params...)
	defer rows.Close()
	if err != nil {
		seelog.Error("调用存储过程出错了.", runSql, params, "\n\t", err)
		return data, err
	}
	node := 0
	for {
		var list []map[string]string
		cols, err := rows.Columns()
		if err != nil {
			seelog.Error("获取结果集columns出错:", err)
			return data, err
		}
		for rows.Next() {
			list = append(list, rowToMap(rows, cols))
		}
		nodeStr := strconv.Itoa(node)
		data["r"+nodeStr ] = list
		//data["c"+nodeStr ] = cols
		node++
		if !rows.NextResultSet() {
			break
		}
	}
	return data, nil
}
