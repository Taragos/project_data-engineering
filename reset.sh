#!/bin/bash

docker compose down
docker volume rm project_data-engineering_clickhouse_data
docker volume rm project_data-engineering_kafka_data
docker volume rm project_data-engineering_postgres_data
docker volume rm project_data-engineering_minio_data

docker compose build 

docker compose up -d
