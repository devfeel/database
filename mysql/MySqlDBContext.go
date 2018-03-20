package mysql

const (
	Default_OPEN_CONNS = 50
	Default_IDLE_CONNS = 50
)

type MySqlDBContext struct {
	DBCommand        *MySqlCommand
	DefaultTableName string
}

func NewMySqlDBContext(connString string) *MySqlDBContext {
	db := new(MySqlDBContext)
	db.Init(connString)
	return db
}

func (ctx *MySqlDBContext) Init(conn string){
	ctx.DBCommand = new(MySqlCommand)
	ctx.DBCommand.DriverName = DriverName
	ctx.DBCommand.Connection = conn
	ctx.DBCommand.PoolOpenConnsCount = Default_OPEN_CONNS
	ctx.DBCommand.PoolIdleConnsCount = Default_IDLE_CONNS
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

func (ctx *MySqlDBContext) FindOne(sql string, args ...interface{}) (result map[string]interface{}, err error) {
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
func (ctx *MySqlDBContext) FindList(dest interface{}, sql string, args ...interface{}) error {
	return ctx.DBCommand.Select(dest, sql, args...)
}

// FindMap query data with sql and return []map[string]interface{}
func (ctx *MySqlDBContext) FindMap(sql string, args ...interface{}) (results []map[string]interface{}, err error){
	return ctx.DBCommand.Query(sql, args...)
}

