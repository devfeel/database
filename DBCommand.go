package database

import "database/sql"

type DBCommand interface {
	SetOnTrace(log func(content interface{}))
	SetOnDebug(log func(content interface{}))
	SetOnInfo(log func(content interface{}))
	SetOnWarn(log func(content interface{}))
	SetOnError(log func(err error, content interface{}))

	Trace(content interface{})
	Debug(content interface{})
	Info(content interface{})
	Warn(content interface{})
	Error(err error, content interface{})

	ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error)
	Exec(commandText string, args ...interface{}) (result sql.Result, err error)
	Select(dest interface{}, commandText string, args ...interface{}) (rowsNum int, err error)
	Query(commandText string, args ...interface{}) (records []map[string]interface{}, err error)
	QueryCount(commandText string, args ...interface{}) (int64, error)
}
