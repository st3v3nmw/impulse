package main

type KeyValueStore interface {
	Put(key string, value string) bool
	Get(key string) (string, bool)
	Delete(key string) bool
}
