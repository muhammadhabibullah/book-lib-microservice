build-app:
	go build -o user-service-app ./cmd/app

build-migration:
	go build -o user-service-migration ./cmd/migration

run-app: build-app
	./user-service-app

run-migration: build-migration
	./user-service-migration
