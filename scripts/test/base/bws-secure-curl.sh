source ./.env

docker run --rm --network=application-network curlimages/curl -H "Host:secure.niij.fml.org" nginx:80$1
