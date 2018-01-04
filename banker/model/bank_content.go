package model

import (
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// BankContent encapsulates required bank record content
type BankContent struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Date    time.Time     `bson:"date" json:"date"`
	Notes   string        `bson:"notes" json:"notes"`
	Branch  string        `bson:"branch" json:"branch"`
	Amount  int           `bson:"amount" json:"amount"`
	Factor  int           `bson:"factor" json:"factor"`
	Balance int           `bson:"balance" json:"balance"`
}

//SheetHeader is array of string that represent the title for each column
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

//SheetContent returns the value of BankContent object
func (b *BankContent) SheetContent() []interface{} {
	return []interface{}{
		b.ID.String(),
		b.Date.String(),
		b.Notes,
		b.Branch,
		strconv.Itoa(b.Amount),
		strconv.Itoa(b.Factor),
		strconv.Itoa(b.Balance),
	}
}
