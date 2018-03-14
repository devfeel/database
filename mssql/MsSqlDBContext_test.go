package mssql

import (
	"testing"
)

var(
	db = NewMsSqlDBContext("server=127.0.0.1;port1433;database=test;user id=sa;password=123456;encrypt=disable")
)

func TestMsSqlDBContext_FindOne(t *testing.T) {
	result, err:=db.FindOne("SELECT TOP 10 * FROM [Demo]")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestMsSqlDBContext_FindList(t *testing.T) {
	result, err:=db.FindList("SELECT TOP 10 * FROM [Demo]")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
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

func TestMsSqlDBContext_ExecProc(t *testing.T) {
	result, err:=db.ExecProc("InsertDemo",889,"insert proc")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}