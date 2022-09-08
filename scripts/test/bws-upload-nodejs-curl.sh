cd $(dirname $0)

zip -r nodejs nodejs-project -x \*/.git/\* \*/nodejs_modules/\*

docker run --rm --network=application-network -v $PWD/nodejs.zip:/nodejs.zip curlimages/curl -H "Host:secure.niij.fml.org" -X POST -F user_name=totto -F application_name=nodejs-project -F runtime=nodejs -F application_file=@nodejs.zip nginx:80/api/1/control-panel-backend/upload
