# K8s service discovery

**Project state: POC**

This POC create an application that get IP:port of pods and print them in traefik yaml format.

Possible evolutions:
* add more storage backends (current storage backend: S3 / minIO)
* add more output formats (current output format: traefik yaml)

**Architecture**

* this app is installed in each Kubernetes cluster where data must be extracted
* this app watch K8s endpoints to extract couples of IP:port, called `Backends`
    * this extract can be filtered by K8s labels
* this app save `Backends` lists in the defined output formats, in a storage

Then, a bridge must be create for each load balancer solution. In this POC, a minio client download the backends file and push it in a directory that traefik use for dynamic configuration. The file `traefik/cfg/main.yml` is created once by the administrator.

## Development

Requirements:
* minikube (installed)
* skaffold (installed)
* a docker network (bridge) `my-local-env`: `docker network create my-local-env`

1. Run minikube with CNI bridge
```
$ minikube start --network-plugin=cni --cni bridge
```

2. Run the app:

```
$ skaffold dev
```

### Organization

* Structs:
	* `Backend`: a group of pods
		* `Name`: string
		* `Targets`: array of type `Target`
	* `Target`: a target for the load balancer
		* `Ip`: IP address of the pod
		* `Port`: port of the pod
		* `Protocol`: protocol to use
* Functions:
	* `main()`: the controller
	* `discovery*(clientset *kubernetes.Clientset, backends *[]backend)`: Use the clientset to populate `backends`
	* `convertTo*(backends []Backend) *bytes.Buffer`: Return a buffer with the content in a specific format (ex: `TraefikYaml`)

### Tests

1. Configure your host network to be able to access to pod (Pod IP range: 10.244.0/24, minikube IP: 192.168.49.2)
```
# ip route add 10.244.0/24 via 192.168.49.2
# ip route add 10.96.0.0/12 via 192.168.49.2
```

2. Get a pod IP:port from the output of the app, and try to get an nginx welcome page (with `curl` or from your web browser)
```
$ curl http://IP:port
```

3. Scale nginx:
```
$ kk scale --replicas=3 deployment angular-green
```

4. Add Traefic
```
$ cd traefik && docker compose up -d
$ curl http://localhost -H 'Host: angular.local'
```

### Guide

Add go package:
```
docker run -it -v ${PWD}/src:/app -w /app golang:1.19 go get github.com/minio/minio-go/v7
```
