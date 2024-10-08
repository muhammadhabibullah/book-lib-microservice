env:
	cp .docker-compose.env.example .docker-compose.env &
	cd ./api-gateway && $(MAKE) env &
	cd ./book-service && $(MAKE) env &
	cd ./lending-service && $(MAKE) env &
	cd ./user-service && $(MAKE) env

check-docker-env:
	docker-compose --env-file .docker-compose.env config -q

run-docker: check-docker-env
	docker-compose --env-file .docker-compose.env up --build -d $(SERVICE)

run-db-docker: check-docker-env
	docker-compose --env-file .docker-compose.env up --build -d mongo

down-docker: check-docker-env
	docker-compose --env-file .docker-compose.env down

stop-app-docker: check-docker-env
	docker-compose --env-file .docker-compose.env stop -t 5 api-gateway user-service book-service lending-service

stop-db-docker: check-docker-env
	docker-compose --env-file .docker-compose.env stop -t 5 mongo

run-migration-local:
	cd ./book-service && $(MAKE) run-migration &
	cd ./lending-service && $(MAKE) run-migration &
	cd ./user-service && $(MAKE) run-migration

build-app-local:
	cd ./api-gateway && $(MAKE) build-app &
	cd ./book-service && $(MAKE) build-app &
	cd ./lending-service && $(MAKE) build-app &
	cd ./user-service && $(MAKE) build-app

run-app-local: build-app-local
	cd ./api-gateway && $(MAKE) run-app &
	cd ./book-service && $(MAKE) run-app &
	cd ./lending-service && $(MAKE) run-app &
	cd ./user-service && $(MAKE) run-app
