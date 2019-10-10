#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lapkins-api .
docker build -t anabiozz/lapkins-api .
docker push anabiozz/lapkins-api
ssh root@165.22.92.145 'bash -s' < docker-reloader.sh;