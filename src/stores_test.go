package main

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestLevelDBStore(t *testing.T) {
	temp_dir, err := os.MkdirTemp("", "impulse")
	if err != nil {
		log.Panic(fmt.Sprintf("Could not make temp dir: %s", err))
	}

	disk_db, err := leveldb.OpenFile(temp_dir+"/test_level.db", nil)
	if err != nil {
		log.Panic(fmt.Sprintf("Failed to open LevelDB: %s", err))
	}
	defer disk_db.Close()

	store := LevelDBStore{disk: disk_db}

	store.Put("name", "Stephen")
	read_result, success := store.Get("name")
	assert.Equal(t, "Stephen", read_result)
	assert.Equal(t, true, success)

	store.Delete("name")
	read_result, success = store.Get("name")
	assert.Equal(t, "", read_result)
	assert.Equal(t, false, success)

	store.Put("name", "John Doe")
	store.Put("name", "Jane Doe")
	read_result, success = store.Get("name")
	assert.Equal(t, "Jane Doe", read_result)
	assert.Equal(t, true, success)
}

func TestInMemoryMapStore(t *testing.T) {
	store := InMemoryMapStore{mem: make(map[string]string)}

	store.Put("name", "Stephen")
	read_result, success := store.Get("name")
	assert.Equal(t, "Stephen", read_result)
	assert.Equal(t, true, success)

	store.Delete("name")
	read_result, success = store.Get("name")
	assert.Equal(t, "", read_result)
	assert.Equal(t, false, success)

	store.Put("name", "John Doe")
	store.Put("name", "Jane Doe")
	read_result, success = store.Get("name")
	assert.Equal(t, "Jane Doe", read_result)
	assert.Equal(t, true, success)
}
