# <div align="center">Customizable Load Balancer</div>

> Assignment 1 of Distributed Systems course (CS60002) offered in Spring semester 2024, Department of CSE, IIT Kharagpur.

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
    <a href="https://github.com/outer-rim/Query-Optimiser/issues">Report Bug</a>
    ·
    <a href="https://github.com/outer-rim/Query-Optimiser/issues">Request Feature</a>
  </p>
</div>


## About The Project

This project is a customisable load balancer.

Team Members:
- Kartik Pontula (20CS10031)
- Prakhar Singh (20CS10045)
- Shiladitya De (20CS30061)
- Sourabh Soumyakanta Das (20CS30051)

<!-- <p align="right">(<a href="#top">back to top</a>)</p> -->



<!-- GETTING STARTED -->
## Getting Started

..............

..............

.............


# Server Implementation

## 1. Overview

- Implemented in Go using the Echo web framework.
- Designed for simplicity with two main endpoints: `/home` and `/heartbeat`.

## 2. Functional Endpoints

### a. `/home`

- **Method:** GET
- **Description:** Returns a JSON response with a greeting message and server number from the `SERVER_NUMBER` environment variable (defaulting to an empty string).
- **Example Response:** `{ "message": "Hello from Server [ID]", "status": "successful" }`

### b. `/heartbeat`

- **Method:** GET
- **Description:** Health check endpoint returning a 200 OK response with no content.

## 3. Middleware

- Logger Middleware: Logs incoming requests and their details.
- Recover Middleware: Recovers from panics during request handling.

## 4. Configuration

- Environment Variables:
  - `PORT`: HTTP port on which the server listens (default: `5000`).
  - `SERVER_NUMBER`: Server number for the `/home` endpoint response (default: empty string).

## 5. Dockerfile Explanation

###  Build Stage:

- **FROM golang:latest AS build:**
  - Uses the official Golang image as the base image for the build stage.

- **WORKDIR /server:**
  - Sets the working directory to `/server`.

- **COPY go.mod go.sum ./ :**
  - Copies the Go module files to the working directory.

- **RUN go mod download:**
  - Downloads the Go modules.

- **COPY &nbsp; \*.go ./ :**
  - Copies the Go source code to the working directory.

- **RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-server:**
  - Builds the Go application with cross-compilation for Linux, disabling CGO, and outputs the executable as `/docker-server`.

### Deployment Stage:

- **FROM scratch AS deploy:**
  - Uses an empty scratch image as the base image for the deployment stage.

- **WORKDIR /server:**
  - Sets the working directory to `/server`.

- **COPY --from=build /docker-server /docker-server:**
  - Copies the built executable from the build stage to the deployment stage.

- **EXPOSE 5000:**
  - Exposes port 5000 for external connections.

- **ENTRYPOINT ["/docker-server"]:**
  - Sets the entry point for the container to run the `/docker-server` executable.

This Dockerfile follows a multi-stage build approach, optimizing the final image size by using a minimal scratch image for deployment. It first builds the Go application in one stage and then copies only the necessary artifacts to the final deployment image. The resulting container is lightweight and contains only the compiled executable and minimal dependencies.



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

- **Method:** `Add(ids []int, Names []string) int`
  - Adds servers to the consistent hash ring.
  - Takes an array of server IDs and corresponding names.
  - Checks for name uniqueness.
  - Checks if the size limit is exceeded.
  - Returns 1 on success, 0 on failure.

### 2. Get Configuration

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


# Load Balancer Performance Analysis

This document presents an analysis of the load balancer's performance based on the Python script provided. The script launches asynchronous requests to the load balancer under different scenarios and gathers metrics for load distribution and scalability.

## Experiment A-1: Load Distribution with N = 3

### Test Description:

- Send 10,000 asynchronous requests to a load balancer with N = 3 server containers.
- Report the request count handled by each server instance in a bar chart.

### Results:

**Bar Chart:** 

![Experiment A-1 Bar Chart](./Analysis/A1.png)

**Overall Metrics:**

- Time taken: 3.5 seconds (average)
- Memory usage: 2.5 GB

**Observations:**

- ......
- ..........

**Conclusion:**
- ..............

---

## Experiment A-2: Scalability with Incrementing N

### Test Description:

- Increment N from 2 to 6 and launch 10,000 requests on each increment.
- Report the average load of the servers at each run in a line chart.

### Results:

**Line Chart:** 

![Experiment A-2 Line Chart](./Analysis/A2.png)

**Overall Metrics:**

- Time taken: 25.5 seconds (average)
- Memory usage: 2.5 GB

**Observations:**

- ..............
- ..................

**Conclusion:**
- ..............

---
