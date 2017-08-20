package model

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
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

func SheetHeader() []interface{} {
	return []interface{}{
		"ID",
		"Date",
		"Notes",
		"Branch",
		"Amount",
		"Factor",
		"Balance",
	}
}

func (b *BankContent) SheetContent() []interface{} {
	return []interface{}{
		b.ID.String(),
		b.Date.String(),
		b.Notes,
		b.Branch,
		strconv.Itoa(b.Amount),
		strconv.Itoa(b.Factor),
		strconv.Itoa(b.Amount),
	}
}
