# REDis INSert data controller

This repository implements a simple controller for watching new created redis pod and inserting data to it.

## Details

The sample controller uses [client-go library](https://github.com/kubernetes/client-go/tree/master/tools/cache) extensively.

## Running

**Prerequisite**: Since the controller uses `apps/v1` deployments, the Kubernetes cluster version should be greater than 1.9.

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
go build ./

# for run on development
export IN_CLUSTER=false

# run with default kubeconfig path
./redins

# or provide path to kubeconfig manually
./redins -kubeconfig=$HOME/.kube/config
```

```sh
# run controller in cluster
kubectl apply -f deployment/deployment.yaml
```
 
## What happens under the hood?

Controller connects to a Kubernetes cluster, sets up an informer for Pods in default namespace and with label selector `"app=ads-redis-statistic"`, and then starts the Informer run loop. When pods with matched criteria (and the initial warmup of pods when the Store syncs) are added to the cluster, controller initialize redis client, connect to provided redis-server, set key/value in database, get and print key/value from DB, and finally close connection to redis-server.
