package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	engine := flag.String("engine", "LEVELDB", "Name of Storage Engine")
	disk_db_path := flag.String("leveldb", "", "Path to LevelDB")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.TraceLevel)
	}

	var store KeyValueStore
	switch *engine {
	case "LEVELDB":
		disk_db, err := leveldb.OpenFile(*disk_db_path, nil)
		if err != nil {
			log.Panic(fmt.Sprintf("Failed to open LevelDB: %s", err))
		}
		defer disk_db.Close()
		store = LevelDBStore{disk: disk_db}
	case "IN_MEMORY_MAP":
		log.Panic(fmt.Sprintf("Storage engine %s not implemented:", *engine))
	case "SSTABLE", "LSM_TREE", "B_TREE":
		log.Panic(fmt.Sprintf("Storage engine %s not implemented:", *engine))
	default:
		log.Panic(fmt.Sprintf("Storage engine %s does not exist:", *engine))
	}

	server := Server{store: store}
	http.ListenAndServe(":3000", server)
}
