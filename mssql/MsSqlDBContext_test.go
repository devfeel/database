package mssql

import (
	"testing"
	"database/sql"
)

type Demo struct{
	ID int
	DemoID int
	DemoName string
}

var(
	db = NewMsSqlDBContext("server=192.168.8.92;port1433;database=test;user id=sa;password=123456;encrypt=disable")
)

func TestMsSqlDBContext_FindOne(t *testing.T) {
	result := new(Demo)
	err:=db.FindOne(result, "SELECT * FROM [Demo] WHERE DemoID = 3")
	if err!= nil{
		if err == sql.ErrNoRows{
			t.Log(err.Error())
		}else{
			t.Error(err)
		}
	}else{
		t.Log(result)
	}
}

func TestMsSqlDBContext_Count(t *testing.T) {
	count, err:=db.Count("SELECT count(0) FROM [Demo]")
	if err!= nil{
		t.Error(err)
	}else{
		t.Log(count)
	}
}

func TestMsSqlDBContext_QuerySum(t *testing.T) {
	sum, err:=db.Count("SELECT Sum(DemoID) FROM [Demo]")
	if err!= nil{
		t.Error(err)
	}else{
		t.Log(sum)
	}
}


func TestMsSqlDBContext_QueryMax(t *testing.T) {
	max, err:=db.QueryMax("SELECT Max(DemoID) FROM [Demo]")
	if err!= nil{
		t.Error(err)
	}else{
		t.Log(max)
	}
}

func TestMsSqlDBContext_QueryMin(t *testing.T) {
	min, err:=db.QueryMin("SELECT Min(DemoName) FROM [Demo]")
	if err!= nil{
		t.Error(err)
	}else{
		t.Log(min)
	}
}

func TestMsSqlDBContext_FindOneMap(t *testing.T) {
	result, err:=db.FindOneMap("SELECT TOP 10 * FROM [Demo]")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}


func TestMsSqlDBContext_FindMap(t *testing.T) {
	result, err:=db.FindListMap("SELECT TOP 10 * FROM [Demo]")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestQueryByPage(t *testing.T){
	skip := 5
	take := 10
	fields := "*"
	tableName := "Demo"
	where := "DemoID = ?"
	orderBy := "ID ASC, ID DESC"
	var demos []*Demo
	err:=db.FindListByPage(&demos, tableName, fields, where, orderBy, skip, take, 10000)
	if err!= nil{
		t.Error(err)
	}else{
		for _, v:=range demos{
			t.Log(*v)
		}
	}
}

func TestMsSqlDBContext_FindList(t *testing.T) {
	var demos []*Demo
	err := db.FindList(&demos, "SELECT TOP 10 * FROM [Demo]")
	if err!= nil{
		t.Error(err)
	}else{
		for _, v:=range demos{
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

func TestMsSqlDBContext_ExecProc(t *testing.T) {
	result, err:=db.ExecProc("InsertDemo",889,"insert proc")
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(result)
}