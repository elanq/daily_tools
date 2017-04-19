package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type BankContent struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Date    time.Time     `bson:"date" json:"date"`
	Notes   string        `bson:"notes" json:"notes"`
	Branch  string        `bson:"branch" json:"branch"`
	Amount  int           `bson:"amount" json:"amount"`
	Factor  int           `bson:"factor" json:"factor"`
	Balance int           `bson:"balance" json:"balance"`
}
