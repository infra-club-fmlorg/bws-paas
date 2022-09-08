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
	sh $(tests)/test-reset.sh

upload-binary:
	sh $(tests)/upload-binary.sh

curl-binary:
	sh $(tests)/test-binary.sh

upload-nodejs:
	sh $(tests)/upload-nodejs.sh

test-nodejs:
	sh $(tests)/test-nodejs.sh

upload-html:
	sh $(tests)/upload-html.sh
	
curl-html:
	sh $(tests)/test-html.sh
