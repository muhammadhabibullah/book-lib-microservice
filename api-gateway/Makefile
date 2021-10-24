env:
	cp .env.example .env &
	cp .docker.env.example .docker.env

build-app:
	go build -o api-gateway-app

run-app: build-app
	./api-gateway-app
