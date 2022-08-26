#/bin/bash

source ./.env

docker volume create --name=application-active
docker volume create --name=application-incoming
docker network create application-network
