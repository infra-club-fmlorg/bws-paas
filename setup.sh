#/bin/bash

source ./.env

docker volume create --name=${USER}-application-active
docker volume create --name=${USER}-application-incoming
docker network create ${USER}-application-network