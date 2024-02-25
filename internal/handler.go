package internal

import (
	"encoding/json"
	"main/model"
	"net/http"
)

type Handler struct {
	service Actions
}

type Actions interface {
	GetInfo(request model.GetInfoRequest) model.GetInfoResponse
}

func NewHandler(service Actions) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/info", h.GetInfo)
}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response := model.GetInfoResponse{
			Code:    http.StatusMethodNotAllowed,
			Msg:     "Method Not Allowed",
			Records: nil,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonResponse)
		return
	}

	request := model.GetInfoRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		response := model.GetInfoResponse{
			Code:    http.StatusBadRequest,
			Msg:     err.Error(),
			Records: nil,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	getInfoResponse := h.service.GetInfo(request)

	jsonResponse, err := json.Marshal(getInfoResponse)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
	return
}
