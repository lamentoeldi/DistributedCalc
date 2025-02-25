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
   - [Expressions](#expression)
4. [Examples of Use](#examples-of-use)
   - [/api/v1/calculate](#apiv1calculate)
   - [/api/v1/expressions](#apiv1expressions)
   - [/api/v1/expressions/{id}](#apiv1expressionsid)
5. [Future Plans](#future-plans)
  


# Start Up
You can run calculation cluster in several ways

## Command Line
Though it is advised to use Docker Compose to run app, you can still use console commands to run it

Use the following commands to download dependencies and run app with default configuration

```shell
cd backend
go mod download
```

```shell
go run cmd/orchestrator/main.go & go run cmd/agent/main.go
```

Use this to run frontend with default configuration

```shell
cd frontend
npm run build && npm run start
```

## Taskfile
Also you can use Taskfile to run app with default configuration

Use this to run backend with default configuration

```shell
task run-backend
```

Use this to run frontend with default configuration

```shell
task run-frontend
```

## Docker CLI
You can use Docker CLI to build images and then run containers

Use this to build images

```shell
docker build -t orchestrator:latest -f ./backend/build/package/orchestrator/Dockerfile ./backend & docker build -t agent:latest -f ./backend/build/package/agent/Dockerfile ./backend
```

Use this to run app with default configuration and forward orchestrator:8080 to localhost:8080

```shell
docker run -d --name orchestrator -p 8080:8080 orchestrator:latest && docker run -d --name agent --link orchestrator:orchestrator -e MASTER_URL=http://orchestrator:8080 agent:latest
```

Use this to build frontend image with default backend URL

```shell
docker build -t frontend:latest --build-arg BACKEND_URL=http://orchestrator:8080 ./frontend
```

Use this to run frontend with default configuration (run orchestrator first)

```shell
docker run -d --name frontend -p 3000:3000 --link orchestrator:orchestrator frontend:latest
```

## Docker Compose
Docker Compose is the most preferable way to run app. As mentioned in [compose file](docker-compose.yaml), on default 
the port 8080 of ***Orchestrator*** is bound on 8080 port of local machine and the port 3000 of ***Frontend*** is bound on 3000 port of local machine

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

It decomposes the expressionTable to run in parallel tasks on ***Agent*** instances

### Configuration
Orchestrator can be configured via environment variables

`HOST`: Host to run on (default: `0.0.0.0`)

NOTICE: Do not change host if you run in docker, otherwise it may not work properly 

`PORT`: Port to run on (default: `8080`), must be in range between `1` and `65535`

`LOG_LEVEL`: Level of logging (default: `info`)

`TIME_ADDITION_MS`: Time in milliseconds which `+` operation takes (default: `1`), must be non-negative integer

`TIME_SUBTRACTION_MS`: Time in milliseconds which `-` operation takes (default: `1`), must be non-negative integer

`TIME_MULTIPLICATION_MS`: Time in milliseconds which `*` operation takes (default: `1`), must be non-negative integer

`TIME_DIVISION_MS`: Time in milliseconds which `/` operation takes (default: `1`), must be non-negative integer

## Agent
Agent is a worker node of calculation cluster

It uses long polling to receive tasks via ***Orchestrator*** API

NOTICE: On start up, agent will try to connect to orchestrator. It will exit immediately on failure after retries

### Configuration
Agent can be configured via environment variables

`LOG_LEVEL`: Level of logging (default: `info`)

`COMPUTING_POWER`: Amount of active workers per agent instance (default: `10`), must be positive integer

`BUFFER_SIZE`: Size of task buffer (default: `128`), must be positive integer

`POLL_TIMEOUT`: Polling interval in milliseconds (default: `50`), must be positive integer

`MAX_RETRIES`: Maximum retries on failed requests (default: `3`), must be positive integer

`MASTER_URL`: Orchestrator URL in `protocol://host:port` format (default: `http://localhost:8080`)

## Frontend
Frontend is a web-interface for DistributedCalc

NOTICE: Due to Next.js specifics, `BACKEND_URL` build arg in `protocol://host:port` format must be provided

# Good to Know

## General

- Currently, the system keeps all data in-memory, that means that all data will be lost on restart
- Currently, the system is stateful, that means that data you receive depends on which node you have accessed
- Agents use long polling to receive tasks from orchestrator
- It is possible to use proxy like [envoy](https://www.envoyproxy.io), 
[nginx](https://nginx.org) or 
[traefik](https://doc.traefik.io/traefik/) to balance incoming requests between running nodes
- If result of expressions has more than `8` decimal places, they are thrown away
- Notice that expressions like `2 2 + 3` will be processed as `22+3` due to system design

## Expressions
1. During the evaluation, field `result` in Expressions schema is `0` until expressionTable is evaluated
2. May have several statuses:
   - `pending`: the expressionTable is being processed
   - `completed`: the expressionTable is processed and result is ready for use
   - `failed`: the system failed to process the expressionTable

# Examples of Use
- [API Specification](backend/api/v1/api.yaml)
- [More examples](backend/examples)

## /api/v1/calculate
Send expressionTable to start evaluation

### Request
```http request
POST localhost:8080/api/v1/calculate
Content-Type: application/json

{
  "expressionTable": "3*4+7"
}
```
`expressionTable`: string

### Response
```json
{
   "id": 1996284807462067036
}
```
`id`: int 

### Example

#### Success

```shell
curl -X POST "http://localhost:8080/api/v1/calculate" \
     -H "Content-Type: application/json" \
     -d '{"expressionTable": "2 + 2 * 2"}'
```

#### Bad Request

```shell
curl -X POST "http://localhost:8080/api/v1/calculate" \
     -H "Content-Type: application/json" \
     -d '"corrupted json"'
```

#### Unprocessable Entity

```shell
curl -X POST "http://localhost:8080/api/v1/calculate" \
     -H "Content-Type: application/json" \
     -d '{"expressionTable": "2++3"}'
```

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

### Example

```shell
curl -X GET "http://localhost:8080/api/v1/expressions"
```

## /api/v1/expressions/{id}
Receive specific expressionTable by id

### Request
```http request
GET http://localhost:8080/api/v1/expressions/1996284807462067036
```

### Response
```json
{
   "expressionTable": {
      "id": 1996284807462067036, 
      "status": "completed", 
      "result": 19.0
   }
}
```
`id`: int

`status`: string

`result`: double

### Example

#### Success or not found

```shell
curl -X GET "http://localhost:8080/api/v1/expressions/1996284807462067036"
```

#### Bad Request

```shell
curl -X GET "http://localhost:8080/api/v1/expressions/invalidpath"
```

## /internal/tasks

### GET

#### Request

```http request
GET http://localhost:8080/internal/task
```

#### Response

```json
{
   "task": {
      "id": 0,
      "arg1": 3.5,
      "arg2": 2,
      "operation": "+",
      "operation_time": 0.01
   }
}
```

#### Example

##### Success or not found

```shell
curl -X GET "http://localhost:8080/internal/task"
```

### POST

#### Request

```http request
POST http://localhost:8080/internal/task
Content-Type: application/json

{
    "id": 90913132,
    "result": 0.7
}
```

#### Example

##### Success or not found

````shell
curl -X POST "http://localhost:8080/internal/tasks" \
     -H "Content-Type: application/json" \
     -d '{"id": 999913183, "result": 0.7}'
````

##### Bad Request

```shell
curl -X POST "http://localhost:8080/internal/tasks" \
     -H "Content-Type: application/json" \
     -d '"corrupted json"'
```

# Future Plans

- Implement durable task queue using message broker (RabbitMQ, NATS, Redis, etc.)
- Implement persistent data storage using database (PostgreSQL, SQLite, MySQL, etc.)
- Implement pagination