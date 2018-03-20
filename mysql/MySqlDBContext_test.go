package mysql

import (
	"testing"
	"github.com/devfeel/mapper"
)

var(
	db = NewMySqlDBContext("test:123456@tcp(127.0.0.1:3306)/test?charset=utf8&allowOldPasswords=1")
)

type PullEventLog struct{
	FileID string `mapper:"fileID"`
	MsgHandle string `mapper:"msgHandle"`
	EventType string`mapper:"eventType"`
	Version string`mapper:"version"`
	Data string `mapper:"data"`
}


func TestMySqlDBContext_FindOne(t *testing.T) {
	result, err:=db.FindOne("SELECT * FROM PullEventLog limit 1")
	if err!= nil{
		t.Error(err)
		return
	}
	info := new(PullEventLog)
	mapper.MapperMap(result, info)
	t.Log(info)
}

func TestMySqlDBContext_FindList(t *testing.T) {
	var results []*PullEventLog
	err:=db.FindList(&results, "SELECT * FROM PullEventLog limit 10")
	if err!= nil{
		t.Error(err)
		return
	}else {
		for _, v := range results {
			t.Log(*v)
		}
	}
}

func TestMsSqlDBContext_Insert(t *testing.T) {
	result, err:=db.Insert("INSERT INTO Demo VALUES(?, ?)",888, "insert ")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestMsSqlDBContext_Update(t *testing.T) {
	result, err:=db.Insert("UPDATE Demo set DemoName = ? where DemoID = ?","asdfasf", 1)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestMsSqlDBContext_Delete(t *testing.T) {
	result, err:=db.Delete("Delete Demo where DemoID = ?",888)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}
