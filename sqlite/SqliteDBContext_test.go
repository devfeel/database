package sqlite

import (
	"fmt"
	"github.com/devfeel/database"
	"github.com/devfeel/mapper"
	"testing"
	"time"
)

var (
	db = NewSqliteDBContext("d:/test.db")
)

type DemoInfo struct {
	UID        int64     `mapper:"uid"`
	UserName   string    `mapper:"username"`
	DepartName string    `mapper:"departname"`
	Created    time.Time `mapper:"created"`
}

func TestSqliteDBContext_InitTable(t *testing.T) {
	//创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS Demo(
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NULL,
        departname VARCHAR(64) NULL,
        created DATE NULL
    );
    `
	fmt.Println(db.GetCommand().Exec(sql_table))

}

func TestSqliteDBContext_ShowCountData(t *testing.T) {
	result := new(DemoInfo)
	db.FindOne(result, "SELECT * FROM Demo limit 1")
	db.FindOne(result, "SELECT * FROMDemo limit 1")
	db.Insert("INSERT INTO Demo(username, departname, created) VALUES(?, ?, ?)", "name1", "dev", "2019-01-01")
	db.Update("UPDATE Demo set departname = ? where uid = ?", "test", 1)
	wantItems := 4
	if wantItems == len(database.ShowStateData()) {
		t.Log("DataBase Count Data success", database.ShowStateData())
	} else {
		t.Error("DataBase Count Data failed", database.ShowStateData())
	}
}

func TestSqliteDBContext_FindOne(t *testing.T) {
	result := new(DemoInfo)
	err := db.FindOne(result, "SELECT * FROM Demo limit 1")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

func TestSqliteDBContext_FindOneMap(t *testing.T) {
	result, err := db.FindOneMap("SELECT * FROM Demo limit 1")
	if err != nil {
		t.Error(err)
		return
	}
	info := new(DemoInfo)
	mapper.MapperMap(result, info)
	t.Log(info)
}

func TestSqliteDBContext_FindList(t *testing.T) {
	var results []*DemoInfo
	err := db.FindList(&results, "SELECT * FROM Demo limit 10")
	if err != nil {
		t.Error(err)
		return
	} else {
		for _, v := range results {
			t.Log(*v)
		}
	}
}

func TestSqliteDBContext_Insert(t *testing.T) {
	result, err := db.Insert("INSERT INTO Demo(username, departname, created) VALUES(?, ?, ?)", "name1", "dev", "2019-01-01")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestSqliteDBContext_Update(t *testing.T) {
	result, err := db.Update("UPDATE Demo set departname = ? where uid = ?", "test", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestSqliteDBContext_Delete(t *testing.T) {
	result, err := db.Delete("Delete Demo where uid = ?", 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestSqliteDBContext_FindListByPage(t *testing.T) {
	var results []*DemoInfo
	err := db.FindListByPage(&results, "Demo", "*", "username = ?", "created ASC", 0, 10, "name1")
	fmt.Println(err, len(results))
}

func TestSqliteDBContext_Scalar(t *testing.T) {
	count, err := db.Scalar("SELECT count(0) FROM [Demo]")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
	}
}

func TestSqliteDBContext_Count(t *testing.T) {
	count, err := db.Count("SELECT count(0) FROM [Demo]")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
	}
}

func TestSqliteDBContext_QueryMax(t *testing.T) {
	max, err := db.QueryMax("SELECT Max(UID) FROM [Demo]")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(max)
	}
}

func TestSqliteDBContext_QueryMin(t *testing.T) {
	min, err := db.QueryMin("SELECT Min(username) FROM [Demo]")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(min)
	}
}
