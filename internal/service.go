package internal

import (
	"main/model"
	"main/store"
	"net/http"
)

// Service defines the business logic layer.
type Service struct {
	store store.Store
}

// Store defines methods for accessing data.
type Store interface {
	GetInfo(request model.GetInfoRequest) (model.DBResponse, error)
}

// NewService creates a new Service instance.
func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) GetInfo(request model.GetInfoRequest) model.GetInfoResponse {
	// Retrieve information from the store.
	infoData, err := s.store.GetInfo(request)
	if err != nil {
		return model.GetInfoResponse{
			Code:    http.StatusInternalServerError,
			Msg:     "Error occurs while getting data from database",
			Records: nil,
		}
	}

	// Convert the retrieved data to the appropriate response format.
	infoResponse := convertToGetInfoResponse(infoData)

	return model.GetInfoResponse{
		Code:    0,
		Msg:     "Success",
		Records: infoResponse.Records,
	}
}

// convertToGetInfoResponse converts database response data to GetInfoResponse format.
func convertToGetInfoResponse(data model.DBResponse) model.GetInfoResponse {
	records := make([]model.Record, 0, len(data.Records))
	for _, info := range data.Records {
		records = append(records, model.Record{
			Key:        info.Key,
			CreatedAt:  info.CreatedAt,
			TotalCount: info.TotalCount,
		})
	}

	return model.GetInfoResponse{
		Records: records,
	}
}
