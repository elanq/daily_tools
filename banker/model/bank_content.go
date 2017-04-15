package model

import "gopkg.in/mgo.v2/bson"

type BankContent struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Date    string        `bson:"date"`
	Notes   string        `bson:"notes"`
	Branch  string        `bson:"branch"`
	Amount  int           `bson:"amount"`
	Factor  int           `bson:"factor"`
	Balance int           `bson:"balance"`
}
