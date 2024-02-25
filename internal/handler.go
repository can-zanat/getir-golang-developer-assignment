package internal

import (
	"encoding/json"
	"main/model"
	"net/http"
	"sync"
)

// CacheStore is a simple in-memory store with a mutex to handle concurrent access.
type CacheStore struct {
	sync.RWMutex
	store map[string]string
}

// NewCacheStore initializes a new CacheStore.
func NewCacheStore() *CacheStore {
	return &CacheStore{
		store: make(map[string]string),
	}
}

// Set adds a key-value pair to the store.
func (c *CacheStore) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.store[key] = value
}

// Get retrieves the value for a key from the store.
func (c *CacheStore) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.store[key]
	return val, ok
}

type Handler struct {
	service Actions
	cache   *CacheStore
}

type Actions interface {
	GetInfo(request model.GetInfoRequest) model.GetInfoResponse
}

func NewHandler(service Actions) *Handler {
	return &Handler{
		service: service,
		cache:   NewCacheStore(),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/info", h.GetInfo)
	mux.HandleFunc("/set", h.SetCache())
	mux.HandleFunc("/get", h.GetCache())
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

func (h *Handler) SetCache() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.cache.Set(data.Key, data.Value)

		response, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(response)
	}
}

func (h *Handler) GetCache() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if value, ok := h.cache.Get(key); ok {
			response, err := json.Marshal(map[string]string{"key": key, "value": value})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)
		} else {
			http.Error(w, "Key not found", http.StatusNotFound)
		}
	}
}
