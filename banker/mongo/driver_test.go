package mongo_test

import (
	"os"
	"testing"
	"time"

	"github.com/elanq/daily_tools/banker/model"
	"github.com/elanq/daily_tools/banker/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"gopkg.in/mgo.v2/bson"
)

type DriverSuite struct {
	suite.Suite
	DBName         string
	CollectionName string
	Driver         *mongo.MongoDriver
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(DriverSuite))
}

func (d *DriverSuite) SetupSuite() {
	gotenv.Load("../env.sample")
	d.DBName = os.Getenv("DB_NAME")
	d.CollectionName = "banker_test_record"
	d.Driver = mongo.NewMongoDriver(d.DBName, d.CollectionName)
}

func buildBankContent() *model.BankContent {
	return &model.BankContent{
		Date:    time.Now(),
		Notes:   "Notes",
		Branch:  "1234",
		Amount:  1000000,
		Factor:  1,
		Balance: 2000000,
	}
}

func (d *DriverSuite) TestNewMongoDriver() {
	assert.NotNil(d.T(), d.Driver.Session, "Session should not be nil")
	assert.NotNil(d.T(), d.Driver.Database, "Driver should not be nil")
	assert.NotNil(d.T(), d.Driver.Collection, "Collection should not be nil")
}

func TestInvalidMongoDriver(t *testing.T) {
	os.Setenv("MONGODB_HOSTS", "127.0.0.1:12345")
	assert.Panics(t, func() {
		mongo.NewMongoDriver("derp", "durpp")
	}, "should be panic")
}

func (d *DriverSuite) TestInsert() {
	sampleContent := buildBankContent()
	err := d.Driver.Insert(sampleContent)
	assert.Nil(d.T(), err, "Should not return any error")
}

func (d *DriverSuite) TestFindOne() {
	//only generated for testing. Not intented for real world use
	objId := bson.NewObjectId()
	sampleContent := buildBankContent()
	sampleContent.ID = objId
	err := d.Driver.Insert(sampleContent)
	assert.Nil(d.T(), err, "Should not return any error")

	var result model.BankContent

	err = d.Driver.FindOne(sampleContent.ID, &result)
	assert.NotNil(d.T(), result, "Should return result")
	assert.Nil(d.T(), err, "Should not return any error")
}

func (d *DriverSuite) TestFind() {
	var results []model.BankContent
	objId := bson.NewObjectId()
	sampleContent := buildBankContent()
	sampleContent.ID = objId
	d.Driver.Insert(sampleContent)

	//unsorted find
	err := d.Driver.Find(bson.M{}, &results)
	assert.NotEmpty(d.T(), results, "should not empty")
	assert.Nil(d.T(), err, "Should not return any error")

	//sorted find
	err = d.Driver.Find(bson.M{}, &results, "-date")
	assert.NotEmpty(d.T(), results, "should not empty")
	assert.Nil(d.T(), err, "Should not return any error")
}

func (d *DriverSuite) TestRemove() {
	var results []model.BankContent

	err := d.Driver.Find(bson.M{}, &results)
	assert.NotEmpty(d.T(), results, "should not empty")
	assert.Nil(d.T(), err, "Should not return any error")

	for _, result := range results {
		err = d.Driver.Remove(result.ID)
		assert.Nil(d.T(), err, "Should not return any error")
	}
}
