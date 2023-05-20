package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

// Interfaces

type KeyValueStore interface {
	Put(key string, value string) bool
	Get(key string) (string, bool)
	Delete(key string) bool
}

// LevelDB

type LevelDBStore struct {
	disk *leveldb.DB
}

func (kv LevelDBStore) Put(key string, value string) bool {
	err := kv.disk.Put([]byte(key), []byte(value), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).Error(fmt.Sprintf("Failed to PUT: %s", err))
	}
	return err == nil
}

func (kv LevelDBStore) Get(key string) (string, bool) {
	data, err := kv.disk.Get([]byte(key), nil)
	var value string
	if err == nil {
		value = string(data)
	} else {
		log.WithFields(log.Fields{
			"key": key,
		}).Info(fmt.Sprintf("Failed to GET: %s", err))
	}
	return value, err == nil
}

func (kv LevelDBStore) Delete(key string) bool {
	err := kv.disk.Delete([]byte(key), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"key": key,
		}).Warn(fmt.Sprintf("Failed to DELETE: %s", err))
	}
	return err == nil
}

// In Memory Map

type InMemoryMapStore struct {
	mem map[string]string
}

func (kv InMemoryMapStore) Put(key string, value string) bool {
	kv.mem[key] = value
	return true
}

func (kv InMemoryMapStore) Get(key string) (string, bool) {
	value, ok := kv.mem[key]
	return value, ok
}

func (kv InMemoryMapStore) Delete(key string) bool {
	delete(kv.mem, key)
	return true
}
