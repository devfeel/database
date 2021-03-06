package sqlite

import (
	"fmt"
	"github.com/devfeel/database"
	"github.com/devfeel/database/internal/counter"
	"strconv"
)

const (
	Default_OPEN_CONNS    = 50
	Default_IDLE_CONNS    = 50
	Default_FieldName_MAX = "MAX"
	Default_FieldName_MIN = "MIN"
)

type SqliteDBContext struct {
	dbCommand        *SqliteCommand
	passwordCallback database.PasswordCallback
	rawConnString    string
	DefaultTableName string
}

func NewSqliteDBContext(connString string) database.DBContext {
	db := new(SqliteDBContext)
	db.Init(connString)
	return db
}

// GetCommand return DBCommand
func (ctx *SqliteDBContext) GetCommand() database.DBCommand {
	return ctx.dbCommand
}

func (ctx *SqliteDBContext) Init(conn string) {
	ctx.dbCommand = new(SqliteCommand)
	ctx.dbCommand.DriverName = DriverName
	ctx.dbCommand.Connection = conn
	ctx.dbCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.dbCommand.PoolIdleConnsCount = Default_IDLE_CONNS

	ctx.rawConnString = conn
	if ctx.passwordCallback != nil {
		ctx.dbCommand.Connection = ctx.passwordCallback(conn)
	}
}

func (ctx *SqliteDBContext) SetPasswordCallback(callback database.PasswordCallback) {
	ctx.passwordCallback = callback
}

// ExecProc executes proc with name
func (ctx *SqliteDBContext) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	records, err = ctx.dbCommand.ExecProc(procName, args...)
	counter.IncHandler(counter.TOKEN_EXECPROC, err, 1)
	return records, err
}

func (ctx *SqliteDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_INSERT, err, 1)
		return 0, err
	}
	rows, err := result.RowsAffected()
	counter.IncHandler(counter.TOKEN_INSERT, err, 1)
	return rows, err
}

func (ctx *SqliteDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_UPDATE, err, 1)
		return 0, err
	}
	rows, err := result.RowsAffected()
	counter.IncHandler(counter.TOKEN_UPDATE, err, 1)
	return rows, err
}

func (ctx *SqliteDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_DELETE, err, 1)
		return 0, err
	}
	rows, err := result.RowsAffected()
	counter.IncHandler(counter.TOKEN_DELETE, err, 1)
	return rows, err
}

// FindOne query data with sql and return dest struct
func (ctx *SqliteDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *SqliteDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
	results, err := ctx.dbCommand.Query(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
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
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *SqliteDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error) {
	results, err = ctx.dbCommand.Query(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return results, err
}

// Scalar executes a query that returns first row.
func (ctx *SqliteDBContext) Scalar(sql string, args ...interface{}) (result interface{}, err error) {
	result, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return result, err
}

// Count query count data with sql, return int64
func (ctx *SqliteDBContext) Count(sql string, args ...interface{}) (count int64, err error) {
	result, err := ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	if err == nil {
		count = result.(int64)
	}
	return count, err
}

// QueryMax query max value with sql, return interface{}
func (ctx *SqliteDBContext) QueryMax(sql string, args ...interface{}) (data interface{}, err error) {
	data, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return data, err
}

// QueryMin query min value with sql, return interface{}
func (ctx *SqliteDBContext) QueryMin(sql string, args ...interface{}) (data interface{}, err error) {
	data, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return data, err
}

// FindListByPage query single table data by skip and take
// args: args for where string param
// call demo:
// var results []*Demo
// FindListByPage(&results, "Demo", "*", "DemoID = ?", "ID ASC, DemoName DESC", 10, 10, 10000)
// Sql demo: select * from Demo where DemoID = 1 order by id limit 10 offset 10;
func (ctx *SqliteDBContext) FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{}) error {
	if fields == "" {
		fields = "*"
	}
	if where != "" {
		where = "WHERE " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}

	skipStr := ""
	if skip > 0 {
		skipStr = " offset " + strconv.Itoa(skip)
	}

	sql := "SELECT " + fields + " FROM " + tableName + " " + where + " " + orderBy + " limit " + strconv.Itoa(take) + skipStr
	fmt.Println(sql)
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}
