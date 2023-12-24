package server

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestLevelDBStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "impulse")
	if err != nil {
		log.Panic(fmt.Sprintf("Could not make temp dir: %s", err))
	}

	diskDB, err := leveldb.OpenFile(tempDir+"/test_level.db", nil)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to open LevelDB: %s", err))
	}
	defer diskDB.Close()

	store := LevelDBStore{disk: diskDB}

	store.Put("name", "Stephen")
	readResult, success := store.Get("name")
	assert.Equal(t, "Stephen", readResult)
	assert.Equal(t, true, success)

	store.Delete("name")
	readResult, success = store.Get("name")
	assert.Equal(t, "", readResult)
	assert.Equal(t, false, success)

	store.Put("name", "John Doe")
	store.Put("name", "Jane Doe")
	readResult, success = store.Get("name")
	assert.Equal(t, "Jane Doe", readResult)
	assert.Equal(t, true, success)
}

func TestHashMapStore(t *testing.T) {
	store := HashMapStore{hashMap: make(map[string]string)}

	store.Put("name", "Stephen")
	readResult, success := store.Get("name")
	assert.Equal(t, "Stephen", readResult)
	assert.Equal(t, true, success)

	store.Delete("name")
	readResult, success = store.Get("name")
	assert.Equal(t, "", readResult)
	assert.Equal(t, false, success)

	store.Put("name", "John Doe")
	store.Put("name", "Jane Doe")
	readResult, success = store.Get("name")
	assert.Equal(t, "Jane Doe", readResult)
	assert.Equal(t, true, success)
}
