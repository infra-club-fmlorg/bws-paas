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

test-reset:
	sh $(tests)/bws-test-reset.sh

upload-binary:
	sh $(tests)/bws-upload-curl.sh

curl-binary:
	sh $(tests)/bws-test-curl.sh

archive-nodejs:
	sh $(tests)/bws-archive-nodejs.sh

upload-nodejs: archive-nodejs
	sh $(tests)/bws-upload-nodejs-curl.sh

curl-nodejs:
	sh $(tests)/bws-test-nodejs-curl.sh

archive-html:
	sh $(tests)/bws-archive-html.sh

upload-html: archive-html
	sh $(tests)/bws-upload-html-curl.sh
	
curl-html:
	sh $(tests)/bws-test-html-curl.sh
