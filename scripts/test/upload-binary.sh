cd $(dirname $0)

docker run --rm --network=application-network -v $PWD/serverx86-64:/serverx86-64 curlimages/curl -H "Host:secure.niij.fml.org" -X POST -F user_name=totto -F application_name=binary -F runtime=binary -F application_file=@/serverx86-64 nginx:80/api/1/control-panel-backend/upload
