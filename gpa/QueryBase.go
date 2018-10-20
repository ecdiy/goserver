package gpa

import (
	"database/sql"
	"github.com/cihub/seelog"
)

func (dao *Gpa) QueryInt64(sqlString string, param ...interface{}) (int64, bool, error) {
	rows, err := dao.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			var r int64
			rows.Scan(&r)
			return r, true, nil
		}
	} else {
		seelog.Error(dao.dsn, err)
		seelog.Flush()
	}
	return 0, false, nil
}
func (dao *Gpa) QueryInt32(sqlString string, param ...interface{}) (int32, bool, error) {
	rows, err := dao.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			var r int32
			rows.Scan(&r)
			return r, true, nil
		}
	} else {
		seelog.Error(dao.dsn, err)
		seelog.Flush()
	}
	return 0, false, nil
}

func (dao *Gpa) QueryInt(sqlString string, param ...interface{}) (int, bool, error) {
	rows, err := dao.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			var r int
			rows.Scan(&r)
			return r, true, nil
		}
	} else {
		seelog.Error(dao.dsn, err)
		seelog.Flush()
	}
	return 0, false, nil
}

func (dao *Gpa) QueryString(sqlString string, param ...interface{}) (string, bool, error) {
	rows, err := dao.Conn.Query(sqlString, param...)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			var r sql.NullString
			rows.Scan(&r)
			return r.String, true, nil
		}
	} else {
		seelog.Error("QueryString Fail:", dao.dsn, "\n\t", sqlString, param, err)
		seelog.Flush()
	}
	return "", false, nil
}
