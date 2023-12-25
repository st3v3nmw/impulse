package server

import (
	"encoding/json"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/valyala/fasthttp"
)

// enums

type ReplicationMode int

const (
	NoReplication ReplicationMode = iota
	SingleLeaderReplication
	LeaderlessReplication
)

// Server

type Server struct {
	store           KeyValueStore
	replicationMode ReplicationMode
	cluster         *Cluster
}

func NewHTTPServer(engine string, replMode string, diskDBPath string) *Server {
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

	var replicationMode ReplicationMode
	switch replMode {
	case "NONE":
		replicationMode = NoReplication
	case "SINGLE_LEADER":
		replicationMode = SingleLeaderReplication
	case "LEADERLESS":
		log.Panicf("Replication mode %s not implemented yet.", replMode)
	default:
		log.Panicf("Unknown replication mode %s", replMode)
	}

	server := Server{store: store, replicationMode: replicationMode, cluster: NewCluster()}
	if replicationMode != NoReplication {
		server.cluster.discover()
	}
	if replicationMode == SingleLeaderReplication {
		server.cluster.raftState.initialize()
	}
	return &server
}

func (server *Server) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
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
	case "PUT":
		success := server.store.Put(key, value)
		if success {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
	case "DELETE":
		server.store.Delete(key)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	case "POST":
		// Leader election, etc calls are handled via POST
		server.handlePOST(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func (server *Server) handlePOST(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/raft/requestVote":
		voteRequest := RaftHeartbeat{}
		json.Unmarshal(ctx.PostBody(), &voteRequest)
		voteResponse, _ := json.Marshal(server.cluster.raftState.handleVoteRequest(voteRequest))
		ctx.SetContentType("application/json")
		ctx.SetBody(voteResponse)
	case "/raft/heartbeat":
		heartbeat := RaftHeartbeat{}
		json.Unmarshal(ctx.PostBody(), &heartbeat)
		server.cluster.raftState.handleHeartbeat(heartbeat)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}
}
