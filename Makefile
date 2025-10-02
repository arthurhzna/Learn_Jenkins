GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

watch-prepare: 
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

watch: 
	bin/air

build: 
	go build -o learn_jenkins

docker-compose: 
	docker-compose up -d --build --force-recreate

docker-build: 
	docker build --platform linux/amd64 -t arthurhozanna/learn_jenkins:$(tag) .

docker-push: 
	docker push arthurhozanna/learn_jenkins:$(tag)
