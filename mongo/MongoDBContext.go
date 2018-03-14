package mongo

type MongoDBContext struct {
	DBCommand             *MongoCommand
	DefaultCollectionName string
}

func NewMongoDBContext(conn, dbName string) *MongoDBContext{
	db := new(MongoDBContext)
	db.Init(conn, dbName)
	return db
}

func (ctx *MongoDBContext) Init(conn, dbName string){
	ctx.DBCommand = new(MongoCommand)
	ctx.DBCommand.DriverName = DriverName
	ctx.DBCommand.SetConn(conn, dbName)
}


/*插入Json数据
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (ctx *MongoDBContext) InsertJson(jsonData string) error {
	return ctx.DBCommand.InsertJson(ctx.DefaultCollectionName, jsonData)
}

/*新增实体数据插入指定Collection
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (ctx *MongoDBContext) InsertBlob(data interface{}) error {
	return ctx.DBCommand.InsertBlob(ctx.DefaultCollectionName, data)
}

/*更新指定的实体数据
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (ctx *MongoDBContext) UpdateBlob(selector interface{}, data interface{}) error {
	return ctx.DBCommand.UpdateBlob(ctx.DefaultCollectionName, selector, data)
}

/*更新指定字段
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (ctx *MongoDBContext) Update(selector interface{}, update interface{}) error {
	return ctx.DBCommand.Update(ctx.DefaultCollectionName, selector, update)
}

/*移除指定查询条件的记录
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，返回具体error，成功则返回nil
 */
func (ctx *MongoDBContext) Remove(selector interface{}) error {
	return ctx.DBCommand.Remove(ctx.DefaultCollectionName, selector)
}

/*查询指定查询条件的第一条单条数据
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，result为具体的document并返回具体error，成功则返回nil
 */
func (ctx *MongoDBContext) FindOne(selector interface{}, result interface{}) (err error) {
	return ctx.DBCommand.FindOne(ctx.DefaultCollectionName, selector, result)
}

/*查询指定查询条件的最后一条数据
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，result为具体的document并返回具体error，成功则返回nil
 */
func (ctx *MongoDBContext) FindLastOne(selector interface{}, result interface{}) (err error) {
	return ctx.DBCommand.FindLastOne(ctx.DefaultCollectionName, selector, result)
}

/*查询指定查询条件的批量数据
* LastUpdateTime: 2018-01-02 10:00
* 支持分页，如果不需要分页，skip、limit请传入0
* 如果失败，则返回具体的document list，成功则返回nil
 */
func (ctx *MongoDBContext) FindList(selector interface{}, skip, limit int, result interface{}) (err error) {
	return ctx.DBCommand.FindList(ctx.DefaultCollectionName, selector, skip, limit, result)
}

/*更新指定条件的数据，如果没有找到则直接插入数据
* LastUpdateTime: 2018-01-02 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (ctx *MongoDBContext) UpsertBlob(selector interface{}, data interface{}) error {
	return ctx.DBCommand.UpsertBlob(ctx.DefaultCollectionName, selector, data)

}
