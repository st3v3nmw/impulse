package main

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	store *KeyValueStore
}

func (server Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	key := req.URL.Path[1:]
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error("Could not read request body")
	}
	value := string(body)

	log.WithFields(log.Fields{
		"method": req.Method,
		"key":    key,
		"value":  value,
	}).Debug("Received a request")

	switch req.Method {
	case "GET":
		value, success := server.store.Get(key)
		if success {
			fmt.Fprint(res, value)
		} else {
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	case "PUT":
		success := server.store.Put(key, value)
		if success {
			res.WriteHeader(http.StatusNoContent)
		} else {
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	case "DELETE":
		server.store.Delete(key)
		res.WriteHeader(http.StatusNoContent)
	default:
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
