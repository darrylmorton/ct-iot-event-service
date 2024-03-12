# ct-iot-event-service

## Description
The `ct-iot-event-service` is a RESTful microservice that also consumes messages from a `thing events queue`, and stores them in the `thing events database`.

[Diagrams](./docs/DIAGRAMS.md)

## Requirements
- Go >= v1.20
- Docker Desktop >= v4.17.0 (99724)

## Build
```
make build
```

## Docker Compose Startup
```
docker compose up
```

## Test
```
make test
```

## Docker Compose Shutdown
```
docker compose down -v
```

## Documentation
This can be found here: [docs/DOCUMENTATION.md](docs/DOCUMENTATION.md).
