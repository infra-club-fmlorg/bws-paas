source ./.env

docker run --rm --network=${USER}-application-network curlimages/curl -H "Host:secure.niij.fml.org" $1