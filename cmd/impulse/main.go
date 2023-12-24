package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/st3v3nmw/impulse/internal/server"
)

func main() {
	engine := flag.String("engine", "LEVELDB", "Name of Storage Engine")
	diskDBPath := flag.String("leveldb", "", "Path to LevelDB")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.TraceLevel)
	}

	server := server.NewHTTPServer(*engine, *diskDBPath)
	log.Info("Starting server...")
	fasthttp.ListenAndServe(":3000", server.HandleFastHTTP)
}
