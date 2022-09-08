scripts = scripts
tests = $(scripts)/test

setup:
	bash $(scripts)/setup.sh

build:
	docker-compose build

up:
	docker-compose up -d

stop:
	docker-compose stop

reset:
	sh $(scripts)/reset.sh

# 以下テスト用
dev: test-reset
	docker-compose up --build

dev-no-cache: test-reset
	docker-compose up --build --no-cache

test-reset:
	sh $(tests)/bws-test-reset.sh

upload-binary:
	sh $(tests)/bws-upload-curl.sh

curl-binary:
	sh $(tests)/bws-test-curl.sh

upload-nodejs:
	sh $(tests)/bws-upload-nodejs-curl.sh

test-nodejs:
	sh $(tests)/bws-test-nodejs-curl.sh

archive-html:
	sh $(tests)/bws-archive-html.sh

upload-html:
	sh $(tests)/bws-upload-html-curl.sh
	
curl-html:
	sh $(tests)/bws-test-html-curl.sh
