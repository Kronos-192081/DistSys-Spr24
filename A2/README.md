# <div align="center">Scalable Sharded Database</div>

> Assignment 2 of Distributed Systems course (CS60002) offered in Spring semester 2024, Department of CSE, IIT Kharagpur.

<!-- PROJECT LOGO -->
<div align="center">
  <p align="center">
    <br />
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
    <img src="https://img.shields.io/badge/Python-FFD43B?style=for-the-badge&logo=python&logoColor=blue"/>
    <img src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white">
    <img src="https://img.shields.io/badge/-Bash-1f425f.svg?style=for-the-badge&logo=image%2Fpng%3Bbase64%2CiVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyZpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw%2FeHBhY2tldCBiZWdpbj0i77u%2FIiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8%2BIDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuNi1jMTExIDc5LjE1ODMyNSwgMjAxNS8wOS8xMC0wMToxMDoyMCAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTUgKFdpbmRvd3MpIiB4bXBNTTpJbnN0YW5jZUlEPSJ4bXAuaWlkOkE3MDg2QTAyQUZCMzExRTVBMkQxRDMzMkJDMUQ4RDk3IiB4bXBNTTpEb2N1bWVudElEPSJ4bXAuZGlkOkE3MDg2QTAzQUZCMzExRTVBMkQxRDMzMkJDMUQ4RDk3Ij4gPHhtcE1NOkRlcml2ZWRGcm9tIHN0UmVmOmluc3RhbmNlSUQ9InhtcC5paWQ6QTcwODZBMDBBRkIzMTFFNUEyRDFEMzMyQkMxRDhEOTciIHN0UmVmOmRvY3VtZW50SUQ9InhtcC5kaWQ6QTcwODZBMDFBRkIzMTFFNUEyRDFEMzMyQkMxRDhEOTciLz4gPC9yZGY6RGVzY3JpcHRpb24%2BIDwvcmRmOlJERj4gPC94OnhtcG1ldGE%2BIDw%2FeHBhY2tldCBlbmQ9InIiPz6lm45hAAADkklEQVR42qyVa0yTVxzGn7d9Wy03MS2ii8s%2BeokYNQSVhCzOjXZOFNF4jx%2BMRmPUMEUEqVG36jo2thizLSQSMd4N8ZoQ8RKjJtooaCpK6ZoCtRXKpRempbTv5ey83bhkAUphz8fznvP8znn%2B%2F3NeEEJgNBoRRSmz0ub%2FfuxEacBg%2FDmYtiCjgo5NG2mBXq%2BH5I1ogMRk9Zbd%2BQU2e1ML6VPLOyf5tvBQ8yT1lG10imxsABm7SLs898GTpyYynEzP60hO3trHDKvMigUwdeaceacqzp7nOI4n0SSIIjl36ao4Z356OV07fSQAk6xJ3XGg%2BLCr1d1OYlVHp4eUHPnerU79ZA%2F1kuv1JQMAg%2BE4O2P23EumF3VkvHprsZKMzKwbRUXFEyTvSIEmTVbrysp%2BWr8wfQHGK6WChVa3bKUmdWou%2BjpArdGkzZ41c1zG%2Fu5uGH4swzd561F%2BuhIT4%2BLnSuPsv9%2BJKIpjNr9dXYOyk7%2FBZrcjIT4eCnoKgedJP4BEqhG77E3NKP31FO7cfQA5K0dSYuLgz2TwCWJSOBzG6crzKK%2BohNfni%2Bx6OMUMMNe%2Fgf7ocbw0v0acKg6J8Ql0q%2BT%2FAXR5PNi5dz9c71upuQqCKFAD%2BYhrZLEAmpodaHO3Qy6TI3NhBpbrshGtOWKOSMYwYGQM8nJzoFJNxP2HjyIQho4PewK6hBktoDcUwtIln4PjOWzflQ%2Be5yl0yCCYgYikTclGlxadio%2BBQCSiW1UXoVGrKYwH4RgMrjU1HAB4vR6LzWYfFUCKxfS8Ftk5qxHoCUQAUkRJaSEokkV6Y%2F%2BJUOC4hn6A39NVXVBYeNP8piH6HeA4fPbpdBQV5KOx0QaL1YppX3Jgk0TwH2Vg6S3u%2BdB91%2B%2FpuNYPYFl5uP5V7ZqvsrX7jxqMXR6ff3gCQSTzFI0a1TX3wIs8ul%2Bq4HuWAAiM39vhOuR1O1fQ2gT%2F26Z8Z5vrl2OHi9OXZn995nLV9aFfS6UC9JeJPfuK0NBohWpCHMSAAsFe74WWP%2BvT25wtP9Bpob6uGqqyDnOtaeumjRu%2ByFu36VntK%2FPA5umTJeUtPWZSU9BCgud661odVp3DZtkc7AnYR33RRC708PrVi1larW7XwZIjLnd7R6SgSqWSNjU1B3F72pz5TZbXmX5vV81Yb7Lg7XT%2FUXriu8XLVqw6c6XqWnBKiiYU%2BMt3wWF7u7i91XlSEITwSAZ%2FCzAAHsJVbwXYFFEAAAAASUVORK5CYII%3D">
    <br />
    <br />
    <a href="https://github.com/Kronos-192081/DistSys-Spr24/issues">Report Bug</a>
    ·
    <a href="https://github.com/Kronos-192081/DistSys-Spr24/issues">Request Feature</a>
  </p>
</div>


# About The Project

This project is an implementation of a load balanced sharded database system in Go that utilizes consistent hashing. The problem statement can be viewed [here](./DS_assign1_LB_2024.pdf).

Team Members:
- Kartik Pontula (20CS10031)
- Prakhar Singh (20CS10045)
- Shiladitya De (20CS30061)
- Sourabh Soumyakanta Das (20CS30051)


# Server Implementation

The server facilitates various operations related to a distributed database system, including configuration management, data manipulation, and interaction via HTTP endpoints.

## Implementation Details

### Main Function (`main`)

- Initializes the HTTP server and sets up routes for different endpoints.
- Establishes a connection to the PostgreSQL database using the provided credentials.
- Starts the server to listen for incoming HTTP requests on the specified port.

### HTTP Handlers

#### Home Handler (`homeHandler`)

- Handles requests to the `/home` endpoint.
- Returns a JSON response with a greeting message including the server number obtained from the environment variable.

#### Heartbeat Handler (`heartbeatHandler`)

- Handles requests to the `/heartbeat` endpoint.
- Responds with an HTTP status OK (200) to indicate the server's availability.

#### Config Handler (`configHandler`)

- Manages configuration requests to the `/config` endpoint.
- Decodes incoming JSON payload containing database configuration details.
- Sets up database shards and initializes corresponding tables based on the provided configuration.
- Returns a JSON response indicating the success or failure of the configuration process.

#### Read Handler (`readHandler`)

- Processes requests to the `/read` endpoint for retrieving data from a specified shard within a given student ID range.
- Decodes incoming JSON payload containing shard name and student ID range.
- Executes a SQL query to fetch data from the specified shard and constructs a JSON response.

#### Copy Handler (`copyHandler`)

- Handles requests to the `/copy` endpoint for copying data from multiple shards.
- Decodes incoming JSON payload specifying shards to copy from.
- Queries the database for data from each shard and compiles the results into a JSON response.

#### Write Handler (`writeHandler`)

- Manages requests to the `/write` endpoint for inserting new data entries into a specified shard.
- Decodes incoming JSON payload containing shard name, current index, and data to be inserted.
- Constructs SQL `INSERT` queries and executes them to add new data entries to the database.

#### Update Handler (`updateHandler`)

- Processes requests to the `/update` endpoint for updating existing data entries in a specified shard.
- Decodes incoming JSON payload containing shard name, student ID, and updated data.
- Constructs SQL `UPDATE` query and executes it to update the corresponding data entry in the database.

#### Delete Handler (`deleteHandler`)

- Handles requests to the `/del` endpoint for deleting data entries from a specified shard.
- Decodes incoming JSON payload containing shard name and student ID.
- Constructs SQL `DELETE` query and executes it to remove the specified data entry from the database.

### Database Setup Function (`dbSetup`)

- Configures the database based on the provided configuration (`dbConfig`).
- Creates tables for each shard according to the specified schema.
- Returns a boolean value indicating the success or failure of the setup process.

### Database Models

- Defines Go structs representing database rows (`Row`) and configuration details (`dbSchema`, `dbConfig`).

### HTTP Server Configuration

- Sets up HTTP routes, listens on the specified port, and handles server startup errors.

## Dependencies

- `database/sql`: Standard Go package for working with SQL databases.
- `github.com/lib/pq`: PostgreSQL driver for Go.

## Environment Variables

- `PORT`: Specifies the port number on which the server listens for incoming HTTP requests.
- `SERVER_NUMBER`: Identifies the server number for differentiating responses in a distributed environment.

## Dockerfile Explaination

This Dockerfile defines the build and deploy stages for a distributed database server implemented in Go. It uses multi-stage builds to separate the build environment from the final deployment environment.

### Build Stage

The build stage (`FROM golang:latest AS build`) sets up the environment for compiling the Go code.

1. **Work Directory**: Sets the working directory to `/server`.
2. **Copy Source Code**: Copies the `main.go` file into the working directory.
3. **Go Modules Setup**:
   - Initializes a Go module named `distri-server`.
   - Tidies up and downloads module dependencies using `go mod tidy` and `go mod download`.
4. **Build Binary**:
   - Sets `CGO_ENABLED=0` and `GOOS=linux` to compile a statically linked binary for Linux.
   - Builds the server binary named `distri-server`.

### Deploy Stage

The deploy stage (`FROM postgres:latest AS deploy`) defines the environment for deploying the server along with a PostgreSQL database.

1. **Copy Binary**:
   - Copies the compiled `distri-server` binary from the build stage into the deploy stage.
2. **Copy Initialization Script**:
   - Copies an initialization script (`init.sh`) into the PostgreSQL image's `docker-entrypoint-initdb.d` directory.
   - This script is used to set up the database schema and perform any initializations required.
3. **Environment Variables**:
   - Sets environment variables for the PostgreSQL database (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`).
4. **Expose Port**:
   - Exposes port 5000 to allow external communication with the server.

<b> Design Choice: </b>The Dockerfile employs a multi-stage build approach to optimize the image size and improve security. By using separate build and deploy stages, it ensures that only necessary dependencies are included in the final deployment image, resulting in a smaller footprint. The build stage utilizes the official Golang image to compile the server binary, while the deploy stage utilizes the official PostgreSQL image to set up the database environment. Additionally, the Dockerfile exposes port 5000 to allow external communication with the server. Overall, this design choice streamlines the Docker image creation process, enhances portability, and promotes consistency across different deployment environments.

# Consistent Hashing Implementation

## Overview

- Consistent hashing is implemented in Go using a custom `ConHash` structure.
- The implementation includes functionalities to add servers, remove servers, and allocate servers based on client IDs.
- Implemented as a go module which can be imported.

## Data Structures

### 1. Node

- Represents a node in the `ConHash` structure.
- Has two attributes:
  - `Occ` (Occupancy): Indicates whether the node is occupied or not.
  - `Name`: Represents the name of the server.

### 2. ConHash

- Consistent hashing structure.
- Attributes:
  - `HashD`: Array of nodes representing the hash ring.
  - `Size`: Size of the hash ring.
  - `VirtServ`: Number of virtual servers per physical server.
  - `Nserv`: Number of servers in the hash ring.
  - `AllServers`: Map to track all server names.
  - `ServerID`: Map to track server IDs.

## Functional Endpoints

### 1. Add Servers
<b>Note: </b> Linear probing is used to resolve conflicts.
- **Method:** `Add(ids []int, Names []string) int`
  - Adds servers to the consistent hash ring.
  - Takes an array of server IDs and corresponding names.
  - Checks for name uniqueness.
  - Checks if the size limit is exceeded.
  - Returns 1 on success, 0 on failure.

### 2. Get Configuration (For debugging purposes)

- **Method:** `GetConfig()`
  - Prints the configuration of the consistent hash ring.
  - Displays index, status, and server information.

### 3. Add Single Server

- **Method:** `AddServer(id int, Name string) int`
  - Adds a single server to the consistent hash ring.
  - Checks for name uniqueness.
  - Checks if the size limit is exceeded.
  - Returns 1 on success, 0 on failure.

### 4. Remove Server
<b> Time Complexity:</b>  &nbsp; $\mathcal{O(K)}$ time complexity, where $\mathcal{K}$ is the number of virtual servers. 
- **Method:** `RemoveServer(Name string) int`
  - Removes a server from the consistent hash ring.
  - Returns 1 on success, 0 on failure.

### 5. Get Server Allocation

- **Method:** `GetServer(id int) string`
  - Returns the server for the given client ID.
  - Handles the case of no allocable server.
  - Prints the hash for the given client ID.

## Usage

- Initialize a new `ConHash` instance using `NewConHash(m, k)` with the desired size and virtual servers.
- Use the provided methods to interact with the consistent hash ring.

## Example

```go
package main

import (
    "fmt"
    "prakhar/conhash"
)

// Create a new ConHash instance
ch := conhash.NewConHash(100, 3)

// Add servers to the hash ring
ch.Add([]int{1, 2, 3}, []string{"server1", "server2", "server3"})

// Get and print the configuration
ch.GetConfig()

// Add a single server
ch.AddServer(4, "server4")

// Remove a server
ch.RemoveServer("server2")

// Get server allocation for a client ID
server := ch.GetServer(123)
fmt.Println("Allocated Server:", server)
```