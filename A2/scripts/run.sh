#!/bin/bash
docker compose up -d lb
echo -e "\nLoad balancer up and running at http://localhost:5000/. Execute 'make test' to initialize servers and perform testing."