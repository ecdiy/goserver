package gpa

//
//func (g *Gpa) Sp(sql string, params ...interface{}) (map[string][]map[string]string, error) {
//	data := make(map[string][]map[string]string)
//	//runSql := sql
//	//if strings.Index(sql, "(") < 0 {
//	//	runSql = "call " + runSql
//	//	pSql := ""
//	//	for i := 0; i < len(params); i++ {
//	//		if len(pSql) > 0 {
//	//			pSql += ","
//	//		}
//	//		pSql += "?"
//	//	}
//	//	runSql += "(" + pSql + ")"
//	//}
//	rows, err := g.Conn.Query(sql, params...)
//	defer rows.Close()
//	if err != nil {
//		seelog.Error("调用存储过程出错了.", sql, params, "\n\t", err)
//		return data, err
//	}
//	node := 0
//	for {
//		var list []map[string]string
//		cols, err := rows.Columns()
//		if err != nil {
//			seelog.Error("获取结果集columns出错:", err)
//			return data, err
//		}
//		for rows.Next() {
//			list = append(list, rowToMap(rows, cols))
//		}
//		nodeStr := strconv.Itoa(node)
//		data["r"+nodeStr] = list
//		//data["c"+nodeStr ] = cols
//		node++
//		if !rows.NextResultSet() {
//			break
//		}
//	}
//	return data, nil
//}
