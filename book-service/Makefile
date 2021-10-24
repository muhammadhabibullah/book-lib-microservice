env:
	cp .env.example .env &
	cp .docker.env.example .docker.env

build-app:
	go build -o book-service-app ./cmd/app

build-migration:
	go build -o book-service-migration ./cmd/migration

run-app: build-app
	./book-service-app

run-migration: build-migration
	./book-service-migration
