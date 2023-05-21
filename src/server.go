package main

import (
	"github.com/valyala/fasthttp"
)

type Server struct {
	store KeyValueStore
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
