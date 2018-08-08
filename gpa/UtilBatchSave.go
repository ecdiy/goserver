package gpa

/**
批量保存数据,性能最优
*/
//func (g *Gpa) BatchSaveByMap(userId interface{}, ds []interface{}, tableName, pk, cols string) *base.ST {
//	return g.BatchSave(userId, func() string {
//		fs := strings.Split(cols, ",")
//		var sql = ""
//		for _, item := range ds {
//			m := item.(map[string]interface{})
//			sql += `(` + fmt.Sprint(userId)
//			for _, f := range fs {
//				itemV, ib := m[strings.TrimSpace(f)]
//				if ib && itemV != nil {
//					sql += ",'" + fmt.Sprint(itemV) + "'"
//				} else {
//					sql += ",null"
//				}
//			}
//			sql = sql + `),`
//		}
//		return sql
//	}, tableName, pk, cols)
//}
//
//func (g *Gpa) BatchSaveByArray(userId interface{}, ds [][]string, tableName, pk, cols string) *base.ST {
//	return g.BatchSave(userId, func() string {
//		var sql = ""
//		for _, item := range ds {
//			sql += `(` + fmt.Sprint(userId)
//			for i := 0; i < len(item); i++ {
//				itemV := strings.TrimSpace(item[i])
//				sql += ",'" + fmt.Sprint(itemV) + "'"
//			}
//			sql = sql + `),`
//		}
//		return sql
//	}, tableName, pk, cols)
//}
//
//func (g *Gpa) BatchSave(userId interface{}, data func() string, tableName, pk, cols string) *base.ST {
//	chkSql := `select count(*) from ` + tableName + ` where UserId=?`
//	//ds := data.([]interface{})
//	c, _, ef := g.QueryInt(chkSql, userId)
//	if ef != nil {
//		return base.StErrorDb
//	}
//	sql := "INSERT INTO "
//	fs := strings.Split(cols, ",")
//	if c > 0 {
//		tmpSql := `CREATE temporary TABLE tmp` + tableName + `(UserId bigint(20)`
//		for _, f := range fs {
//			tmpSql += "," + f + " varchar(240) "
//		}
//		_, et := g.Exec(tmpSql + ")")
//		if et != nil {
//			return base.StErrorDb
//		}
//		sql += "tmp" + tableName
//	}
//	if c == 0 {
//		sql += tableName
//	}
//	sql += "(UserId," + cols + ")VALUES"
//	sql += data()
//	_, e := g.Exec(sql[0 : len(sql)-1])
//	if e != nil {
//		return base.StErrorDb
//	}
//	if c > 0 {
//		upSql := `update ` + tableName + ` s, tmp` + tableName + ` t set `
//		for _, f := range fs {
//			upSql += "s." + f + "=t." + f + ","
//		}
//		tblCols := "UserId," + cols
//		_, ue := g.Exec(upSql[0:len(upSql)-1] + ` where s.UserId=t.UserId and s.` + pk + `=t.` + pk)
//		if ue != nil {
//			return base.StErrorDb
//		}
//		_, bie := g.Exec(`insert into ` + tableName + ` (` + tblCols + `) select ` + tblCols+
//			` from tmp`+ tableName+ ` where `+ pk+ ` not in (select `+ pk+ ` from `+ tableName+ ` where UserId=?)`, userId)
//		if bie != nil {
//			return base.StErrorDb
//		}
//		g.Exec("drop table tmp" + tableName)
//	}
//	return base.OK
//}
