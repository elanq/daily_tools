package mongo

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

//Mongodb driver struct. it contains
//1. session data
//2. database name
//3. collection name
type MongoDriver struct {
	Session    *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
}

//init mongodb driver
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

//make insert command to mongodb
func (d *MongoDriver) Insert(document interface{}) error {
	err := d.Collection.Insert(document)
	return err
}

//search one data in database
func (d *MongoDriver) FindOne(selector interface{}, result interface{}) error {
	err := d.Collection.FindId(selector).One(result)
	return err
}

//find data by query
func (d *MongoDriver) Find(selector interface{}, result interface{}) error {
	err := d.Collection.Find(selector).All(result)
	return err
}

//remova data
func (d *MongoDriver) Remove(selector interface{}) error {
	err := d.Collection.RemoveId(selector)
	return err
}
