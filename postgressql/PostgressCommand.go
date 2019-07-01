package postgressql

import (
	"strconv"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/a526757124/database/internal"
	_ "github.com/lib/pq"
)

var (
	sqlPoolMap   map[string]*sql.DB
	sqlPoolMutex *sync.RWMutex
)

const (
	DriverName = "pgsql"
)

func init() {
	sqlPoolMap = make(map[string]*sql.DB)
	sqlPoolMutex = new(sync.RWMutex)
}

func getSqlPool(connString string) (*sql.DB, bool) {
	sqlPoolMutex.RLock()
	pool, exists := sqlPoolMap[connString]
	sqlPoolMutex.RUnlock()
	return pool, exists
}

func setSqlPool(connString string, openConnsCount, idleConnsCount int) (*sql.DB, error) {
	dbPool, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	dbPool.SetMaxIdleConns(idleConnsCount)
	dbPool.SetMaxOpenConns(openConnsCount)

	sqlPoolMutex.Lock()
	sqlPoolMap[connString] = dbPool
	sqlPoolMutex.Unlock()

	return dbPool, nil
}

type PostgressCommand struct {
	SqlPool            *sql.DB
	Connection         string
	PoolOpenConnsCount int
	PoolIdleConnsCount int
	internal.BaseCommand
}

// getSqlPool get global conn pool
func (command *PostgressCommand) getSqlPool() (*sql.DB, error) {
	var err error
	pool, exists := getSqlPool(command.Connection)
	if !exists {
		pool, err = setSqlPool(command.Connection, command.PoolOpenConnsCount, command.PoolIdleConnsCount)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}

// ExecProc executes proc with name
func (command *PostgressCommand) ExecProc(procName string, args ...interface{}) (records []map[string]interface{}, err error) {
	var keyValue string
	for i,_:= range args {
		if keyValue != "" {
			keyValue += ","
		}
		keyValue += "$"+ strconv.Itoa((i+1))
	}
	sqlStmt := "SELECT  " + procName + " (#KEY=VALUE#)"
	sqlStmt = strings.Replace(sqlStmt, "#KEY=VALUE#", keyValue, -1)
	logTitle := getLogTitle("ExecProc", sqlStmt)
	records, err = command.Query(sqlStmt, args...)
	if err != nil {
		command.Error(err, logTitle+" error - "+err.Error())
	} else {
		command.Debug(logTitle + " success")
	}
	return records, err
}

// Exec executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (command *PostgressCommand) Exec(commandText string, args ...interface{}) (result sql.Result, err error) {
	logTitle := getLogTitle("Exec", commandText+fmt.Sprint(args...))
	sqlPool, err := command.getSqlPool()
	if err != nil {
		command.Error(err, logTitle+" getSqlPool error - "+err.Error())
		return nil, err
	}
	stmt, err := sqlPool.Prepare(commandText)
	if err != nil {
		command.Error(err, logTitle+" Prepare error - "+err.Error())
		return nil, err
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	result, err = stmt.Exec(args...)
	if err != nil {
		command.Error(err, logTitle+" Exec error - "+err.Error())
	} else {
		command.Debug(logTitle + " Exec success")
	}
	return result, err
}

// Select executes a query that returns dest interface{}, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (command *PostgressCommand) Select(dest interface{}, commandText string, args ...interface{}) (rowsNum int, err error) {
	logTitle := getLogTitle("Select", commandText+fmt.Sprint(args...))
	sqlPool, err := command.getSqlPool()
	if err != nil {
		command.Error(err, logTitle+" getSqlPool error - "+err.Error())
		return internal.Zero, err
	}
	rows, err := sqlPool.Query(commandText, args...)
	if err != nil {
		command.Error(err, logTitle+" Query error - "+err.Error())
		return internal.Zero, err
	} else {
		command.Debug(logTitle + " Query success")
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	return command.StructScan(rows, dest)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (command *PostgressCommand) Query(commandText string, args ...interface{}) (records []map[string]interface{}, err error) {
	logTitle := getLogTitle("Query", commandText+fmt.Sprint(args...))
	sqlPool, err := command.getSqlPool()
	if err != nil {
		command.Error(err, logTitle+" getSqlPool error - "+err.Error())
		return nil, err
	}
	rows, err := sqlPool.Query(commandText, args...)
	if err != nil {
		command.Error(err, logTitle+" Query error - "+err.Error())
		return nil, err
	} else {
		command.Debug(logTitle + " Query success")
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	records = make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		dest := make(map[string]interface{})
		err = command.MapScan(rows, dest)
		if err != nil {
			continue
		}
		records = append(records, dest)
	}
	return records, err
}

// QueryCount executes a query that returns count column
func (command *PostgressCommand) QueryCount(commandText string, args ...interface{}) (int64, error) {
	logTitle := getLogTitle("QueryCount", commandText+fmt.Sprint(args...))
	sqlPool, err := command.getSqlPool()
	if err != nil {
		command.Error(err, logTitle+" getSqlPool error - "+err.Error())
		return 0, err
	}
	rows, err := sqlPool.Query(commandText, args...)
	if err != nil {
		command.Error(err, logTitle+" Query error - "+err.Error())
		return 0, err
	} else {
		command.Debug(logTitle + " Query success")
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	count := int64(0)
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			command.Error(err, logTitle+" scan count error - "+err.Error())
			return 0, err
		}
	}
	return count, nil
}

// getLogTitle return log title
func getLogTitle(commandName, commandText string) string {
	return "database.PostgressCommand:" + commandName + " [" + commandText + "]"
}
