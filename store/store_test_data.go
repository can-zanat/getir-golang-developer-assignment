package store

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type RecordDb struct {
	Key       string    `bson:"key"`
	CreatedAt time.Time `bson:"createdAt"`
	Counts    bson.M    `bson:"counts"`
}

var (
	sampleRecord = &RecordDb{
		Key:       "TAKwGc6Jr4i8Z487",
		CreatedAt: time.Date(2017, time.January, 28, 1, 22, 14, 0, time.UTC),
		Counts:    bson.M{"$numberLong": "2800"},
	}
)
