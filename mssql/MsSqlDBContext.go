package mssql

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

func (ctx *MsSqlDBContext) FindList(sql string, args ...interface{}) (results []map[string]interface{}, err error) {
	return ctx.DBCommand.Query(sql, args...)
}
