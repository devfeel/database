# database

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