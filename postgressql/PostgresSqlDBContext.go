package postgressql

import (
	"errors"
	"strconv"
	"github.com/devfeel/database"
)

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type PostgresSqlDBContext struct {
	DBCommand        *PostgressCommand
	DefaultTableName string
}

func NewPostgresSqlDBContext(connString string) database.DBContext {

	db := new(PostgresSqlDBContext)
	db.Init(connString)
	return db
}

// GetCommand return DBCommand
func (ctx *PostgresSqlDBContext) GetCommand() database.DBCommand {
	return ctx.DBCommand
}

func (ctx *PostgresSqlDBContext) Init(conn string) {
	ctx.DBCommand = new(PostgressCommand)
	ctx.DBCommand.DriverName = DriverName
	ctx.DBCommand.Connection = conn
	ctx.DBCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.DBCommand.PoolIdleConnsCount = Default_IDLE_CONNS
}

// ExecProc executes proc with name
func (ctx *PostgresSqlDBContext) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	return ctx.DBCommand.ExecProc(procName, args...)
}

func (ctx *PostgresSqlDBContext) Insert(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func (ctx *PostgresSqlDBContext) Update(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func (ctx *PostgresSqlDBContext) Delete(sql string, args ...interface{}) (n int64, err error) {
	result, err := ctx.DBCommand.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

// FindOne query data with sql and return dest struct
func (ctx *PostgresSqlDBContext) FindOne(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindOneMap query data with sql and return map[string]interface{}
func (ctx *PostgresSqlDBContext) FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error) {
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
func (ctx *PostgresSqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	_, err := ctx.DBCommand.Select(dest, sql, args...)
	return err
}

// FindListMap query data with sql and return []map[string]interface{}
func (ctx *PostgresSqlDBContext) FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error) {
	return ctx.DBCommand.Query(sql, args...)
}

// FindListByPage query single table data by skip and take
// args: args for where string param
// call demo:
// var results []*Demo
// FindListByPage(&results, "Demo", "*", "DemoID = ?", "ID ASC, DemoName DESC", 10, 10, 10000)
// Sql demo: select * from Demo where DemoID = 1 order by id limit 10,10;
func (ctx *PostgresSqlDBContext) FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{}) error {
	if fields == "" {
		fields = "*"
	}
	if where != "" {
		where = "WHERE " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	sql := "SELECT " + fields + " FROM " + tableName + " " + where + " " + orderBy + " limit " + strconv.Itoa(take) + " OFFSET " + strconv.Itoa(skip)
	_, err := ctx.DBCommand.Select(dest, sql, args...)
	return err
}

func (ctx *PostgresSqlDBContext) Scalar(sql string,args ...interface{})(result interface{},err error){
	return ctx.DBCommand.Scalar(sql, args...)
}

// Count query count data with sql, return int64
func (ctx *PostgresSqlDBContext) Count(sql string, args ...interface{}) (count int64, err error) {
	return ctx.DBCommand.QueryCount(sql, args...)
}

// QuerySum query sum data with sql, return int64
func (ctx *PostgresSqlDBContext) QuerySum(sql string, args ...interface{}) (sum int64, err error) {
	return ctx.DBCommand.QueryCount(sql, args...)
}

// QueryMax query max value with sql, return interface{}
func (ctx *PostgresSqlDBContext) QueryMax(sql string, args ...interface{}) (data interface{}, err error) {
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
func (ctx *PostgresSqlDBContext) QueryMin(sql string, args ...interface{}) (data interface{}, err error) {
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
