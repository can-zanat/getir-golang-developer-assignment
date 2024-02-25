package internal

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"main/model"
	"net/http"
	"testing"
	"time"
)

var (
	expectedDbResponse = model.DbResponse{
		Records: []model.Record{
			{
				Key:        "TAKwGc6Jr4i8Z487",
				CreatedAt:  time.Date(2017, time.January, 28, 1, 22, 14, 0, time.UTC),
				TotalCount: 1998,
			},
		},
	}
)

func TestService_GetInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStore := NewMockStore(mockCtrl)
	service := NewService(mockStore)

	t.Run("should return info properly", func(t *testing.T) {

		mockStore.EXPECT().GetInfo(expectedRequest).Return(expectedDbResponse, nil)

		response := service.GetInfo(expectedRequest)

		assert.Equal(t, 0, response.Code)
		assert.Equal(t, "Success", response.Msg)
		assert.Equal(t, expectedDbResponse.Records, response.Records)
	})

	t.Run("should return error when error occurs while getting data from db", func(t *testing.T) {

		mockStore.EXPECT().GetInfo(gomock.Any()).Return(model.DbResponse{}, errors.New("database error")).Times(1)

		response := service.GetInfo(expectedRequest)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "Error occurs while getting data from database", response.Msg)
		assert.Empty(t, response.Records)
	})
}
