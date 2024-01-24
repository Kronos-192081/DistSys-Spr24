#!/bin/bash
docker compose up -d lb
printf "\nLoad balancer up and running at http://localhost:5000/home"