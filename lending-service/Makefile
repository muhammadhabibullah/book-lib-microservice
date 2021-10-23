build-app:
	go build -o lending-service-app ./cmd/app

build-migration:
	go build -o lending-service-migration ./cmd/migration

run-app: build-app
	./lending-service-app

run-migration: build-migration
	./lending-service-migration
