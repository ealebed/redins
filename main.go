package main

import (
	"time"

	"github.com/ealebed/redins/k8s"
)

func main() {
	kubeInformerFactory := k8s.NewInformerFactory()

	stop := make(chan struct{})
	defer close(stop)

	kubeInformerFactory.Start(stop)
	for {
		time.Sleep(time.Second)
	}
}
