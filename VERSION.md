# database

#### version 0.6
* New Feature: DBContext add SetPasswordCallback(PasswordCallback)
* New Feature: add Counter used to stat db query
* New Feature: add database.ShowStateData to show stat data with map[string]int64
* How to get stat data:
    - call database.ShowStateData() to get current stat data
    - data like map[INSERT:1 SELECT:1 SELECT_ERROR:1 UPDATE:1] 
* 2019-08-04 12:00

#### version 0.5
* New Feature: add postgreSql support
* Architecture: add Scalar in DBContext & DBCommand used to query data on first column&row
* Architecture: remove QuerySum in DBContext & DBCommand
* 2019-07-01 23:00

#### version 0.4.3
* Architecture: add interface DBCommand and DBContext for mssql\mysql\sqlite
* Architecture: make mssql\mysql\sqlite uniform and implement interface DBCommand and DBContext
* 2019-06-27 23:00

#### version 0.4.2
* Bug Fixed: modify mysql to sqlite in logTitle
* 2019-06-14 14:00

#### version 0.4.1
* Bug Fixed: modify mysql to sqlite3 in sqlite module
* Bug Fixed: add Default_FieldName_Max\Min for QueryMax\QueryMin
* 2019-06-05 14:00

#### version 0.4
* New Feature: add sqlite support
* 2019-05-30 18:00

#### version 0.3.4
* New Feature: add mongodb Count
* 2018-11-21 10:00

#### version 0.3.3
* Fixed Bug: fixed mysql Count() can not get correct value
* 2018-09-26 10:00

#### version 0.3.2
* New Feature: add mysql ExecProc
* 2018-08-15 10:00

#### version 0.3.1
* Fixed Bug: fixed mongodb conn pool not effective
* 2018-06-15 10:00

#### version 0.3
* MsSqlDBContext\MySqlDBContext新增Count\QuerySum，返回int64, error
* MsSqlDBContext\MySqlDBContext新增QueryMax\QueryMin，返回interface{}, error
* 2018-05-25 13:00

#### version 0.2.2
* 修复ExecProc无论有无数据均返回空数据BUG
* 2018-05-21 13:00

#### version 0.2.1
* FindOne增加sql.ErrNoRows，当无对应结果时，返回sql.ErrNoRows
* 2018-03-28 13:00

#### version 0.2
* DBCommand新增Select接口，支持传入切片指针，直接mapper struct slice
* DBContext新增FindOneMap\FindListMap，返回map[string]interface{}\[]map[string]interface{}
* DBContext调整FindOne\FindList，支持传入单对象指针\切片指针，直接mapper struct\struct slice
* 2018-03-20 11:00

#### version 0.1
* 初始版本，目前支持mongodb、mysql、mssql
* 2018-03-14 12:00