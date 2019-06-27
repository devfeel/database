package mysql

import (
	"errors"
	"github.com/devfeel/database"
	"strconv"
)

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type MySqlDBContext struct {
	DBCommand        *MySqlCommand
	DefaultTableName string
}

func NewMySqlDBContext(connString string) database.DBContext {
	db := new(MySqlDBContext)
	db.Init(connString)
	return db
}

// GetCommand return DBCommand
func (ctx *MySqlDBContext) GetCommand() database.DBCommand {
	return ctx.DBCommand
}

func (ctx *MySqlDBContext) Init(conn string) {
	ctx.DBCommand = new(MySqlCommand)
	ctx.DBCommand.DriverName = DriverName
	ctx.DBCommand.Connection = conn
	ctx.DBCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.DBCommand.PoolIdleConnsCount = Default_IDLE_CONNS
}

// ExecProc executes proc with name
func (ctx *MySqlDBContext) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	return ctx.DBCommand.ExecProc(procName, args...)
}

func (ctx *MySqlDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.LastInsertId()
	return rows, err
}

func (ctx *MySqlDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func (ctx *MySqlDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

// FindOne query data with sql and return dest struct
func (ctx *MySqlDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *MySqlDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
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
func (ctx *MySqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *MySqlDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error) {
	return ctx.DBCommand.Query(sql, args...)
}

// Count query count data with sql, return int64
func (ctx *MySqlDBContext) Count(sql string, args ...interface{}) (count int64, err error) {
	return ctx.DBCommand.QueryCount(sql, args...)
}

// QuerySum query sum data with sql, return int64
func (ctx *MySqlDBContext) QuerySum(sql string, args ...interface{}) (sum int64, err error) {
	return ctx.DBCommand.QueryCount(sql, args...)
}

// QueryMax query max value with sql, return interface{}
func (ctx *MySqlDBContext) QueryMax(sql string, args ...interface{}) (data interface{}, err error) {
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
func (ctx *MySqlDBContext) QueryMin(sql string, args ...interface{}) (data interface{}, err error) {
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
// Sql demo: select * from Demo where DemoID = 1 order by id limit 10,10;
func (ctx *MySqlDBContext) FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{}) error {
	if fields == "" {
		fields = "*"
	}
	if where != "" {
		where = "WHERE " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	sql := "SELECT " + fields + " FROM " + tableName + " " + where + " " + orderBy + " limit " + strconv.Itoa(skip) + "," + strconv.Itoa(take)
	_, err := ctx.DBCommand.Select(dest, sql, args...)
	return err
}
