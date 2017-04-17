package mongo

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

type MongoDriver struct {
	Session    *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
}

func NewMongoDriver(dbName string, collectionName string) *MongoDriver {
	session, err := mgo.DialWithTimeout(os.Getenv("MONGODB_HOSTS"), 3*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to DB"))
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbName)
	collection := db.C(collectionName)
	return &MongoDriver{
		Session:    session,
		Database:   db,
		Collection: collection,
	}
}

func (d *MongoDriver) Insert(document interface{}) error {
	err := d.Collection.Insert(document)
	return err
}

func (d *MongoDriver) BulkInsert(documents ...interface{}) error {
	bulk := d.Collection.Bulk()
	bulk.Unordered()
	bulk.Insert(documents)
	_, err := bulk.Run()

	return err
}

func (d *MongoDriver) FindOne(selector interface{}, result interface{}) error {
	err := d.Collection.FindId(selector).One(result)
	return err
}

func (d *MongoDriver) Find(selector interface{}, result interface{}) error {
	err := d.Collection.Find(selector).All(result)
	return err
}

func (d *MongoDriver) Remove(selector interface{}) error {
	err := d.Collection.RemoveId(selector)
	return err
}
