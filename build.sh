#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lapkins .
docker build -t anabiozz/lapkins .
docker push anabiozz/lapkins
ssh root@165.22.92.145 'bash -s' < docker-reloader.sh;