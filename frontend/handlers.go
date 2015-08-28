package main

import (
	"github.com/abates/vpanel"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func NewRouter(manager *vpanel.Manager) *httprouter.Router {
	h := new(handlers)
	h.manager = manager
	router := httprouter.New()
	router.Handler("GET", "/api/host/stats", h.HostStats)
	router.Handler("GET", "/api/container/templates", h.Templates)

	router.Handler("POST", "/api/container", h.CreateContainer)
	router.Handler("GET", "/api/container/:id", h.GetContainer)
	router.Handler("UPDATE", "/api/container/:id", h.UpdateContainer)
	router.Handler("DELETE", "/api/container/:id", h.DestroyContainer)
	router.Handler("GET", "/api/container/:id/stats", h.ContainerStats)

	return router
}

type appHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, *APIError)

func (fn appHandler) Handle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	i, err := fn(w, r, p)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		if e := json.NewEncoder(w).Encode(i); e != nil {
			log.Printf("Failed to encode output: %v", e)
		}
	}
}

func decode(r *http.Request, limit int64, v interface{}) *APIError {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return ioError(err)
	}

	if err := r.Body.Close(); err != nil {
		return ioError(err)
	}

	return decoderError(json.Unmarshal(body, v))
}

type handlers struct {
	manager *vpanel.Manager
}

func (h *handlers) HostStats(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	return monitor.HostStats, nil
}

func (h *handlers) Templates(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	return vpanel.Templates()
}

func (h *handlers) CreateContainer(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	metadata := vpanel.NewContainerMetadata()
	if err := decode(r, 10240, &metadata); err != nil {
		return nil, err
	}

	if !metadata.Valid() {
		return nil, metadata.Err
	}

	return metadata, manager.CreateContainer(metadata)
}

func (h *handlers) GetContainer(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	id := p.ByName("id")
}

func (h *handlers) UpdateContainer(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	id := p.ByName("id")
}

func (h *handlers) DestroyContainer(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	id := p.ByName("id")
}

func (h *handlers) ContainerStats(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, *APIError) {
	id := p.ByName("id")
}
