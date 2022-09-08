scripts = scripts
tests = $(scripts)/test

setup: reset
	bash $(scripts)/setup.sh

build: setup
	docker-compose build

up: setup
	docker-compose up -d

stop:
	docker-compose stop

reset: test-reset
	sh $(scripts)/reset.sh

# 以下テスト用
dev: setup
	docker-compose up --build up

dev-no-cache: setup
	docker-compose up --build --no-cache up

archive-nodejs:
	sh $(tests)/bws-archive-nodejs.sh

upload-binary:
	sh $(tests)/bws-upload-curl.sh

upload-nodejs: archive-nodejs
	sh $(tests)/bws-upload-nodejs-curl.sh

curl-binary:
	sh $(tests)/bws-test-curl.sh

curl-nodejs:
	sh $(tests)/bws-test-nodejs-curl.sh

test-reset:
	sh $(tests)/bws-test-reset.sh
