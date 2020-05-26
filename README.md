# Cartographer

[![Go](https://img.shields.io/badge/go-1.14-00E5E6.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/docker-19.03-2885E4.svg)](https://www.docker.com/)
[![Seabolt](https://img.shields.io/badge/seabolt-1.7.4-2885E4.svg)](https://github.com/neo4j-drivers/seabolt)
[![Build Status](https://travis-ci.org/dynastymasra/cartographer.svg?branch=master)](https://travis-ci.org/dynastymasra/cartographer)
[![Coverage Status](https://coveralls.io/repos/github/dynastymasra/cartographer/badge.svg?branch=master)](https://coveralls.io/github/dynastymasra/cartographer?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dynastymasra/cartographer)](https://goreportcard.com/report/github.com/dynastymasra/cartographer)

Service to serve information about Administrative Division a country

## Libraries

Use [Go Module](https://blog.golang.org/using-go-modules) for install all dependencies required this application.

## How To Run and Deploy

Before run this service. Make sure all requirements dependencies has been installed likes **Golang, Docker, and database Neo4J**

### Local

Use command go ```go run main.go``` in root folder for run this application.

### Docker

**cartographer** uses docker multi stages build, minimal docker version is **17.05**. If docker already installed use command.

This command will build the images.
```bash
docker build -f Dockerfile -t cartographer:$(VERSION) .
```

To run service use this command
```bash
docker run --name cartographer -d -e ADDRESS=:8080 -e <environment> $(IMAGE):$(VERSION)
```

## Test

For run unit test, from root project you can go to folder or package and execute command
```bash
go test -v -cover -coverprofile=coverage.out -covermode=set
go tool cover -html=coverage.out
```
`go tool` will generate GUI for test coverage. Available package or folder can be tested

- `/country`
- `/country/handler`
- `/region`
- `/region/handler`
- `/infrastructure/web/handler`

## Environment Variables

+ `SERVER_PORT` - Address application is used default is `8080`
+ `LOGGER_LEVEL` - Log level(debug, info, error, warn, etc)
+ `LOGGER_FORMAT` - Format specific for log
  - `text` - Log format will become standard text output, this used for development
  - `json` - Log format will become *JSON* format, usually used for production
+ `NEO4J_ADDRESS` - Neo4J database address `bolt+routing://<host>:<port>`
  - `bolt+routing://` - Used with causal cluster
  - `bolt://` - Used with single server
+ `NEO4J_USERNAME` - Neo4J database username
+ `NEO4J_PASSWORD` - Neo4J database password
+ `NEO4J_MAX_CONN_POOL` - Neo4j maximum number of connections per URL to allow on this driver
+ `NEO4J_ENCRYPTED` - Neo4J whether to turn on/off TLS encryption (`true`/`false`)
+ `NEO4J_LOG_ENABLED` - Neo4J database log enabled (`true`/`false`)
+ `NEO4J_LOG_LEVEL` Neo4J type that default logging implementations use for available default `0`
  - `0` - Doesn't generate any output
  - `1` - Level error
  - `2` - Level warning
  - `3` - Level info
  - `4` - Level debug

## API Documentation

This service use [GraphQL](https://graphql.org/) to serve the request, [![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/45953192904281df47f8)

## Available Administrative Division

+ **Indonesia** - Base on `PMDN 72 TH 2019`, Reference:
  - [Ministry of Home Affairs](https://www.kemendagri.go.id/files/2020/PMDN%2072%20TH%202019+lampiran.pdf)
  - [Github cahyadsn](https://github.com/cahyadsn/wilayah)