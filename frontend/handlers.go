package main

import (
	"encoding/json"
	"github.com/abates/vpanel"
	//"github.com/julienschmidt/httprouter"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request, map[string]string) (interface{}, error)

func NewRouter() http.Handler {
	m := map[string]map[string]appHandler{
		"GET": {
			"/api/host/stats":          HostStats,
			"/api/container/templates": Templates,
			"/api/container/:id":       GetContainerMetadata,
			"/api/container/:id/stats": ContainerStats,
		},
		"POST": {
			"/api/container": CreateContainer,
		},
		"DELETE": {
			"/api/container/:id": DestroyContainer,
		},
	}

	router := mux.NewRouter()
	for method, routes := range m {
		for route, ah := range routes {
			router.HandleFunc(route, makeHandleFunc(ah)).Methods(method)
		}
	}

	return router
}

func makeHandleFunc(handler appHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		v := mux.Vars(r)
		status := http.StatusOK
		i, err := handler(w, r, v)

		if err != nil {
			switch err.(type) {
			case vpanel.ContainerNotFoundError:
				status = http.StatusNotFound
			case vpanel.ValidationError:
				status = 422 // Unprocessable Entity
			default:
				status = http.StatusInternalServerError
			}
			i = err
		}

		if err = encode(w, status, i); err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

func encode(w http.ResponseWriter, code int, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(code)
		w.Write(bytes)
	}
	return err
}

func decode(r *http.Request, limit int64, v interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, limit))
	if err == nil {
		if err = r.Body.Close(); err == nil {
			err = json.Unmarshal(body, v)
		}
	}

	return err
}

func HostStats(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	return monitor.HostStats(), nil
}

func Templates(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	return vpanel.ContainerTemplates()
}

func CreateContainer(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	metadata := vpanel.NewContainerMetadata()
	if err := decode(r, 10240, &metadata); err != nil {
		return nil, err
	}

	if !metadata.IsValid() {
		return nil, metadata.Err
	}

	return metadata, manager.CreateContainer(metadata)
}

func GetContainerMetadata(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	var metadata vpanel.ContainerMetadata

	container, err := manager.GetContainer(vars["id"])
	if err == nil {
		metadata = container.Metadata
	}
	return metadata, err
}

func DestroyContainer(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	return nil, manager.DestroyContainer(vars["id"])
}

func ContainerStats(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	//id := p.ByName("id")
	return nil, nil
}
