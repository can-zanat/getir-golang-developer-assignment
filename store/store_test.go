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
	"go.mongodb.org/mongo-driver/bson"
)

const (
	mongoImage = "mongo:7.0.4"
)

// prepareTestStore creates a MongoDB container for testing and initializes a store instance connected to it.
// It returns the store instance and a cleanup function to terminate the container after testing.
func prepareTestStore() (store *Mongo, clean func()) {
	ctx := context.Background()

	// Start a MongoDB container.
	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage(mongoImage))
	if err != nil {
		panic(err)
	}

	// Define a cleanup function to terminate the container after testing.
	cleanFunc := func() {
		if cleanErr := mongodbContainer.Terminate(ctx); cleanErr != nil {
			panic(err)
		}
	}

	// Retrieve the MongoDB connection string.
	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	// Create a store instance connected to the MongoDB container.
	s, _ := NewStore(config.Mongo{URI: uri, Database: "getircase-study"})

	if err != nil {
		panic(err)
	}

	return s, cleanFunc
}

// TestStore_GetInfo tests the GetInfo method of the store.Mongo type.
func TestStore_GetInfo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Setup test environment.
	store, clean := prepareTestStore()
	defer clean()

	// Prepare MongoDB collection.
	collection := store.Client.Database("getircase-study").Collection("records")

	// Insert a sample record into the MongoDB collection.
	_, err := collection.InsertOne(context.Background(), bson.D{
		{Key: "key", Value: "TAKwGc6Jr4i8Z487"},
		{Key: "createdAt", Value: time.Date(2017, time.January, 28, 1, 22, 14, 0, time.UTC)},
		{Key: "counts", Value: []int64{2800}},
	})
	assert.Nil(t, err)

	// Define test cases.
	test := []struct {
		Name           string
		StartDate      time.Time
		EndDate        time.Time
		MinCount       int
		MaxCount       int
		ExpectedResult model.DBResponse
		ExpectedError  error
	}{
		{
			Name:      "should return record properly",
			StartDate: time.Date(2015, 1, 26, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2019, 2, 2, 0, 0, 0, 0, time.UTC),
			MinCount:  2700,
			MaxCount:  3000,
			ExpectedResult: model.DBResponse{
				Records: []model.Record{
					{
						Key:        "TAKwGc6Jr4i8Z487",
						CreatedAt:  time.Date(2017, time.January, 28, 1, 22, 14, 0, time.UTC),
						TotalCount: 2800,
					},
				},
			},
		},
		{
			Name:      "should return nil when record not found",
			StartDate: time.Date(2018, 1, 26, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2019, 2, 2, 0, 0, 0, 0, time.UTC),
			MinCount:  2700,
			MaxCount:  3000,
			ExpectedResult: model.DBResponse{
				Records: []model.Record(nil),
			},
		},
	}

	// Iterate over test cases.
	for _, tc := range test {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := store.GetInfo(model.GetInfoRequest{
				StartDate: tc.StartDate.Format("2006-01-02"),
				EndDate:   tc.EndDate.Format("2006-01-02"),
				MinCount:  tc.MinCount,
				MaxCount:  tc.MaxCount,
			})

			// Verify the result and error match the expected ones.
			assert.Equal(t, tc.ExpectedError, err)
			assert.Equal(t, tc.ExpectedResult, result)
		})
	}
}
