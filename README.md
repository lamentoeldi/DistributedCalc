This is distributed calculator

(Or a bad Apache Spark cosplay?)

The system consists of nodes of 2 types:
- Orchestrator
- Agent

# Table of Contents
1. [Start Up](#start-up)
   - [Command Line](#command-line)
   - [Taskfile](#taskfile)
   - [Docker Compose](#docker-compose)
   - [Docker CLI](#docker-cli)
2. [Services](#services)
   - [Orchestrator](#orchestrator)
   - [Agent](#agent)
3. [Good to Know](#good-to-know)
   - [General](#general)
   - [Expression](#expression)
4. [Examples of Use](#examples-of-use)
   - [/api/v1/calculate](#apiv1calculate)
   - [/api/v1/expressions](#apiv1expressions)
   - [/api/v1/expressions/{id}](#apiv1expressionsid)
5. [Future Plans](#future-plans)
  


# Start Up
You can run calculation cluster in several ways

## Command Line
Though it is advised to use Docker Compose to run app, you can still use console commands to run it

```shell
go mod download
```

```shell
go run cmd/orchestrator/main.go
```

```shell
go run cmd/agent/main.go
```

## Taskfile
Also you can use Taskfile to run app 

```shell
task run
```

## Docker CLI
You can use Docker CLI to build images and then run containers

Use this to build orchestrator image
```shell
docker build -t orchestrator:latest -f ./build/package/orchestrator/Dockerfile ./
```

Use this to build agent image
```shell
docker build -t agent:latest -f ./build/package/agent/Dockerfile ./
```

Use this to run orchestrator
```shell
docker run -d --name orchestrator -p 8080:8080 orchestrator:latest
```

Use this to run agent
```shell
docker run -d --name agent --link orchestrator:orchestrator -e MASTER_URL=http://orchestrator:8080 agent:latest
```

## Docker Compose
Docker Compose is the most preferable way to run app. As mentioned in [compose file](docker-compose.yaml), on default 
the port 8080 of ***Orchestrator*** is bound on 8080 port of local machine

Use this to build and run app
```shell
docker compose up --build
```

Use this to build images
```shell
docker compose build
```

Use this to run app
```shell
docker compose up
```

# Services

## Orchestrator
Orchestrator is a master node of calculation cluster

It decomposes the expression to run in parallel tasks on ***Agent*** instances

### Configuration
Orchestrator can be configured via environment variables

`HOST`: Host to run on (default: `0.0.0.0`)

NOTICE: Do not change host if you run in docker, otherwise it may not work properly 

`PORT`: Port to run on (default: `8080`)

`LOG_LEVEL`: Level of logging (default: `info`)

`TIME_ADDITION_MS`: Time in milliseconds which `+` operation takes (default: `1`)

`TIME_SUBTRACTION_MS`: Time in milliseconds which `-` operation takes (default: `1`)

`TIME_MULTIPLICATION_MS`: Time in milliseconds which `*` operation takes (default: `1`)

`TIME_DIVISION_MS`: Time in milliseconds which `/` operation takes (default: `1`)

## Agent
Agent is a worker node of calculation cluster

It uses long polling to receive tasks via ***Orchestrator*** API

NOTICE: On start up, agent will try to connect to orchestrator. It will exit immediately on failure after retries

### Configuration
Agent can be configured via environment variables

`LOG_LEVEL`: Level of logging (default: `info`)

`COMPUTING_POWER`: Amount of active workers per agent instance (default: `10`)

`POLL_TIMEOUT`: Polling interval in milliseconds (default: `50`)

`MAX_RETRIES`: Maximum retries on failed requests (default: `3`)

`MASTER_URL`: Orchestrator URL in `protocol://host:port` format (default: `http://localhost:8080`)

# Good to Know

## General

- Currently, the system keeps all data in-memory, that means that all data will be lost on restart
- Currently, the system is stateful, that means that data you receive depends on which node you have accessed
- Agents use long polling to receive tasks from orchestrator
- It is possible to use proxy like [envoy](https://www.envoyproxy.io), 
[nginx](https://nginx.org) or 
[traefik](https://doc.traefik.io/traefik/) to balance incoming requests between running nodes
- If result of expressions has more than `8` decimal places, they are thrown away

## Expression
1. During the evaluation, field `result` in Expression schema is `0` until expression is evaluated
2. May have several statuses:
   - `pending`: the expression is being processed
   - `completed`: the expression is processed and result is ready for use
   - `failed`: the system failed to process the expression

# Examples of Use
- [API Specification](api/v1/api.yaml)
- [More examples](examples/api/v1)

## /api/v1/calculate
Send expression to start evaluation

### Request
```http request
POST localhost:8080/api/v1/calculate
Content-Type: application/json

{
  "expression": "3*4+7"
}
```
`expression`: string

### Response
```json
{
   "id": 1996284807462067036
}
```
`id`: int 

## /api/v1/expressions
Receive all expressions from orchestrator store
- Currently, does not support pagination, just returns all expressions kept in the store

### Request
```http request
GET http://localhost:8080/api/v1/expressions
```

### Response
```json
{
   "expressions": [
      {
         "id": 1996284807462067036,
         "status": "completed",
         "result": 19.0
      },
      {
         "id": 1798228190132811771,
         "status": "pending",
         "result": 0
      }
   ]
}
```
`id`: int

`status`: string

`result`: double

## /api/v1/expressions/{id}
Receive specific expression by id

### Request
```http request
GET http://localhost:8080/api/v1/expressions/123
```

### Response
```json
{
   "expression": {
      "id": 1996284807462067036, 
      "status": "completed", 
      "result": 19.0
   }
}
```
`id`: int

`status`: string

`result`: double

# Future Plans

- Implement durable task queue using message broker (RabbitMQ, NATS, Redis, etc.)
- Implement persistent data storage using database (PostgreSQL, SQLite, MySQL, etc.)
- Implement pagination