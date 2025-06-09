# Gateway Service

The gateway is the only service exposed to the "Internet".
For the scope of the assessment, this refers to Tec de Monterrey's WAN.
We leveraged [Kong Gateway](https://konghq.com/products/kong-gateway) to 
handle the routing and to send POST requests to a logging service after
each request.

## Traffic

The gateway listens for HTTPS requests on port 443 and routes traffic with
a round-robin balancing algorithm to internal services (it uses HTTP
internally). For the assessment we balanced both the frontend and the
backend with 2 instances each.

SSL certificates were generated with OpenSSL using the `config/cert.sh` bash
script, which contains a single OpenSSL command. The `config/kong.conf` file
tells Kong to allow external connections on port 443 with SSL using the
self-signed certificates and it also allows HTTP/2 connections.

```
# config/kong.conf
# ...

proxy_listen = 0.0.0.0:443 http2 ssl reuseport
ssl_cert = /etc/kong/cert.crt
ssl_cert_key = /etc/kong/cert.key
```

## Local Testing
To test the configurations, we created the `compose.yaml` file located
in this directory, which sets up our logging service wrapper, a MongoDB node,
2 dummy "hello world" services and a Kong gateway (obviously). The container
already has self-signed certificates for HTTPS connections and the container 
port has been mapped to `433` (the container uses `8443` for HTTPS by default).

To try this local approach run `docker compose up` in this directory. You
can look at the container logs to verify the load balancing scheme and you 
can also use `mongosh` or another MongoDB client to verify that logs are stored
in the database.

## Declarative Configuration

Kong was configured to use a static declarative approach with a YAML file
to define routes and plugins. Although it is also possible to connect Kong
to a database for dynamic changes.

Creating routes and balancing load is very straight forward. Just group multiple
tcp sockets as targets under an upstream:

```{yaml}
_format_version: "3.0"

upstreams:
  - name: my-upstream
    algorithm: round-robin # can change
    targets:
      - target: $IP_1:$PORT
      - target: $IP_2:$PORT
```

Then add a service that uses the upstream as a `host` and add a route 

```{yaml}
services:
  - name: my-service
    host: my-upstream
    port: $PORT
    routes:
      - name: my-route
        paths:
          - /my/path
```

### Production Routes
The routes that were used during the assessment were stored under the VCS in
`config/kong.prod.yaml`, but in essence moving from the dummy endpoints to the
"real" services was a very straight-forward task.

## Logging
Since Kong Gateway doesn't have a MongoDB plugin for logs, we used the HTTP log plugin
and wrapped the MongoDB operations under a Golang API, this is also safer since the
database is hidden/abstracted away from the gateway which is exposed to the Internet.

The logging wrapper is located in the `logging` directory. It was dockerized and
published to Docker Hub (to avoid building multiple times).

The `logger/models` directory contains structs to deserialize Kong requests to objects
and to serialize our responses to JSON. Serialization is done with the `encoding/json`
package and its abstracted away using a function in `logger/lib/utils.go`.

The `logger/routes` directory contains a single file with a single Post handler.

The `logger/db` creates a MongoDB connection pool singleton to manage DB operations.

`logger/main.go` is the entrypoint, which creates the HTTP/2 multiplexer, mounts
the POST handler for logs, initializes the DB connection pool singleton and handles
shutdowns gracefully (closing the connection).

```
logger
├── db
│   └── mongo.go
├── Dockerfile
├── go.mod
├── go.sum
├── lib
│   └── utils.go
├── main.go
├── models
│   ├── kong_logs.go
│   └── response.go
├── README.Docker.md
└── routes
    └── logs.go
```
