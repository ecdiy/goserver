 
## 
```
go get github.com/ecdiy/gpa 
```

 
## SQL 操作
```
 
 type SqlAction struct {
 	SysRoleDel func(roleId int64, roleId2 int64) (int64, error) `delete from SysRole where id=? and creator!=0 and 0=(SELECT count(*) from SysUserRole where roleId=?)`
 	AllAppId   func() ([]int, error)                            `select id from App  `
 	IntArray2  func(appId int) ([]int, error)                   `select Id,AppEnable from App where Id=?`
 	FindRole2  func() (SysRole, error)                          `select id, createAt  from SysRole where id=3` 	
 }
 
```

##  GPA Object Save Insert
```
type Gpa interface {
	 Save(model interface{}) (int64, error)
	 Insert(s string, param ... interface{}) (int64, error)
	 Exec(s string, param ... interface{}) (int64, error)
}
```

## Demo
```

type SysRole struct {
	Id       int64 `@Id,AutoIncrement`
	RoleName string
	Creator  int64
	Remark   string
	CreateAt time.Time
}

type SqlAction struct {
	SysRoleDel func(roleId int64) (int64, error) `delete from SysRole where id=? `
	FindRole2  func() (SysRole, error)           `select id, createAt  from SysRole where id=3`
}

func Test_Gpa(t *testing.T) {
	defer func() {
		seelog.Flush()
	}()
	
	sqlAction := &SqlAction{}
	
	orm := GetGpa("mysql", "root:root@tcp(127.0.0.1:3306)/base-sys-user?timeout=30s&charset=utf8&parseTime=true",
	 sqlAction)

	sqlAction.SysRoleDel(48)
	sr := &SysRole{RoleName: "TestXX", Creator: 1, CreateAt: time.Now(), Remark: "test"}
	orm.Save(sr)
	fmt.Print(sr.Id)
	row, _ := sqlAction.SysRoleDel(sr.Id)
	if row != 1 {
		t.Error("delete ite fail.")
	}
}

```

#### QQ群: 620063196
