env:
	cp .env.example .env &
	cp .docker.env.example .docker.env

build-app:
	go build -o bin/lending-service-app ./cmd/app

build-migration:
	go build -o bin/lending-service-migration ./cmd/migration

run-app: build-app
	bin/lending-service-app

run-migration: build-migration
	bin/lending-service-migration
