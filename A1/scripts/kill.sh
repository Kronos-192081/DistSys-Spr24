#!/bin/bash
docker kill $(docker network inspect -f '{{range .Containers}}{{.Name}} {{end}}' net1)
docker compose down