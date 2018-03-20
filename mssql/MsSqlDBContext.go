package mssql

import "strconv"

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

func (ctx *MsSqlDBContext) FindOne(sql string, args ...interface{}) (result map[string]interface{}, err error) {
	results, err := ctx.DBCommand.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(results) <= 0 {
		return make(map[string]interface{}), nil
	}
	return results[0], nil
}

// FindList query data with sql and return struct slice
// slice's elem type must ptr
func (ctx *MsSqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	return ctx.DBCommand.Select(dest, sql, args...)
}

// FindMap query data with sql and return []map[string]interface{}
func (ctx *MsSqlDBContext) FindMap(sql string, args ...interface{}) (results []map[string]interface{}, err error){
	return ctx.DBCommand.Query(sql, args...)
}

// FindListByPage query single table data by skip and take
// args: args for where string param
// call demo:
// var results []*Demo
// FindListByPage(&results, "Demo", "*", "DemoID = ?", "ID ASC, DemoName DESC", 10, 10, 10000)
func (ctx *MsSqlDBContext) FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{})  (err error) {
	if fields == ""{
		fields = "*"
	}
	where = "WHERE " + where
	orderBy = "ORDER BY " + orderBy
	sql := "SELECT * FROM ( "
	sql += "SELECT ROW_NUMBER() OVER ("+orderBy+") AS [ROW_NUMBER], "+fields+" FROM " +tableName+" AS t0 WITH(NOLOCK) " + where
	sql += ") AS tp WHERE [tp].[ROW_NUMBER] BETWEEN "+ strconv.Itoa(skip)+" + 1 AND "+ strconv.Itoa(take)+" + "+ strconv.Itoa(skip)+" ORDER BY [tp].[ROW_NUMBER]"
	return ctx.DBCommand.Select(dest, sql, args...)
}