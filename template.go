package main

var indexTemplate string = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>gRPC Gateway Demo</title>
    <link rel="stylesheet" href="https://fonts.xz.style/serve/inter.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@exampledev/new.css@1.1.2/new.min.css">
</head>

<body>
    <header>
        <h1>gRPC Gateway Demo</h1>
    </header>

    <h3>About</h3>
    <br>
    <p>
        DigitalOcean App Platform gRPC Gateway demo.
    </p>
    <h3>Usage</h3>
    <br>
    <pre>http POST \
    '{{ .PublicURL }}/api/hello' \
    name=foo</pre>
    <h3>Proto</h3>
    <pre>
syntax = "proto3";

option go_package = "github.com/bojand/sample-grpc-gateway/proto/helloworld";

package helloworld;

import "google/api/annotations.proto";

service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {
        option (google.api.http) = {
            post: "/api/hello"
            body: "*"
        };
    }
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}</pre>
</body>

</html>`
