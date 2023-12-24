package server

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/valyala/fasthttp"
)

type Server struct {
	store KeyValueStore
}

func NewHTTPServer(engine string, diskDBPath string) Server {
	var store KeyValueStore
	switch engine {
	case "LEVELDB":
		diskDB, err := leveldb.OpenFile(diskDBPath, nil)
		if err != nil {
			log.Panicf("Failed to open LevelDB: %s", err)
		}
		defer diskDB.Close()
		store = &LevelDBStore{disk: diskDB}
	case "HASH_MAP":
		store = &HashMapStore{hashMap: make(map[string]string)}
	case "LSM_TREE", "B_TREE":
		log.Panicf("Storage engine %s not implemented yet.", engine)
	default:
		log.Panicf("Storage engine %s does not exist.", engine)
	}

	return Server{store: store}
}

func (server Server) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	key := string(ctx.Path()[1:])
	value := string(ctx.PostBody())
	method := string(ctx.Method())

	switch method {
	case "GET":
		value, success := server.store.Get(key)
		if success {
			ctx.WriteString(value)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	case "PUT", "POST":
		success := server.store.Put(key, value)
		if success {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
	case "DELETE":
		server.store.Delete(key)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}
