![Version](https://img.shields.io/badge/version-0.0.3-orange.svg)
![Test](https://github.com/bdemirpolat/kubecd/actions/workflows/kubecd.yml/badge.svg)
![Go](https://img.shields.io/github/go-mod/go-version/bdemirpolat/kubecd)
![MIT License](https://img.shields.io/github/license/bdemirpolat/kubecd)
[![Documentation](https://godoc.org/github.com/bdemirpolat/kubecd?status.svg)](https://pkg.go.dev/github.com/bdemirpolat/kubecd)
[![Go Report Card](https://goreportcard.com/badge/github.com/bdemirpolat/kubecd)](https://goreportcard.com/report/github.com/bdemirpolat/kubecd)

## What is kubecd?
kubecd is a continous deployment tool for kubernetes.

## How its work?
kubecd check your git repository at regular intervals and
if detect any changes in your specified k8s manifest directory, sends new manifests to kubernetes.




# Development
This part contains information about how to start develop kubecd. 

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

## Run CLI
```
go run cmd/appcli/cli.go application --help
```



