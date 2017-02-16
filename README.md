# Annotator

Kubernetes [Operator](https://coreos.com/blog/introducing-operators.html) that annotates pods with
Docker labels. Gets the labels from the [MicroBadger API](https://microbadger.com/api).

[![](https://images.microbadger.com/badges/image/rossf7/annotator.svg)](https://microbadger.com/images/rossf7/annotator "Get your own image badge on microbadger.com") [![](https://images.microbadger.com/badges/commit/rossf7/annotator.svg)](https://microbadger.com/images/rossf7/annotator "Get your own image badge on microbadger.com")

## Warning

This is an experimental operator at an alpha stage. Please do not deploy to production.

## License

Apache-2.0

## Limitations

* Only labels from the first container are annotated.
* Only public images on Docekr Hub are currently suppported.

## Installation

* Create the operator.

```
$ kubectl apply -f deployment.yaml
```

## Building

Build the Go binary and Docker image. Developed using Go 1.7 and Kubernetes
1.5 using Minikube.

```
$ make
```

## Run tests

Run the unit tests.

```
$ make test
```
