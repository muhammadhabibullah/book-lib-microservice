env:
	cp .env.example .env &
	cp .docker.env.example .docker.env

build-app:
	go build -o bin/api-gateway-app

run-app: build-app
	bin/api-gateway-app
