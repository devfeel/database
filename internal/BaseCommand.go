package internal

import (
	"database/sql"
	"reflect"
	"github.com/devfeel/mapper"
	"github.com/devfeel/mapper/reflectx"
	"errors"
)

const(
	Zero = 0
)

type BaseCommand struct{
	DriverName string
	OnTrace func(content interface{})
	OnDebug func(content interface{})
	OnInfo func(content interface{})
	OnWarn func(content interface{})
	OnError func(err error, content interface{})
}

func (cmd *BaseCommand) Trace(content interface{}){
	if cmd.OnTrace != nil{
		cmd.OnTrace(content)
	}
}

func (cmd *BaseCommand) Debug(content interface{}){
	if cmd.OnDebug != nil{
		cmd.OnDebug(content)
	}
}

func (cmd *BaseCommand) Info(content interface{}){
	if cmd.OnInfo != nil{
		cmd.OnInfo(content)
	}
}

func (cmd *BaseCommand) Warn(content interface{}){
	if cmd.OnWarn != nil{
		cmd.OnWarn(content)
	}
}

func (cmd *BaseCommand) Error(err error, content interface{}){
	if cmd.OnError != nil{
		cmd.OnError(err, content)
	}
}


// ColScanner is an interface used by MapScan
type ColScanner interface {
	Columns() ([]string, error)
	Scan(dest ...interface{}) error
	Err() error
}

// MapScan scans a single Row into the dest map[string]interface{}.
// Use this to get results for SQL that might not be under your control
// (for instance, if you're building an interface for an SQL server that
// executes SQL from input).  Please do not use this as a primary interface!
// This will modify the map sent to it in place, so reuse the same map with
// care.  Columns which occur more than once in the result will overwrite
// each other!
func  (cmd *BaseCommand) MapScan(r ColScanner, dest map[string]interface{}) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}

func (cmd *BaseCommand) StructScan(rows *sql.Rows, dest interface{})(rowsNum int, err error){
	rowsNum = 0
	var isSlice bool
	destValue := reflect.ValueOf(dest)

	_, errCheckSlice := reflectx.BaseType(destValue.Type(), reflect.Slice)
	if errCheckSlice!= nil{
		isSlice = false
	}else{
		isSlice = true
	}

	if isSlice {
		destSlice := reflect.Indirect(destValue)
		elemType := reflectx.GetSliceType(destValue)
		if elemType.Kind() != reflect.Ptr {
			return Zero, errors.New("slice elem must ptr ")
		}

		for rows.Next() {
			rowsNum += 1
			rowMap := make(map[string]interface{})
			err = cmd.MapScan(rows, rowMap)
			if err != nil {
				return Zero, errors.New("map data error" + err.Error())
			}
			elem := reflect.New(elemType.Elem())
			err = mapper.MapperMap(rowMap, elem.Interface())
			if err != nil {
				return Zero, errors.New("map data error" + err.Error())
			} else {
				destSlice.Set(reflect.Append(destSlice, elem))
			}
		}
	}else{
		elemType := destValue.Type()
		if elemType.Kind() != reflect.Ptr {
			return Zero, errors.New("slice elem must ptr ")
		}
		hasData := rows.Next()
		if hasData{
			rowMap := make(map[string]interface{})
			err = cmd.MapScan(rows, rowMap)
			if err != nil {
				return Zero, errors.New("map data error" + err.Error())
			}
			err = mapper.MapperMap(rowMap, dest)
			if err != nil {
				return Zero, errors.New("map data error" + err.Error())
			}
		}else{
			return Zero, sql.ErrNoRows
		}
	}
	return rowsNum, err
}