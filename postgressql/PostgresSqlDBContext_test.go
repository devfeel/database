package postgressql

import (
	"testing"
)

var db = NewPostgresSqlDBContext("host=127.0.0.1 port=5432 user=postgres password=123456 dbname=mytest sslmode=disable")

type UserInfo struct {
	ID   string `mapper:"id"`
	Name string `mapper:"name"`
}

func TestPostgresSqlDBContext_ExecProc(t *testing.T) {

}

func TestPostgresSqlDBContext_FindOne(t *testing.T) {
	result := new(UserInfo)
	err := db.FindOne(result, "SELECT * FROM userinfo limit 1")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestPostgresSqlDBContext_FindOneMap(t *testing.T) {

}

func TestPostgresSqlDBContext_FindList(t *testing.T) {

}

func TestPostgresSqlDBContext_Insert(t *testing.T) {

}

func TestPostgresSqlDBContext_Update(t *testing.T) {

}

func TestPostgresSqlDBContext_Delete(t *testing.T) {

}

func TestPostgresSqlDBContext_FindListByPage(t *testing.T) {

}
