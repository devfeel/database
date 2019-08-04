package mssql

import (
	"github.com/devfeel/database"
	"github.com/devfeel/database/internal/counter"
	"strconv"
)

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type MsSqlDBContext struct {
	dbCommand        *MsSqlCommand
	passwordCallback database.PasswordCallback
	rawConnString    string
	DefaultTableName string
}

func NewMsSqlDBContext(connString string) database.DBContext {
	db := new(MsSqlDBContext)
	db.Init(connString)
	return db
}

// GetCommand return DBCommand
func (ctx *MsSqlDBContext) GetCommand() database.DBCommand {
	return ctx.dbCommand
}

func (ctx *MsSqlDBContext) Init(conn string) {
	ctx.dbCommand = new(MsSqlCommand)
	ctx.dbCommand.DriverName = DriverName
	ctx.dbCommand.Connection = conn
	ctx.dbCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.dbCommand.PoolIdleConnsCount = Default_IDLE_CONNS

	ctx.rawConnString = conn
	if ctx.passwordCallback != nil {
		ctx.dbCommand.Connection = ctx.passwordCallback(conn)
	}
}

func (ctx *MsSqlDBContext) SetPasswordCallback(callback database.PasswordCallback) {
	ctx.passwordCallback = callback
}

func (ctx *MsSqlDBContext) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	records, err = ctx.dbCommand.ExecProc(procName, args...)
	counter.IncHandler(counter.TOKEN_EXECPROC, err, 1)
	return records, err
}

func (ctx *MsSqlDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_INSERT, err, 1)
		return 0, err
	}
	rows, err := result.LastInsertId()
	counter.IncHandler(counter.TOKEN_INSERT, err, 1)
	return rows, err
}

func (ctx *MsSqlDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_UPDATE, err, 1)
		return 0, err
	}
	rows, err := result.RowsAffected()
	counter.IncHandler(counter.TOKEN_UPDATE, err, 1)
	return rows, err
}

func (ctx *MsSqlDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
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
func (ctx *MsSqlDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *MsSqlDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
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
func (ctx *MsSqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *MsSqlDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error) {
	results, err = ctx.dbCommand.Query(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return results, err
}

// Scalar executes a query that returns first row.
func (ctx *MsSqlDBContext) Scalar(sql string, args ...interface{}) (result interface{}, err error) {
	result, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return result, err
}

// Count query count data with sql, return int64
func (ctx *MsSqlDBContext) Count(sql string, args ...interface{}) (count int64, err error) {
	result, err := ctx.dbCommand.Scalar(sql, args...)
	if err == nil {
		count = result.(int64)
	}
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return count, err
}

// QueryMax query max value with sql, return interface{}
func (ctx *MsSqlDBContext) QueryMax(sql string, args ...interface{}) (data interface{}, err error) {
	data, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return data, err
}

// QueryMin query min value with sql, return interface{}
func (ctx *MsSqlDBContext) QueryMin(sql string, args ...interface{}) (data interface{}, err error) {
	data, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return data, err
}

// FindListByPage query single table data by skip and take
// args: args for where string param
// call demo:
// var results []*Demo
// FindListByPage(&results, "Demo", "*", "DemoID = ?", "ID ASC, DemoName DESC", 10, 10, 10000)
func (ctx *MsSqlDBContext) FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{}) error {
	if fields == "" {
		fields = "*"
	}
	if where != "" {
		where = "WHERE " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	sql := "SELECT * FROM ( "
	sql += "SELECT ROW_NUMBER() OVER (" + orderBy + ") AS [ROW_NUMBER], " + fields + " FROM " + tableName + " AS t0 WITH(NOLOCK) " + where
	sql += ") AS tp WHERE [tp].[ROW_NUMBER] BETWEEN " + strconv.Itoa(skip) + " + 1 AND " + strconv.Itoa(take) + " + " + strconv.Itoa(skip) + " ORDER BY [tp].[ROW_NUMBER]"
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}
