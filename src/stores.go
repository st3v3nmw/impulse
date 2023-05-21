package main

import (
	"fmt"
	"sync"

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

func (kv *LevelDBStore) Put(key string, value string) bool {
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

func (kv *LevelDBStore) Delete(key string) bool {
	err := kv.disk.Delete([]byte(key), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"key": key,
		}).Warn(fmt.Sprintf("Failed to DELETE: %s", err))
	}
	return err == nil
}

// Hash Map

type HashMapStore struct {
	lock     sync.RWMutex
	hash_map map[string]string
}

func (kv *HashMapStore) Put(key string, value string) bool {
	kv.lock.Lock()
	kv.hash_map[key] = value
	kv.lock.Unlock()
	return true
}

func (kv *HashMapStore) Get(key string) (string, bool) {
	kv.lock.RLock()
	value, ok := kv.hash_map[key]
	kv.lock.RUnlock()
	return value, ok
}

func (kv *HashMapStore) Delete(key string) bool {
	kv.lock.Lock()
	delete(kv.hash_map, key)
	kv.lock.Unlock()
	return true
}

// Sorted Strings Table

type SSTableStore struct {
}
