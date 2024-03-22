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

# Load Balancer Implementation

This README provides an overview of the load balancer implementation.

## Environment Variables

- **`DATABASE_URL`**: The URL for connecting to the database.
- **`SERVER_IMAGE`**: The Docker image used for creating server containers.
- **`NETWORK_MODE`**: The Docker network mode for communication between containers.
- **`SERVER_SCHEMA`**: The schema used for the server configuration.
- **`HASH_MODULO`**: The modulo value used for consistent hashing.

## Initialization

1. **Database Connection**: The load balancer establishes a connection to the specified database using the `DATABASE_URL` environment variable.
   
2. **Consistent Hashing Initialization**: The load balancer initializes a consistent hashing mechanism using the modulo value specified in the `HASH_MODULO` environment variable.
   
3. **Server Configuration**: The load balancer configures server instances using the Docker image specified in the `SERVER_IMAGE` environment variable and assigns them to the specified network mode (`NETWORK_MODE`).

## Data Structures

1. **Consistent Hashing Ring**: Utilized for distributing data across multiple server instances.
   
2. **Server List**: A map data structure to store the list of active server instances.

## Load Balancer Logic

- **`/rm` Endpoint (DELETE)**:
  - Removes specified server containers based on payload data.
  - Adjusts the number of servers to match the specified count.

- **`/read` Endpoint (POST)**:
  - Reads data from server containers based on the provided payload.
  - Performs range parsing to obtain shard ID list.

- **`/write` Endpoint (POST)**:
  - Writes data to server containers based on the payload.
  - Distributes data entries across appropriate server shards.
  - Implements fault-tolerance handling for failed writes.

- **Update Endpoint (PUT)**:
  - Updates data on server containers based on payload.
  - Implements fault-tolerance for partially successful updates.

- **Del Endpoint (DELETE)**:
  - Deletes data entries from server containers based on payload.
  - Handles fault-tolerance for incomplete deletions.

## Fault-Tolerance and Scalability

- **Server Heartbeat**:
  - Periodically checks the health of server instances.
  - Handles server failures by removing and respawning containers.

- **Add Server Container**:
  - Dynamically adds new server containers based on demand.
  - Configures the new server container with the specified schema and shard assignments.

- **Kill Server Container**:
  - Terminates existing server containers upon failure or removal.
  - Updates consistent hashing and server list accordingly.

## Utility Functions

1. **Consistent Hashing Functions**:
   - `hash_function()`: Computes the hash value for a given key using the specified modulo value.
   - `get_shard_id()`: Determines the shard ID for a given key based on the consistent hashing mechanism.

## Docker API Functions

1. **Container Management Functions**:
   - `create_container()`: Creates a new Docker container based on the specified image and network mode.
   - `remove_container()`: Removes the specified Docker container from the system.

## Other Handler Functions

### Read Handler

- **Endpoint**: `/read` (POST)
- **Description**: Handles read operations by retrieving data from the appropriate server shards based on the shard ID list parsed from the payload.
- **Functionality**:
  1. Parses the payload to extract the list of shard IDs relevant to the read operation.
  2. Sends read requests to the corresponding server instances based on the shard IDs.
  3. Aggregates the responses from multiple servers and returns the combined result to the client.

### Write Handler

- **Endpoint**: `/write` (POST)
- **Description**: Handles write operations by distributing data entries across the server shards using consistent hashing.
- **Functionality**:
  1. Computes the shard ID for the provided key using consistent hashing.
  2. Sends write requests to the server instances responsible for the identified shard IDs.
  3. Implements fault tolerance for partial writes by retrying failed writes on other server instances.

### Update Handler

- **Endpoint**: `/update` (PUT)
- **Description**: Handles update operations by updating data entries on the server shards and implementing fault-tolerance for partially successful updates.
- **Functionality**:
  1. Retrieves the existing data associated with the provided key.
  2. Updates the data value and sends update requests to the server instances responsible for the corresponding shard IDs.
  3. Ensures fault tolerance by retrying failed updates on other server instances.

### Delete Handler

- **Endpoint**: `/delete` (DELETE)
- **Description**: Handles delete operations by removing data entries from the server shards and ensuring fault-tolerance for incomplete deletions.
- **Functionality**:
  1. Sends delete requests to the server instances responsible for the shard IDs associated with the provided key.
  2. Verifies the success of the delete operation and retries on other server instances if necessary to ensure fault tolerance.

## Dockerfile Explaination

This provides information about the Dockerfile used for building and deploying the load balancer.

## Dockerfile Structure

1. **Base Image**:
   - `FROM golang:latest AS build`: Specifies the base image as the latest version of Golang.

2. **Build Stage**:
   - `WORKDIR /lb`: Sets the working directory inside the container to `/lb`.
   - `COPY lb.go ./`: Copies the main Go file (`lb.go`) into the working directory.
   - `COPY conhash ./conhash`: Copies the `conhash` directory containing consistent hashing library.
   - `RUN go mod init distri-lb`: Initializes a Go module named `distri-lb`.
   - `RUN go mod edit -replace prakhar/conhash=./conhash`: Replaces the standard `conhash` module with the local directory `./conhash`.
   - `RUN go mod tidy && go mod download`: Cleans up and downloads module dependencies.
   - `RUN CGO_ENABLED=0 GOOS=linux go build -o distri-lb`: Builds the executable binary named `distri-lb` with Linux OS target.

3. **Deployment Stage**:
   - `FROM postgres:latest AS deploy`: Sets the base image for deployment as the latest version of PostgreSQL.
   - `COPY --from=build /lb/distri-lb ./`: Copies the compiled binary `distri-lb` from the build stage.
   - `COPY init.sh /docker-entrypoint-initdb.d/`: Copies the initialization script to set up the database.

4. **Environment Variables**:
   - `POSTGRES_USER`: Specifies the username for PostgreSQL database (default: `postgres`).
   - `POSTGRES_PASSWORD`: Specifies the password for PostgreSQL database (default: `20CS30061`).
   - `POSTGRES_DB`: Specifies the name of the PostgreSQL database (default: `testdb`).

5. **Exposed Port**:
   - `EXPOSE 5000`: Exposes port 5000 for communication with the load balancer.

# Performance Analysis of Distributed Database

This performance analysis provides a detailed overview of the efficiency of the developed distributed database across multiple test scenarios, focusing on various configurations and their impact on write and read operations.

## Test Results

### Test 1: 
- **Configuration**:
  - Shard Replicas: 3
  - Shards: 4
  - Servers: 6
- **Init Time**: Completed successfully
- **Write Time**: 18.27 seconds
- **Write Speed**: 547.42 writes per second
- **Read Time**: 7.06 seconds
- **Read Speed**: 1416.94 reads per second
- **Analysis Time**: 25.33 seconds

### Test 2: 
- **Configuration**:
  - Shard Replicas: 7
  - Shards: 4
  - Servers: 7
- **Init Time**: Completed successfully
- **Write Time**: 35.67 seconds
- **Write Speed**: 280.36 writes per second
- **Read Time**: 6.83 seconds
- **Read Speed**: 1463.36 reads per second
- **Analysis Time**: 42.50 seconds

### Test 3: 
- **Configuration**:
  - Shard Replicas: 8
  - Shards: 6
  - Servers: 10
- **Init Time**: Completed successfully
- **Write Time**: 32.18 seconds
- **Write Speed**: 310.78 writes per second
- **Read Time**: 6.67 seconds
- **Read Speed**: 1498.75 reads per second
- **Analysis Time**: 38.85 seconds

### Test 4:

- All endpoints of the load balancer and distributed database system were found to be functioning correctly during testing.
- The fault tolerance mechanism of the system was successfully verified, with the load balancer seamlessly handling the failure of a server container by spawning a new container and copying shard entries from other replicas.

The manual testing and fault tolerance verification confirmed the robustness and reliability of the load balancer and distributed database system under normal and failure conditions.


## Detailed Analysis

### Write Operations:
- The write time varies across different configurations, ranging from 18.27 seconds to 35.67 seconds. 
- Test 1, with fewer shard replicas and servers, demonstrates higher write speeds compared to Test 2 and Test 3.
- As the number of shard replicas increases in Test 2 and Test 3, the write speed decreases due to increased synchronization overhead.

### Read Operations:
- Read times across all tests are relatively consistent, with Test 3 exhibiting the lowest read time of 6.67 seconds.
- Test 1 shows a slightly higher read time compared to Test 2 and Test 3, potentially due to fewer servers distributing the read load.

### Impact of Configuration:
- The choice of shard replicas, shards, and servers significantly influences the overall performance of the distributed database.
- Increasing the number of servers can potentially improve parallelism and distribution, leading to higher performance in terms of both write and read operations.
- However, an increase in shard replicas may introduce additional synchronization overhead, impacting overall performance.

## Conclusion

The performance analysis highlights the importance of carefully selecting the configuration parameters for a distributed database system. While higher replication and sharding can enhance fault tolerance and scalability, they also introduce overhead that can affect performance. Overall, understanding the trade-offs and optimizing the configuration based on specific use cases is essential for achieving optimal performance in a distributed database environment.

