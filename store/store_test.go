package store

import (
	"context"
	"main/config"
	"main/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

const (
	mongoImage = "mongo:7.0.4"
)

func prepareTestStore() (store *Mongo, clean func()) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage(mongoImage))
	if err != nil {
		panic(err)
	}

	cleanFunc := func() {
		if cleanErr := mongodbContainer.Terminate(ctx); cleanErr != nil {
			panic(err)
		}
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	s := NewStore(config.Mongo{URI: uri})

	if err != nil {
		panic(err)
	}

	return s, cleanFunc
}

func TestStore_GetInfo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	store, clean := prepareTestStore()
	defer clean()

	collection := store.Client.Database("getircase-study").Collection("records")

	_, err := collection.InsertOne(context.Background(), sampleRecord)
	assert.Nil(t, err)

	test := []struct {
		Name           string
		StartDate      time.Time
		EndDate        time.Time
		MinCount       int
		MaxCount       int
		ExpectedResult model.DBResponse
		ExpectedError  error
	}{
		//{
		//	Name:      "should return record properly",
		//	StartDate: time.Date(2015, 1, 26, 0, 0, 0, 0, time.UTC),
		//	EndDate:   time.Date(2019, 2, 2, 0, 0, 0, 0, time.UTC),
		//	MinCount:  2700,
		//	MaxCount:  3000,
		//	ExpectedResult: model.DBResponse{
		//		Records: []model.Record{
		//			{
		//				Key:        sampleRecord.Key,
		//				CreatedAt:  sampleRecord.CreatedAt,
		//				TotalCount: 5700,
		//			},
		//		},
		//	},
		// },
		{
			Name:      "should return nil when record not found",
			StartDate: time.Date(2015, 1, 26, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2019, 2, 2, 0, 0, 0, 0, time.UTC),
			MinCount:  2700,
			MaxCount:  3000,
			ExpectedResult: model.DBResponse{
				Records: []model.Record(nil),
			},
		},
	}

	for _, tc := range test {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := store.GetInfo(model.GetInfoRequest{
				StartDate: tc.StartDate.Format("2006-01-02"),
				EndDate:   tc.EndDate.Format("2006-01-02"),
				MinCount:  tc.MinCount,
				MaxCount:  tc.MaxCount,
			})
			assert.Equal(t, tc.ExpectedError, err)
			assert.Equal(t, tc.ExpectedResult, result)
		})
	}
}
