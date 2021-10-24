env:
	cp .env.example .env &
	cp .docker.env.example .docker.env

build-app:
	go build -o bin/book-service-app ./cmd/app

build-migration:
	go build -o bin/book-service-migration ./cmd/migration

run-app: build-app
	bin/book-service-app

run-migration: build-migration
	bin/book-service-migration
