package mongo

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/devfeel/database/internal"
)

type MongoCommand struct {
	internal.BaseCommand
	Connection string
	DefaultDB  string
}

type Selector map[string]interface{}

const (
	Field_ID = "_id" //field表达式
	DriverName = "mongodb"
)

func ObjectIdHex(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

func NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}

func CreateUpdateSet(selector Selector) Selector {
	return Selector{"$set": selector}
}



func (cmd *MongoCommand) SetConn(conn, dbName string) error {
	cmd.Connection = conn
	cmd.DefaultDB = dbName

	return nil
}

func (cmd *MongoCommand) GetDataBase() (*mgo.Database, error) {
	session, err := mgo.Dial(cmd.Connection)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	realDBName := cmd.DefaultDB
	db := session.DB(realDBName)
	return db, nil
}

func (impl *MongoCommand) getDataBaseByName(dbName string) (*mgo.Database, error) {
	session, err := mgo.Dial(impl.Connection)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	realDBName := dbName
	if dbName == "" {
		realDBName = impl.DefaultDB
	}
	db := session.DB(realDBName)
	return db, nil
}

/*新增Json数据插入指定的Collection
* Author: Panxinming
* LastUpdateTime: 2016-10-17 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) InsertJson(collectionName, jsonData string) error {
	logTitle := getLogTitle("insertJson", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()
	c := db.C(collectionName)

	var f interface{}
	data := []byte(jsonData)
	//特殊处理，需要先json序列化，成为标准json，再进一步转为bosn，才能成功
	err_jsonunmar := json.Unmarshal(data, &f)
	if err_jsonunmar != nil {
		cmd.Error(err_jsonunmar, logTitle+"json.Unmarshal["+jsonData+"]error - "+err_jsonunmar.Error())
		return err_jsonunmar
	}
	bjsondata, err_jsonmar := json.Marshal(f)
	if err_jsonmar != nil {
		cmd.Error(err_jsonmar, logTitle+"json.Marshal["+jsonData+"]error - "+err_jsonmar.Error())
		return err_jsonmar
	}

	var bf interface{}
	bsonerr := bson.UnmarshalJSON(bjsondata, &bf)
	if bsonerr != nil {
		cmd.Error(bsonerr, logTitle+"bson.UnmarshalJSON["+jsonData+"]error - "+bsonerr.Error())
		return bsonerr
	}
	err = c.Insert(&bf)
	if err != nil {
		cmd.Error(err, logTitle+"["+jsonData+"]error - "+err.Error())
		return err
	} else {
		cmd.Debug(logTitle + "[" + jsonData + "]Success - " + jsonData)
	}
	return nil
}

/*新增实体数据插入指定Collection
* Author: Panxinming
* LastUpdateTime: 2017-02-26 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) InsertBlob(collectionName string, data interface{}) error {
	logTitle := getLogTitle("insertBlob", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	err = c.Insert(data)
	if err != nil {
		cmd.Error(err, logTitle+"["+fmt.Sprint(data)+"]error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success")
	}
	return err
}

/*更新指定的实体数据
* Author: Panxinming
* LastUpdateTime: 2017-11-11 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) UpdateBlob(collectionName string, selector interface{}, data interface{}) error {
	logTitle := getLogTitle("updateBlob", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	err = c.Update(selector, data)
	if err != nil {
		cmd.Error(err, logTitle+"["+fmt.Sprint(data)+"]error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success")
	}
	return err
}

/*更新指定字段
* Author: Panxinming
* LastUpdateTime: 2017-11-11 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) Update(collectionName string, selector interface{}, data interface{}) error {
	logTitle := getLogTitle("update", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	err = c.Update(selector, data)
	if err != nil {
		cmd.Error(err, logTitle+"["+fmt.Sprint(data)+"]error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success")
	}
	return err
}

/*移除指定查询条件的记录
* Author: Panxinming
* selectot:
* LastUpdateTime: 2017-11-11 14:00
* 如果失败，返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) Remove(collectionName string, selector interface{}) error {
	logTitle := getLogTitle("remove", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	err = c.Remove(selector)
	if err != nil {
		cmd.Error(err, logTitle+"error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "Success")
	}
	return err
}

/*查询指定查询条件的第一条单条数据
* Author: Panxinming
* selectot:
* LastUpdateTime: 2017-02-22 10:00
* 如果失败，result为具体的document并返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) FindOne(collectionName string, selector interface{}, result interface{}) (err error) {
	logTitle := getLogTitle("findOne", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	err = c.Find(selector).One(result)
	if err != nil {
		cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "]Success")
	}
	return err
}

/*查询指定查询条件的最后一条数据
* Author: Panxinming
* selectot:
* LastUpdateTime: 2017-11-11 10:00
* 如果失败，result为具体的document并返回具体error，成功则返回nil
 */
func (cmd *MongoCommand) FindLastOne(collectionName string, selector interface{}, result interface{}) (err error) {
	logTitle := getLogTitle("findLastOne", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	err = c.Find(selector).Sort("-" + Field_ID).One(result)
	if err != nil {
		cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "]Success")
	}
	return err
}

/*查询指定查询条件的批量数据
* Author: Panxinming
* selectot:
* LastUpdateTime: 2017-11-11 10:00
* 支持分页，如果不需要分页，skip、limit请传入0
* 如果失败，则返回具体的document list，成功则返回nil
 */
func (cmd *MongoCommand) FindList(collectionName string, selector interface{}, skip, limit int, result interface{}) (err error) {
	logTitle := getLogTitle("findList", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	if limit > 0 {
		err = c.Find(selector).Skip(skip).Limit(limit).All(result)
	} else {
		err = c.Find(selector).Skip(skip).All(result)
	}
	if err != nil {
		cmd.Error(err, logTitle+"["+fmt.Sprint(selector)+"]["+fmt.Sprint(skip, limit)+"]error - "+err.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(selector) + "]["+fmt.Sprint(skip, limit)+"]Success")
	}
	return err
}

/*更新指定条件的数据，如果没有找到则直接插入数据
* Author: Panxinming
* LastUpdateTime: 2017-02-22 10:00
* 如果失败，则返回具体的error，成功则返回nil
 */
func (cmd *MongoCommand) UpsertBlob(collectionName string, selector interface{}, data interface{}) error {
	logTitle := getLogTitle("upsertBlob", collectionName)
	db, err := cmd.GetDataBase()
	if err != nil {
		cmd.Error(err, logTitle+"getDataBase error - "+err.Error())
		return err
	}
	defer db.Session.Close()

	c := db.C(collectionName)
	change, errUpsert := c.Upsert(selector, data)
	if errUpsert != nil {
		cmd.Error(errUpsert, logTitle+"["+fmt.Sprint(data)+"]error - "+errUpsert.Error())
	} else {
		cmd.Debug(logTitle + "[" + fmt.Sprint(data) + "]Success -> " + fmt.Sprint(change))
	}
	return errUpsert
}

func (cmd *MongoCommand) GetStat() (result map[string]interface{}, err error) {
	db, err := cmd.GetDataBase()
	if err != nil {
		return nil, err
	}
	defer db.Session.Close()
	result = make(map[string]interface{})
	err = db.Run(bson.M{"serverStatus": 1}, result)
	return result, err
}

// getLogTitle return log title
func getLogTitle(commandName string, collectionName string) string {
	return "database.MongoCommand:" + commandName + "[" + collectionName + "]:"
}


