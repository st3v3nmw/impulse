package server

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Cluster struct {
	nodes     []string
	events    chan string
	raftState *RaftState
}

func NewCluster() (cluster *Cluster) {
	cluster = &Cluster{nodes: []string{}, events: make(chan string)}
	cluster.raftState = NewRaftState(&cluster.nodes)
	return cluster
}

func (cluster *Cluster) discover() (err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		currentPodIP := os.Getenv("POD_IP")

		for range ticker.C {
			pods, err := clientset.CoreV1().Pods("impulse").List(context.TODO(), metav1.ListOptions{
				FieldSelector: "status.phase=Running",
			})
			if err != nil {
				log.Panic(err.Error())
			}

			cluster.nodes = nil
			for _, pod := range pods.Items {
				podIP := pod.Status.PodIP
				if podIP != currentPodIP {
					cluster.nodes = append(cluster.nodes, podIP)
				}
			}
		}
	}()

	return nil
}
