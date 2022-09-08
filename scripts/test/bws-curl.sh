source ./.env

docker run --rm --network=application-network curlimages/curl $1
