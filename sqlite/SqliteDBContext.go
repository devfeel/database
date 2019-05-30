package sqlite

import (
	"errors"
	"database/sql"
)

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type SqliteDBContext struct {
	DBCommand        *SqliteCommand
	DefaultTableName string
}

func NewSqliteDBContext(connString string) *SqliteDBContext {
	db := new(SqliteDBContext)
	db.Init(connString)
	return db
}

func (ctx *SqliteDBContext) Init(conn string){
	ctx.DBCommand = new(SqliteCommand)
	ctx.DBCommand.DriverName = DriverName
	ctx.DBCommand.Connection = conn
	ctx.DBCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.DBCommand.PoolIdleConnsCount = Default_IDLE_CONNS
}

// Exec executes sql
func (ctx *SqliteDBContext) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	return ctx.DBCommand.Exec(sql, args...)
}

func (ctx *SqliteDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.LastInsertId()
	return rows, err
}

func (ctx *SqliteDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func (ctx *SqliteDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

// FindOne query data with sql and return dest struct
func (ctx *SqliteDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err :=ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *SqliteDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
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
func (ctx *SqliteDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	_, err:= ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *SqliteDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error){
	return ctx.DBCommand.Query(sql, args...)
}

// Count query count data with sql, return int64
func (ctx *SqliteDBContext) Count(sql string, args ...interface{})(count int64, err error) {
	return ctx.DBCommand.QueryCount(sql, args...)
}

// QuerySum query sum data with sql, return int64
func (ctx *SqliteDBContext) QuerySum(sql string, args ...interface{})(sum int64, err error) {
	return ctx.DBCommand.QueryCount(sql, args...)
}

// QueryMax query max value with sql, return interface{}
func (ctx *SqliteDBContext) QueryMax(sql string, args ...interface{})(data interface{}, err error) {
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
func (ctx *SqliteDBContext) QueryMin(sql string, args ...interface{})(data interface{}, err error) {
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
