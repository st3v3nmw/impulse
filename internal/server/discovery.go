package server

import (
	"context"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Discovery struct {
	nodes []string
}

func (discovery Discovery) discover() (err error) {
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

			discovery.nodes = nil
			for _, pod := range pods.Items {
				podIP := pod.Status.PodIP

				if podIP != currentPodIP {
					discovery.nodes = append(discovery.nodes, podIP)
				}
			}
		}
	}()

	return nil
}
