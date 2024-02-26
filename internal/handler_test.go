package internal

import (
	"bytes"
	"encoding/json"
	"main/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	getInfoURL  = "/info"
	getCacheURL = "/get?key=active-tabs"
	setCacheURL = "/set"

	mockSuccessResponse = model.GetInfoResponse{
		Code: 0,
		Msg:  "Success",
		Records: []model.Record{
			{
				Key:        "TAKwGc6Jr4i8Z487",
				CreatedAt:  time.Date(2017, time.January, 28, 1, 22, 14, 0, time.UTC),
				TotalCount: 1998,
			},
		},
	}

	mockStatusMethodNotAllowedResponse = model.GetInfoResponse{
		Code:    http.StatusMethodNotAllowed,
		Msg:     "Method Not Allowed",
		Records: nil,
	}

	mockStatusBadRequestResponse = model.GetInfoResponse{
		Code:    http.StatusBadRequest,
		Msg:     "startDate is not in the valid format YYYY-MM-DD",
		Records: nil,
	}

	expectedRequest = model.GetInfoRequest{
		StartDate: "2016-01-26",
		EndDate:   "2018-02-02",
		MinCount:  2700,
		MaxCount:  3000,
	}
)

type SetCacheRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func TestHandler_GetInfo(t *testing.T) {
	mockService, mockServiceController := createMockService(t)
	defer mockServiceController.Finish()

	handler := NewHandler(mockService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return info properly", func(t *testing.T) {
		mockService.EXPECT().GetInfo(gomock.Any()).Return(model.GetInfoResponse{
			Code:    0,
			Msg:     "Success",
			Records: mockSuccessResponse.Records,
		})

		expectedRequestJSON, _ := json.Marshal(expectedRequest)
		reqBody := bytes.NewBuffer(expectedRequestJSON)

		req, err := http.NewRequest(http.MethodPost, server.URL+getInfoURL, reqBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		var response model.GetInfoResponse
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.Nil(t, err)
		assert.Equal(t, mockSuccessResponse, response)
	})
	t.Run("should return method not allowed error when request type is not Post", func(t *testing.T) {
		expectedRequestJSON, _ := json.Marshal(expectedRequest)
		reqBody := bytes.NewBuffer(expectedRequestJSON)

		req, err := http.NewRequest(http.MethodGet, server.URL+getInfoURL, reqBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
		var response model.GetInfoResponse
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.Nil(t, err)
		assert.Equal(t, mockStatusMethodNotAllowedResponse, response)
	})
	t.Run("should return status bad request error when time format is invalid", func(t *testing.T) {
		expectedRequestJSON, _ := json.Marshal(model.GetInfoRequest{
			StartDate: "invalidStartDate",
			EndDate:   "2018-02-02",
			MinCount:  2700,
			MaxCount:  3000,
		})
		reqBody := bytes.NewBuffer(expectedRequestJSON)

		req, err := http.NewRequest(http.MethodPost, server.URL+getInfoURL, reqBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		var response model.GetInfoResponse
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.Nil(t, err)
		assert.Equal(t, mockStatusBadRequestResponse, response)
	})
}

func TestHandler_SetCache(t *testing.T) {
	mockService, mockServiceController := createMockService(t)
	defer mockServiceController.Finish()

	handler := NewHandler(mockService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return cache properly", func(t *testing.T) {
		expectedRequestJSON, _ := json.Marshal(SetCacheRequest{
			Key:   "active-tabs",
			Value: "getir",
		})
		reqBody := bytes.NewBuffer(expectedRequestJSON)

		req, err := http.NewRequest(http.MethodPost, server.URL+setCacheURL, reqBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("should return bad request when request body is empty", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, server.URL+setCacheURL, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("should return method not allowed error when request type is not Post", func(t *testing.T) {
		expectedRequestJSON, _ := json.Marshal(SetCacheRequest{
			Key:   "active-tabs",
			Value: "getir",
		})
		reqBody := bytes.NewBuffer(expectedRequestJSON)

		req, err := http.NewRequest(http.MethodGet, server.URL+setCacheURL, reqBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	})
}

func TestHandler_GetCache(t *testing.T) {
	mockService, mockServiceController := createMockService(t)
	defer mockServiceController.Finish()

	handler := NewHandler(mockService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("should return bad request when query params is empty", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/get", http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("should return method not allowed error when request type is not Get", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, server.URL+getCacheURL, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	})
}

func createMockService(t *testing.T) (*MockActions, *gomock.Controller) {
	t.Helper()

	mockServiceController := gomock.NewController(t)
	mockService := NewMockActions(mockServiceController)

	return mockService, mockServiceController
}
