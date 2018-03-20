# database

#### version 0.2
* DBCommand新增Select接口，支持传入切片指针，直接mapper struct slice
* DBContext新增FindMap，返回[]map[string]interface{}
* DBContext调整FindList，支持传入切片指针，直接mapper struct slice
* 2018-03-20 11:00

#### version 0.1
* 初始版本，目前支持mongodb、mysql、mssql
* 2018-03-14 12:00