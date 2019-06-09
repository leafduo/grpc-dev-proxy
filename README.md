#  grpc-dev-proxy
grpc-dev-proxy is a proxy that allows debugging gRPC services using HTTP GUI tools like [Postman](https://www.getpostman.com/) or [Paw](https://paw.cloud/) .

## Feature
- Supports gRPC [server reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md)
- Supports metadata

## Requirements
- gRPC [reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) must be enabled on the server. It is easy to enable on supported languages listed [here](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md#known-implementations).

## Installation
```bash
go get github.com/leafduo/grpc-dev-proxy
```

## Usage
Start `grpc-dev-proxy`:

```bash
grpc-dev-proxy
```

By default `grpc-dev-proxy` listens on port 2333, you can change it by passing `--port`:

```bash
grpc-dev-proxy --port 1234 
```

Start your HTTP debugging tool, set the URL to `http://localhost:2333` and method to `POST`.

Set the following HTTP headers:

- `Target`: gRPC host and port
- `Service`: gRPC service name
- `Method`: method you are calling

Put request message in HTTP request body in JSON format, don’t forget to change `Content-Type` to `application/json`.

If gRPC metadata is needed, it should go in the HTTP query string.

## Example
For example, we have the following `.proto` file:

```protobuf
syntax = "proto3";

package helloworld;

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

The gRPC server is located at `localhost:50051`, and we want to make a call to `helloworld.Gretter/SayHello` with `name` set to `world` and metadata `user-id = 1`.

### cURL

```bash
curl -X "POST" "http://localhost:2333/?user-id=1" \
     -H 'Target: localhost:50051' \
     -H 'Service: helloworld.Greeter' \
     -H 'Method: SayHello' \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "name": "world"
}'
```

### Paw

![headers](https://raw.githubusercontent.com/leafduo/grpc-dev-proxy/assets/Paw%202019-06-09%20at%2018.55.56%402x.png)
![body](https://raw.githubusercontent.com/leafduo/grpc-dev-proxy/assets/Paw%202019-06-09%20at%2018.56.07%402x.png)
![metadata](https://raw.githubusercontent.com/leafduo/grpc-dev-proxy/assets/Paw%202019-06-09%20at%2018.56.12%402x.png)
