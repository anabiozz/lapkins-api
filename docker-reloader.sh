#!/bin/bash
cd /root/lapkins; docker-compose stop && docker-compose pull && docker-compose up -d --build;