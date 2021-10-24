env:
	cp .env.example .env &
	cp .docker.env.example .docker.env

build-app:
	go build -o bin/user-service-app ./cmd/app

build-migration:
	go build -o bin/user-service-migration ./cmd/migration

run-app: build-app
	bin/user-service-app

run-migration: build-migration
	bin/user-service-migration
