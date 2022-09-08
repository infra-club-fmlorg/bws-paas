cd $(dirname $0)

docker stop totto-hello
docker stop totto-vite-project
docker stop totto-html-project

# rm application.zip

exit 0
