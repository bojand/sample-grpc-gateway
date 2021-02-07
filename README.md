# sample-grpc-gateway

Sample [gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway) application for DigitalOcean App Platform.

### apigw
The HTTP service that sets up the gRPC gateway.

### samplesvc
The actual gRPC service that handles all the gRPC call.

## Instructions

Clone the repo.

Adjust `./do/app.yaml` to match up your GitHub URLs.

Create app using [doctl](https://github.com/digitalocean/doctl):

```sh
doctl apps create --spec .do/app.yaml
```

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

Happy hacking!
