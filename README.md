# sample-grpc-gateway

Sample [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway) application for DigitalOcean App Platform.

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/bojand/sample-grpc-gateway/tree/main)

## Development

Change Go package paths as required.

Adjust HTTP paths as desired.

Install [buf](https://buf.build/).

Install required gRPC plugins

```sh
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

To generate code:

```sh
buf generate
```

## Extra

This consists of a single service that acts as both the HTTP gRPC Gateway and the actual gRPC service.
We can easily split up this into two separate services for each component. You can view code for this in the [multi branch](https://github.com/bojand/sample-grpc-gateway/tree/multi).

Happy hacking!
