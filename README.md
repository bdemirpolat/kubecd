## What is kubecd?
kubecd is a continous deployment tool for kubernetes.

## How its work?
kubecd check your git repository at regular intervals and
if detect any changes in your specified k8s manifest directory, sends new manifests to kubernetes.


### Specify database file
```
export KUBECD_DBFILE="/Users/X/kubecd/kubecd.db"
```

### Specify clone path for your projects.
```
export KUBECD_CLONE_PATH="/Users/X/kubecd"
```

### Choose k8s cluster type
If kubecd working out of cluster, you should set like below
```
export KUBECD_CLUSTER_TYPE="OUT_OF_CLUSTER"
```

Default cluster type is "IN_CLUSTER". 

If kubecd working in same cluster with your repositories you do not need to set "KUBECD_CLUSTER_TYPE" variable. 


## Run Server
```
go run cmd/appserver/server.go
```