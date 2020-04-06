package mysql

import (
	"github.com/devfeel/database"
	"github.com/devfeel/database/internal/convert"
	"github.com/devfeel/database/internal/counter"
	"strconv"
)

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type MySqlDBContext struct {
	dbCommand        *MySqlCommand
	passwordCallback database.PasswordCallback
	rawConnString    string
	DefaultTableName string
}

func NewMySqlDBContext(connString string) database.DBContext {
	db := new(MySqlDBContext)
	db.Init(connString)
	return db
}

// GetCommand return DBCommand
func (ctx *MySqlDBContext) GetCommand() database.DBCommand {
	return ctx.dbCommand
}

func (ctx *MySqlDBContext) Init(conn string) {
	ctx.dbCommand = new(MySqlCommand)
	ctx.dbCommand.DriverName = DriverName
	ctx.dbCommand.Connection = conn
	ctx.dbCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.dbCommand.PoolIdleConnsCount = Default_IDLE_CONNS

	ctx.rawConnString = conn
	if ctx.passwordCallback != nil {
		ctx.dbCommand.Connection = ctx.passwordCallback(conn)
	}
}

func (ctx *MySqlDBContext) SetPasswordCallback(callback database.PasswordCallback) {
	ctx.passwordCallback = callback
}

// ExecProc executes proc with name
func (ctx *MySqlDBContext) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	records, err = ctx.dbCommand.ExecProc(procName, args...)
	counter.IncHandler(counter.TOKEN_EXECPROC, err, 1)
	return records, err
}

func (ctx *MySqlDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_INSERT, err, 1)
		return 0, err
	}
	rows, err := result.RowsAffected()
	counter.IncHandler(counter.TOKEN_INSERT, err, 1)
	return rows, err
}

func (ctx *MySqlDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.dbCommand.Exec(sql, args...)
	if err != nil {
		counter.IncHandler(counter.TOKEN_UPDATE, err, 1)
		return 0, err
	}
	rows, err := result.RowsAffected()
	counter.IncHandler(counter.TOKEN_UPDATE, err, 1)
	return rows, err
}

func (ctx *MySqlDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
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
func (ctx *MySqlDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *MySqlDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
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
func (ctx *MySqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *MySqlDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error) {
	results, err = ctx.dbCommand.Query(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return results, err
}

// Scalar executes a query that returns first row.
func (ctx *MySqlDBContext) Scalar(sql string, args ...interface{}) (result interface{}, err error) {
	result, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return result, err
}

// Count query count data with sql, return int64
func (ctx *MySqlDBContext) Count(sql string, args ...interface{}) (count int64, err error) {
	result, err := ctx.dbCommand.Scalar(sql, args...)
	if err == nil {
		count, err = convert.UInt8SliceToInt64(result.([]uint8))
	}
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return count, err
}

// QueryMax query max value with sql, return interface{}
func (ctx *MySqlDBContext) QueryMax(sql string, args ...interface{}) (data interface{}, err error) {
	data, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return data, err
}

// QueryMin query min value with sql, return interface{}
func (ctx *MySqlDBContext) QueryMin(sql string, args ...interface{}) (data interface{}, err error) {
	data, err = ctx.dbCommand.Scalar(sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
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
	_, err := ctx.dbCommand.Select(dest, sql, args...)
	counter.IncHandler(counter.TOKEN_SELECT, err, 1)
	return err
}
