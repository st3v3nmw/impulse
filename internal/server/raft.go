package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"
)

type NodeRole int

const (
	Follower NodeRole = iota
	Candidate
	Leader
)

type RaftState struct {
	nodes *[]string

	term          int
	currentRole   NodeRole
	votedFor      string
	votesReceived []string
	leaderIP      string
}

type RaftHeartbeat struct {
	Term   int    `json:"term"`
	NodeIP string `json:"nodeIP"`
}

type RaftVoteResponse struct {
	Term         int  `json:"term"`
	VotedInFavor bool `json:"votedInFavor"`
}

func NewRaftState(nodes *[]string) (state *RaftState) {
	state = &RaftState{
		nodes:       nodes,
		term:        0,
		currentRole: Follower,
	}
	return state
}

func (state *RaftState) initialize() {
	log.Info("Initializing Raft...")
	ticker := time.NewTicker(10*time.Second + time.Duration(rand.Intn(5000))*time.Millisecond)
	go func() {
		currNodeID := os.Getenv("POD_IP")

		for range ticker.C {
			if currNodeID == state.leaderIP {
				// the current node is the leader, broadcast a heartbeat
				state.sendHeartbeat(currNodeID)
			} else if !slices.Contains(*state.nodes, state.leaderIP) && len(*state.nodes) > 0 {
				// no leader detected
				log.Infof("No leader detected, %s starting election...", currNodeID)

				state.term = state.term + 1
				state.currentRole = Candidate
				state.votedFor = currNodeID
				state.votesReceived = []string{currNodeID}

				voteRequest, _ := json.Marshal(RaftHeartbeat{
					Term:   state.term,
					NodeIP: currNodeID,
				})
				for _, nodeIP := range *state.nodes {
					nodeURL := fmt.Sprintf("http://%s:3000/raft/requestVote", nodeIP)
					resp, err := http.Post(nodeURL, "application/json", bytes.NewBuffer(voteRequest))
					if err != nil {
						continue
					}

					voteResponse := RaftVoteResponse{}
					respBody, _ := io.ReadAll(resp.Body)
					json.Unmarshal(respBody, &voteResponse)
					resp.Body.Close()

					log.Infof("%v", voteResponse)
					if state.currentRole == Candidate && voteResponse.Term == state.term && voteResponse.VotedInFavor {
						// not idempotent, need some form of set
						state.votesReceived = append(state.votesReceived, nodeIP)
					} else if voteResponse.Term > state.term {
						state.term = voteResponse.Term
						state.currentRole = Follower
						state.votedFor = ""
					}
				}

				log.Infof("Received %v votes", len(state.votesReceived))
				if float64(len(state.votesReceived)) >= math.Ceil((float64(len(*state.nodes))+2)/2) {
					state.currentRole = Leader
					state.leaderIP = currNodeID
					log.Infof("%s elected", currNodeID)
					state.sendHeartbeat(currNodeID)
				}
			}
		}
	}()
}

func (state *RaftState) handleVoteRequest(voteRequest RaftHeartbeat) RaftVoteResponse {
	log.Infof("Received vote request from %s", voteRequest.NodeIP)

	if voteRequest.Term > state.term {
		state.term = voteRequest.Term
		state.currentRole = Follower
		state.votedFor = ""
	}

	votedInFavor := voteRequest.Term == state.term && (state.votedFor == "" || state.votedFor == voteRequest.NodeIP)
	if votedInFavor {
		state.votedFor = voteRequest.NodeIP
	}
	return RaftVoteResponse{Term: state.term, VotedInFavor: votedInFavor}
}

func (state *RaftState) sendHeartbeat(currNodeID string) {
	log.Infof("Sending heartbeat from %s", currNodeID)

	heartbeat, _ := json.Marshal(RaftHeartbeat{
		Term:   state.term,
		NodeIP: currNodeID,
	})
	for _, nodeIP := range *state.nodes {
		nodeURL := fmt.Sprintf("http://%s:3000/raft/heartbeat", nodeIP)
		http.Post(nodeURL, "application/json", bytes.NewBuffer(heartbeat))
	}
}

func (state *RaftState) handleHeartbeat(heartbeat RaftHeartbeat) {
	log.Infof("Received heartbeat from %s", heartbeat.NodeIP)

	if heartbeat.Term >= state.term {
		state.term = heartbeat.Term
		state.currentRole = Follower
		state.leaderIP = heartbeat.NodeIP
	}
}
