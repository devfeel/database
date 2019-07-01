package postgressql

import (
	"testing"
)


var db = NewPostgresSqlDBContext("host=127.0.0.1 port=5432 user=postgres password=123456 dbname=mytest sslmode=disable")
type UserInfo struct {
	ID   string `mapper:"id"`
	Name string `mapper:"name"`
}

func TestPostgresSqlDBContext_Insert(t *testing.T) {
	n, err := db.Insert("INSERT INTO userinfo(id,name) VALUES($1,$2)", 1, "李四")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(n)
	}
}

func TestPostgresSqlDBContext_Update(t *testing.T) {
	n, err := db.Update("update userinfo set name=$1 where id=$2", "王五", 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(n)
	}
}

func TestPostgresSqlDBContext_ExecProc(t *testing.T) {
	
	result, err := db.ExecProc("add", 1, 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
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
	result,err := db.FindOneMap("SELECT * FROM userinfo where id=$1 limit 1", 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestPostgresSqlDBContext_FindList(t *testing.T) {
	var results []*UserInfo
	err := db.FindList(&results, "SELECT * FROM userinfo")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(results)
	}
}


func TestPostgresSqlDBContext_FindListByPage(t *testing.T) {
	var results []*UserInfo
	err := db.FindListByPage(&results, "userinfo", "*", "id=$1", "id ASC", 10, 0,1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(results)
	}
}

func TestPostgresSqlDBContext_Delete(t *testing.T) {
	n, err := db.Update("delete from userinfo where id=$1", 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(n)
	}
}

func TestPostgresSqlDBContext_Scalar(t *testing.T) {
	sum, err := db.Scalar("select sum(id) from userinfo")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(sum)
    }
}
func TestPostgresSqlDBContext_Count(t *testing.T) {
	count, err := db.Scalar("select count(id) from userinfo")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
    }
}
func TestPostgresSqlDBContext_Max(t *testing.T) {
	max, err := db.Scalar("select max(id) from userinfo")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(max)
    }
}
func TestPostgresSqlDBContext_Min(t *testing.T) {
	min, err := db.Scalar("select min(id) from userinfo")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(min)
    }
}