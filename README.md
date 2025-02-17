This is distributed calculator
(Or a bad Apache Spark cosplay?)

# Table of Contents
1. [Start Up](#start-up)
   - [Command Line](#command-line)
   - [Taskfile](#taskfile)
   - [Docker Compose](#docker-compose)
   - [Docker CLI](#docker-cli)
2. [Services](#services)
   - [Orchestrator](#orchestrator)
   - [Agent](#agent)
  


# Start Up
You can run calculation cluster in several ways

## Command Line
Though it is advised to use Docker Compose to run app, you can still use console commands to run it

## Taskfile
Also you can use Taskfile to run app 

## Docker Compose

## Docker CLI

# Services

## Orchestrator
Orchestrator is a master node of calculation cluster

It decomposes the expression to run in parallel tasks on ***Agent*** instances

### Configuration
Orchestrator can be configured via environment variables

`HOST`: Host to run on (default: `0.0.0.0`)

NOTICE: Do not change host if you run in docker, otherwise it may not work properly 

`PORT`: Port to run on (default: `8080`)

## Agent
Agent is a worker node of calculation cluster

It uses long polling to receive tasks via ***Orchestrator*** API

NOTICE: On start up, agent will try to connect to orchestrator. It will exit immediately on failure

### Configuration
Agent can be configured via environment variables

`COMPUTING_POWER`: Amount of active workers per agent instance (default: `10`)

`POLL_TIMEOUT`: Polling interval in milliseconds (default: `50`)

`MAX_RETRIES`: Maximum number of retries on failed requests (default: `3`)

`MASTER_URL`: Orchestrator URL