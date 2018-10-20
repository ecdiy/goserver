package gpa

func (dao *Gpa) Get(key string) (string, bool, error) {
	return dao.QueryString(`select Val from KeyVal where K=?`, key)
}

func (dao *Gpa) Set(k string, v interface{}) (int64, error) {
	_, b, _ := dao.Get(k)
	if b {
		return dao.Exec(`update KeyVal set Val=? where K=?`, v, k)
	} else {
		return dao.Exec("insert into KeyVal(K,Val)values(?,?)", k, v)
	}
}
