package main

import (
	"context"
	"os"
	"log"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/minio/minio-go/v7"
)

func initKubernetes() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func labelFilter() string {
	var labelFilter = os.Getenv("LABEL_FILTER")
	return labelFilter
}

func discoveryEndpoints(clientset *kubernetes.Clientset, backends *[]Backend) {
	endpoints, err := clientset.CoreV1().Endpoints("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelFilter()})
	if err != nil {
		panic(err.Error())
	}

	for _, item := range endpoints.Items {
		for _, subset := range item.Subsets {
			for _, Port := range subset.Ports {
				port := Port.Port

				var protocol string

				if Port.AppProtocol != nil {
					protocol = *Port.AppProtocol
				} else {
					protocol = string(Port.Protocol)
				}
				backend := Backend{Name: item.ObjectMeta.Name + "-" + strconv.Itoa(int(port))}

				for _, ip := range subset.Addresses {
					target := Target{ Ip: ip.IP, Port: port, Protocol: protocol }
					backend.Targets = append(backend.Targets, target) 
				}

				*backends = append(*backends, backend)
			}
		}
	}
}


func discoveryEndpointsWatch(clientset *kubernetes.Clientset, minioClient *minio.Client) {
	log.Printf("Start to watch endpoints")

	watcher, err := clientset.CoreV1().Endpoints("").Watch(context.TODO(), metav1.ListOptions{LabelSelector: labelFilter()})
    if err != nil {
        log.Fatal(err)
    }

    for event := range watcher.ResultChan() {
        // endpoint := event.Object.(*v1.Endpoint)

        switch event.Type {
	        case watch.Added, watch.Modified, watch.Deleted:
	           process(clientset, minioClient)
        }
    }
}
