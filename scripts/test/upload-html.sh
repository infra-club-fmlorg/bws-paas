cd $(dirname $0)

zip -r html html-project -x \*/.git/\* \*/node_modules/\*

docker run --rm --network=application-network -v $PWD/html.zip:/html.zip curlimages/curl -H "Host:secure.niij.fml.org" -X POST -F user_name=totto -F application_name=html-project -F runtime=html -F application_file=@html.zip nginx:80/api/1/control-panel-backend/upload
