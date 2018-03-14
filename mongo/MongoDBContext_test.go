package mongo

import (
	"testing"
	"gopkg.in/mgo.v2/bson"
)

const (
	testMongoConn          = "mongodb://127.0.0.1:27017/"
	testMongoDbName        = "test"
	testObjectID           = "5a068700d3f05e12e844df86"
	collectionName_Demo = "Demo" //Demo表名
)

type DemoInfo struct {
	ID          bson.ObjectId `bson:"_id"`
	DemoID   int `bson:"DemoID"`
	DemoName string `bson:"DemoName"`
}

var db = NewMongoDBContext(testMongoConn, testMongoDbName)

func init(){
	db.DefaultCollectionName = collectionName_Demo
}

func Test_InsertBlob(t *testing.T) {
	demo := &DemoInfo{
		ID:          NewObjectId(),
		DemoID:          1,
		DemoName:        "key1",
	}
	err:=db.InsertBlob(demo)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(demo)
}

func TestMongoDBContext_FindOne(t *testing.T) {
	selector := bson.M{"DemoID": 1}
	var data DemoInfo
	err := db.FindOne(selector, &data)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestMongoDBContext_FindList(t *testing.T) {
	selector := Selector{"DemoID": 1}
	var data []DemoInfo
	err := db.FindList(selector, 0, 0, &data)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestMongoDBContext_Remove(t *testing.T) {
	selector := Selector{"_id": ObjectIdHex(testObjectID)}
	err := db.Remove(selector)
	if err!= nil{
		t.Error(err)
		return
	}
	t.Log(selector)
}