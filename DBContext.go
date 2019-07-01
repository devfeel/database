package database

type DBContext interface {
	GetCommand() DBCommand
	Init(conn string)
	ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error)
	Insert(sql string, args ...interface{}) (n int64, err error)
	Update(sql string, args ...interface{}) (n int64, err error)
	Delete(sql string, args ...interface{}) (n int64, err error)
	FindOne(dest interface{}, sql string, args ...interface{}) error
	FindOneMap(sql string, args ...interface{}) (result map[string]interface{}, err error)
	FindList(dest interface{}, sql string, args ...interface{}) error
	FindListMap(sql string, args ...interface{}) (results []map[string]interface{}, err error)
	FindListByPage(dest interface{}, tableName, fields, where, orderBy string, skip, take int, args ...interface{}) error
	Scalar(sql string,args ...interface{})(result interface{},err error)
	Count(sql string, args ...interface{}) (count int64, err error)
	QuerySum(sql string, args ...interface{}) (sum int64, err error)
	QueryMax(sql string, args ...interface{}) (data interface{}, err error)
	QueryMin(sql string, args ...interface{}) (data interface{}, err error)
}
