cd $(dirname $0)

docker stop totto-binary
docker stop totto-nodejs-project
docker stop totto-html-project

rm nodejs.zip
rm html.zip
