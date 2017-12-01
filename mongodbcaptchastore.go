package mongodbcaptchastore

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

const (
	DefaultDBName = ""

	ConvertToCappedCollectionCommandName    = "convertToCapped"
	MaxSizeCappedCollectionCommandName      = "size"
	MaxDocumentsCappedCollectionCommandName = "max"

	CaptchaIDPropertyName = "captcha_id"
	CreatedAtPropertyName = "created_at"

	TTLIndexName = "ttl"

	GreaterThanEqual = "$gte"
	Set              = "$set"

	RefreshMode = true
)

type Data struct {
	ID        bson.ObjectId `bson:"_id"`
	CreatedAt time.Time     `bson:"created_at"`
	CaptchaID string        `bson:"captcha_id"`
	Digits    []byte        `bson:"digits"`
}

type Store struct {
	collection *mgo.Collection
	expiration time.Duration
}

func (s *Store) Set(id string, digits []byte) {
	err := s.collection.Insert(&Data{
		ID:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		CaptchaID: id,
		Digits:    digits,
	})
	if err != nil {
		panic(err)
	}
}

func (s *Store) Get(id string, clear bool) (digits []byte) {
	data := &Data{}
	err := s.collection.Find(bson.M{
		CaptchaIDPropertyName: id,
	}).One(data)
	if err != nil {
		panic(err)
	}

	digits = data.Digits

	if clear {
		s.collection.RemoveId(data.ID)
	}

	return
}

func New(urlstring, dbName, collectionName string, collectSize, collectNum int, timeout, expiration time.Duration) (s *Store, err error) {
	sess, err := mgo.DialWithTimeout(urlstring, timeout)
	if err != nil {
		return
	}
	sess.SetMode(mgo.Monotonic, RefreshMode)

	db := sess.Clone().DB(dbName)
	coll := db.C(collectionName)

	// check if collection already exists in the database
	// if not create a new capped collection
	// if found, convert to capped collection
	var colls []string
	colls, err = db.CollectionNames()
	if err != nil {
		return
	}

	exists := false
	for _, c := range colls {
		if c == collectionName {
			exists = true
			break
		}
	}

	if !exists {
		err = coll.Create(&mgo.CollectionInfo{
			ForceIdIndex: true,
			Capped:       true,
			MaxBytes:     collectSize,
			MaxDocs:      collectNum,
		})
	} else {
		var res interface{}
		err = db.Run(bson.D{
			{
				ConvertToCappedCollectionCommandName, collectionName,
			},
			{
				MaxSizeCappedCollectionCommandName, collectSize,
			},
			{
				MaxDocumentsCappedCollectionCommandName, collectNum,
			},
		}, res)
	}
	if err != nil {
		return
	}

	// Ensure TTL index exists
	err = coll.EnsureIndex(mgo.Index{
		Key:         []string{CreatedAtPropertyName},
		Name:        TTLIndexName,
		ExpireAfter: expiration,
	})
	if err != nil {
		return
	}

	s = &Store{
		collection: coll,
		expiration: expiration,
	}

	return
}
