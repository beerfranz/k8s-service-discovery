/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	// "context"
	// "fmt"
	"time"
	"os"
	"strconv"
	"log"
	// "bufio"
    "io"
    "bytes"
    "strings"
    "gopkg.in/yaml.v2"

	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/pkg/api/v1"
	// "k8s.io/apimachinery/pkg/watch"
	// "k8s.io/client-go/rest"

	"github.com/minio/minio-go/v7"
	// "github.com/minio/minio-go/v7/pkg/credentials"
)

type Target struct {
    Ip string          `yaml:"ip"`
    Port  int32        `yaml:"port"`
    Protocol string    
}

type Backend struct {
	Name string        `yaml:"name"`
	Targets []Target   `yaml:"targets"`
}

type Backends struct {
	Backends []Backend `yaml:"backends"`
}

func main() {

	clientset := initKubernetes()

	minioClient := initMinio()

	switch os.Getenv("MODE") {
	case "watch":

		switch os.Getenv("DISCOVERY") {
		case "endpoints", "defaults":
			discoveryEndpointsWatch(clientset, minioClient)
		default:
			log.Fatalf("The discovery '%+v' is not valid. Valid values: endpoints, default", os.Getenv("DISCOVERY"))
		}

	case "loop":
		var sleepString = os.Getenv("SLEEP")
		sleepInt, _ := strconv.Atoi(sleepString)
		var sleep = time.Duration(sleepInt)

		log.Printf("Start in loop mode with a sleep of %+v seconds\n", sleepString)

		for {
			process(clientset, minioClient)
			
			time.Sleep(sleep * time.Second)
		}

	default:
		log.Fatalf("The mode '%+v' is not valid. Valid values: watch or loop", os.Getenv("MODE"))
	}

	log.Printf("Bye bye")
}

func process (clientset *kubernetes.Clientset, minioClient *minio.Client) {
	var object = os.Getenv("MINIO_PATH_KEY")

	var arrBackends []Backend
	
	switch os.Getenv("DISCOVERY") {
		case "endpoints", "defaults":
			discoveryEndpoints(clientset, &arrBackends)
		default:
			log.Fatalf("The discovery '%+v' is not valid. Valid values: endpoints, default", os.Getenv("DISCOVERY"))
	}
	
	if strings.Contains(os.Getenv("OUTPUT_FORMATS"), "traefik_yaml") {
		
		buffer := convertToTraefikYaml(arrBackends)
		log.Println(buffer)

		putObject(minioClient, buffer, object + ".traefik.yaml")
	}

	if strings.Contains(os.Getenv("OUTPUT_FORMATS"), "default_yaml") {
		backends := Backends{Backends: arrBackends}
		yamlData, err := yaml.Marshal(backends)

	    if err != nil {
	        log.Printf("Error while Marshaling. %v", err)
	    }
	    log.Println(string(yamlData))

	    buffer := bytes.NewBufferString("")
	    io.WriteString(buffer, string(yamlData))
	    putObject(minioClient, buffer, object + ".default.yaml")
	}
}