#!/bin/bash
docker compose build server
docker compose build lb
echo "Server and load balancer images built. Execute 'make run' to start the system.\n"