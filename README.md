# agones-shutdown example

## Testing

Using [minikube](https://agones.dev/site/docs/installation/creating-cluster/minikube/)

```
minikube start
make build-image
go test ./...
```