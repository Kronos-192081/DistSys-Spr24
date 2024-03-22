#!/bin/bash
docker compose build server
docker compose build lb
printf "\nServer and load balancer images built. Execute 'make run' to start the load balancer system."