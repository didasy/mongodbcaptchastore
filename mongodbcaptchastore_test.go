package mongodbcaptchastore_test

import (
	"log"
	"testing"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/JesusIslam/mongodbcaptchastore"
	"github.com/stretchr/testify/assert"
)

const (
	MongoDBURL     = "mongodb://localhost:27017"
	DBName         = "test"
	CollectionName = "captcha"
	CollectionSize = 4096       // bytes
	CollectionNum  = 10         // 10 documents
	Timeout        = 5000000000 // 5 seconds in nanoseconds
	Expiration     = 1000000000 // 1 seconds in nanoseconds

	CaptchaID = "captchaid"
	Clear     = true

	ErrorFailedToDialMongoDBMessage         = "Failed to dial MongoDB: %s"
	ErrorFailedToGetMongoDBDatabasesMessage = "Failed to get databases %s"
	ErrorDatabaseNotFoundMessage            = "Database not found: %s"
	ErrorFailedToDropCollectionMessage      = "Failed to drop collection: %s"
	ErrorFailedToRemoveAllDocumentsMessage  = "Failed to remove all documents: %s"
)

var (
	sess *mgo.Session
	s    *mongodbcaptchastore.Store
	err  error

	captchaData = []byte("This is a mock captcha data")
)

func init() {
	checkDatabase()
}

func TestNewSetGet(t *testing.T) {
	defer catchPanic(t)
	defer cleanupCollection()

	s, err = mongodbcaptchastore.New(MongoDBURL, DBName, CollectionName, CollectionSize, CollectionNum, time.Duration(Timeout), time.Duration(Expiration))
	assert.Nil(t, err)
	assert.NotNil(t, s)

	s.Set(CaptchaID, captchaData)

	data := s.Get(CaptchaID, Clear)
	assert.Equal(t, captchaData, data)
}

func checkDatabase() {
	sess, err = mgo.DialWithTimeout(MongoDBURL, time.Duration(Timeout))
	if err != nil {
		log.Fatalf(ErrorFailedToDialMongoDBMessage, err)
	}

	databaseNames, err := sess.DatabaseNames()
	if err != nil {
		log.Fatalf(ErrorFailedToGetMongoDBDatabasesMessage, err)
	}

	found := false
	for _, d := range databaseNames {
		if d == DBName {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf(ErrorDatabaseNotFoundMessage, err)
	}

	return
}

func catchPanic(t *testing.T) {
	e := recover()
	assert.Nil(t, e)
}

func cleanupCollection() {
	err = sess.DB(DBName).C(CollectionName).DropCollection()
	if err != nil {
		log.Fatalf(ErrorFailedToDropCollectionMessage, err)
	}
}
