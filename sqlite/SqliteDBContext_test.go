package sqlite

import (
	"testing"
	"github.com/devfeel/mapper"
	"fmt"
	"time"
)

var(
	db = NewSqliteDBContext("d:/test.db")
)

type DemoInfo struct{
	UID int64 `mapper:"uid"`
	UserName string `mapper:"username"`
	DepartName string`mapper:"departname"`
	Created time.Time `mapper:"created"`
}

func TestSqliteDBContext_InitTable(t *testing.T){
	//创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS Demo(
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NULL,
        departname VARCHAR(64) NULL,
        created DATE NULL
    );
    `
	fmt.Println(db.DBCommand.Exec(sql_table))

}



func TestSqliteDBContext_FindOne(t *testing.T) {
	result := new(DemoInfo)
	err:=db.FindOne(result, "SELECT * FROM Demo limit 1")
	if err!= nil{
		t.Error(err)
	}else{
		t.Log(result)
	}
}

func TestSqliteDBContext_FindOneMap(t *testing.T) {
	result, err:=db.FindOneMap("SELECT * FROM Demo limit 1")
	if err!= nil{
		t.Error(err)
		return
	}
	info := new(DemoInfo)
	mapper.MapperMap(result, info)
	t.Log(info)
}

func TestSqliteDBContext_FindList(t *testing.T) {
	var results []*DemoInfo
	err:=db.FindList(&results, "SELECT * FROM Demo limit 10")
	if err!= nil{
		t.Error(err)
		return
	}else {
		for _, v := range results {
			t.Log(*v)
		}
	}
}

func TestSqliteDBContext_Insert(t *testing.T) {
	result, err:=db.Insert("INSERT INTO Demo(username, departname, created) VALUES(?, ?, ?)","name1", "dev", "2019-01-01")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestSqliteDBContext_Update(t *testing.T) {
	result, err:=db.Insert("UPDATE Demo set departname = ? where uid = ?","test", 1)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestSqliteDBContext_Delete(t *testing.T) {
	result, err:=db.Delete("Delete Demo where uid = ?",2)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}
