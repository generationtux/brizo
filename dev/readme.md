# Brizo Dev

This directory contains development resources for Brizo.

## Kubernetes

To create default Kubernetes resources (namespace, test RCs, pods, etc.) use kubectl apply on the kube folder in this directory
```sh
$ kubectl apply -f ./dev/kube
$ kubectl get pods --namespace=brizo
```
