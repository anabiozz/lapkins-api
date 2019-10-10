#!/bin/bash
cd /root/lapkins-api; docker-compose stop && docker-compose pull && docker-compose up -d --build;