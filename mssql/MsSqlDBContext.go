package mssql

import (
	"strconv"
	"errors"
)

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type MsSqlDBContext struct {
	DBCommand        *MsSqlCommand
	DefaultTableName string
}

func NewMsSqlDBContext(connString string) *MsSqlDBContext {
	db := new(MsSqlDBContext)
	db.Init(connString)
	return db
}

func (ctx *MsSqlDBContext) Init(conn string){
	ctx.DBCommand = new(MsSqlCommand)
	ctx.DBCommand.DriverName = DriverName
	ctx.DBCommand.Connection = conn
	ctx.DBCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.DBCommand.PoolIdleConnsCount = Default_IDLE_CONNS
}

func (ctx *MsSqlDBContext) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	return ctx.DBCommand.ExecProc(procName, args...)
}

func (ctx *MsSqlDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.LastInsertId()
	return rows, err
}

func (ctx *MsSqlDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func (ctx *MsSqlDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

// FindOne query data with sql and return dest struct
func (ctx *MsSqlDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err:=ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *MsSqlDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
	results, err := ctx.DBCommand.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(results) <= 0 {
		return make(map[string]interface{}), nil
	}
	return results[0], nil
}

// FindList query data with sql and return dest struct slice
// slice's elem type must ptr
func (ctx *MsSqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	_, err:= ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *MsSqlDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error){
	return ctx.DBCommand.Query(sql, args...)
}

// Count query count data with sql, return int64
func (ctx *MsSqlDBContext) Count(sql string, args ...interface{})(count int64, err error) {
	result, err := ctx.DBCommand.Query(sql, args...)
	if err != nil {
		return 0, err
	}
	if result == nil || len(result) == 0 {
		return 0, errors.New("no data return")
	}
	count = result[0][""].(int64)
	return count, err
}

// QuerySum query sum data with sql, return int64
func (ctx *MsSqlDBContext) QuerySum(sql string, args ...interface{})(sum int64, err error) {
	result, err := ctx.DBCommand.Query(sql, args...)
	if err != nil {
		return 0, err
	}
	if result == nil || len(result) == 0 {
		return 0, errors.New("no data return")
	}
	sum = result[0][""].(int64)
	return sum, err
}

// QueryMax query max value with sql, return interface{}
func (ctx *MsSqlDBContext) QueryMax(sql string, args ...interface{})(data interface{}, err error) {
	result, err := ctx.DBCommand.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		return nil, errors.New("no data return")
	}
	data = result[0][""]
	return data, err
}

// QueryMin query min value with sql, return interface{}
func (ctx *MsSqlDBContext) QueryMin(sql string, args ...interface{})(data interface{}, err error) {
	result, err := ctx.DBCommand.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		return nil, errors.New("no data return")
	}
	data = result[0][""]
	return data, err
}



// FindListByPage query single table data by skip and take
// args: args for where string param
// call demo:
// var results []*Demo
// FindListByPage(&results, "Demo", "*", "DemoID = ?", "ID ASC, DemoName DESC", 10, 10, 10000)
func (ctx *MsSqlDBContext) FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{})  error {
	if fields == ""{
		fields = "*"
	}
	if where != ""{
		where = "WHERE " + where
	}
	if orderBy != ""{
		orderBy = "ORDER BY " + orderBy
	}
	sql := "SELECT * FROM ( "
	sql += "SELECT ROW_NUMBER() OVER ("+orderBy+") AS [ROW_NUMBER], "+fields+" FROM " +tableName+" AS t0 WITH(NOLOCK) " + where
	sql += ") AS tp WHERE [tp].[ROW_NUMBER] BETWEEN "+ strconv.Itoa(skip)+" + 1 AND "+ strconv.Itoa(take)+" + "+ strconv.Itoa(skip)+" ORDER BY [tp].[ROW_NUMBER]"
	_, err:=ctx.DBCommand.Select(dest, sql, args...)
	return err
}