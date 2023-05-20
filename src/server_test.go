package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestServer(t *testing.T) {
	// Setup LevelDB
	temp_dir, err := os.MkdirTemp("", "impulse")
	if err != nil {
		log.Panic(fmt.Sprintf("Could not make temp dir: %s", err))
	}
	disk_db, err := leveldb.OpenFile(temp_dir+"/test_level.db", nil)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to open LevelDB: %s", err))
	}
	defer disk_db.Close()

	// Setup Server
	store := KeyValueStore{disk: disk_db}
	server := Server{store: &store}

	// Test Put
	req := httptest.NewRequest(http.MethodPut, "/name", strings.NewReader("Stephen"))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", string(data))
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Test Get, Known key
	req = httptest.NewRequest(http.MethodGet, "/name", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	res = w.Result()
	data, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Stephen", string(data))
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Test Get, Unknown key
	req = httptest.NewRequest(http.MethodGet, "/jina", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	res = w.Result()
	data, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Not Found\n", string(data))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// Test Delete
	req = httptest.NewRequest(http.MethodDelete, "/name", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	res = w.Result()
	data, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", string(data))
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Test unsuppported method
	req = httptest.NewRequest(http.MethodPost, "/name", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	res = w.Result()
	data, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Method Not Allowed\n", string(data))
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
