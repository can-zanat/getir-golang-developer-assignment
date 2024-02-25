package internal

import (
	"main/model"
	"main/store"
	"net/http"
)

type Service struct {
	store store.Store
}

type Store interface {
	GetInfo(request model.GetInfoRequest) (model.DbResponse, error)
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) GetInfo(request model.GetInfoRequest) model.GetInfoResponse {

	infoData, err := s.store.GetInfo(request)
	if err != nil {
		return model.GetInfoResponse{
			Code:    http.StatusInternalServerError,
			Msg:     "Error occurs while getting data from database",
			Records: nil,
		}
	}

	var infoResponse model.GetInfoResponse
	infoResponse = convertToGetInfoResponse(infoData)

	return model.GetInfoResponse{
		Code:    0,
		Msg:     "Success",
		Records: infoResponse.Records,
	}
}

func convertToGetInfoResponse(data model.DbResponse) model.GetInfoResponse {
	var records []model.Record
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
