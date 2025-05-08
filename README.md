This is distributed calculator

(Or a bad Apache Spark cosplay?)

The system consists of nodes of 2 types:
- Orchestrator
- Agent

Also, the web interface is provided

# Table of Contents
1. [Getting Started](#getting-started)
   - [Requirements](#requirements)
   - [Start Up](#start-up)
     - [Docker Compose](#docker-compose)
2. [Services](#services)
   - [Scheme](#scheme)
   - [API Gateway](#api-gateway)
   - [Orchestrator](#orchestrator)
   - [Agent](#agent)
   - [BFF](#bff)
3. [Explanatory note](#explanatory-note)
   - [Why MongoDB?](#why-mongodb)
   - [Why Redis?](#why-redis)
   - [Why nginx?](#why-nginx)
   - [Why gRPC streams?](#why-grpc-streams)
4. [Good to Know](#good-to-know)
   - [General](#general)
   - [Expressions](#expressions)
5. [Examples of Use](#examples-of-use)

# Getting Started

## Requirements
Before you start, you may want to have the following dependencies installed:

### Mandatory
- Docker 28.0.0 and above

## Start Up
- [Docker Compose](#docker-compose)

### Docker Compose
Since external dependencies (such as mongodb, redis) are required to run up, 
it is only possible to run it in docker for your convenience

**Build and run all services**

```shell
docker compose up --build
```

***Or***

**Build all images**

```shell
docker compose build
```

**Run all services**

```shell
docker compose up
```

# Services

## Scheme
<img src="backend/docs/assets/images/DistributedCalc.jpg" alt="Prject scheme"/>

System consists of 4 components:
- API Gateway
- Orchestrator
- Agent
- BFF

## API Gateway
API gateway is an entry point for every incoming requests. 
It passes requests to appropriate backend. Nginx is used as api gateway.

## Orchestrator
Orchestrator is the master node of distributed calculator

- It provides REST API for client requests
- It provides gRPC API for internal requests
- It decomposes expressions into atomic tasks and enqueues them for agent to process

### Configuration
Orchestrator can be configured via environment variables

`HOST`: Host to run on (default: `0.0.0.0`)

> **NOTICE**: Do not change host if you run in docker, otherwise it may not work properly 

`HTTP_PORT`: http port to listen to (default: `8080`), must be in range between `1` and `65535`

`GRPC_PORT`: gRPC port to listen to (default: `50051`), must be in range between `1` and `65535`

`TIME_ADDITION_MS`: Time in milliseconds which `+` operation takes (default: `1`), must be non-negative integer

`TIME_SUBTRACTION_MS`: Time in milliseconds which `-` operation takes (default: `1`), must be non-negative integer

`TIME_MULTIPLICATION_MS`: Time in milliseconds which `*` operation takes (default: `1`), must be non-negative integer

`TIME_DIVISION_MS`: Time in milliseconds which `/` operation takes (default: `1`), must be non-negative integer

`MONGO_HOST`: MongoDB host

`MONGO_PORT`: MongoDB port

`MONGO_USER`: MongoDB user

`MONGO_PASSWORD`: MongoDB password

`MONGO_NAME`: MongoDB name

`MONGO_MIGRATIONS_PATH`: MongoDB migrations dir

`REDIS_HOST`: Redis host

`REDIS_PORT`: Redis port

`REDIS_USER`: Redis user (can be omitted)

`REDIS_PASSWORD`: Redis password (can be omitted)

## Agent
Agent is a slave node of distributed calculator

- It pulls tasks from orchestrator to process and sends the result back after processing
- It supports horizontal scaling via ***reverse proxy***

> **NOTICE**: On start up, agent will try to connect to orchestrator. It will exit immediately on failure after retries

### Configuration
Agent can be configured via environment variables

`WORKERS_LIMIT``: Amount of active workers per agent instance (default: `10`), must be positive integer

`BUFFER_SIZE`: Size of task buffer (default: `128`), must be positive integer

`POLL_TIMEOUT`: Defines how often task process results will be sent back to orchestrator (default: `50`), 
must be positive integer

`MAX_RETRIES`: Maximum retries on failed requests (default: `3`), must be positive integer

`ORCHESTRATOR_HOST`: Orchestrator host

`ORCHESTRATOR_PORT`: Orchestrator port

## BFF
BFF is a service which serves static frontend files and proxies requests to primary backend 

### Configuration
BFF can be configured via environment variables

`BACKEND_HOST`: Backend host

`BACKEND_PORT`: Backend port

# Explanatory note

## Why MongoDB?
MongoDB is used to keep users data, expressions and task persistently
- MongoDB is schemaless which means that no schema is required
- MongoDB is performant
- MongoDB is developer-friendly
Overall, MongoDB just was the most convenient choice for this project

## Why Redis?
Redis is used as JWT blacklist
- Redis is performant in-memory key-val db, so it is good to be used as JWT blacklist

## Why nginx?
- Nginx is flexible and convenient
In this project nginx serves as rather an API gateway than load balancer 
to provide single entry point for all incoming requests

## Why gRPC streams?
Orchestrator and agent use gRPC streaming to communicate because it is performant asynchronous realtime 
messaging implementation and because it was required by Task Conditions, 
though in real world we would use rather Kafka/RabbitMQ than gRPC streaming 
due to reliability issues

# Good to Know

## General
> **NOTICE**: Signature private key is kept in memory and is generated on every startup, so all issued tokens become invalid as orchestrator is shut down
- The system is completely stateless
- If result of expressions has more than `8` decimal places, they are thrown away
- Notice that expressions like `2 2 + 3` will be processed as `22+3` due to the system design

## Expressions
1. During the evaluation, field `result` in Expressions schema is `0` until expression is evaluated
2. May have several statuses:
   - `pending`: the expression is being processed
   - `completed`: the expression is processed and result is ready for use
   - `failed`: the system failed to process the expression

# Examples of Use
Since authorization tokens are required on most requests, specific examples are no longer provided. 
Please, use [API docs](backend/api/v1/api.yaml) to enforce your own