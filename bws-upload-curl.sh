docker run --rm --network=application-network -v $PWD/serverx86-64:/serverx86-64 curlimages/curl -H "Host:secure.niij.fml.org" -X POST -F userid=totto -F appname=hello -F userfile=@/serverx86-64 nginx:80/api/1/control-panel-backend/upload
