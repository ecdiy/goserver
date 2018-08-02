package gpa

func (me *Gpa) Get(key string) (string, bool, error) {
	return me.QueryString(`select Val from KeyVal where K=?`, key)
}

func (me *Gpa) Set(k string, v interface{}) (int64, error) {
	_, b, _ := me.Get(k)
	if b {
		return me.Exec(`update KeyVal set Val=? where K=?`, v, k)
	} else {
		return me.Exec("insert into KeyVal(K,Val)values(?,?)", k, v)
	}
}
