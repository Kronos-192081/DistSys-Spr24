#!/bin/bash
docker compose up -d lb
echo "Load balancer up and running at http://localhost:5000/home"