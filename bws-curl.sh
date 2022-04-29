source ./.env

docker run --rm --network=${USER}-application-network curlimages/curl $1