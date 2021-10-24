# Book Library Microservice

## Problem Definition

Book Library Microservice is an online system to manage library operational: managing book catalog, book lending, and
library membership.

### Roles

There are three user roles of this system: Admin, Librarian, and Member.

#### Admin

- Create, read, update, and delete librarian data.

#### Librarian

- Manage book data:
    - Create, read, update, and delete book data
    - Update book stock data
- Create, read, update, and delete member data.
- Create, read, and update lending data by all member.

#### Member

- Read book & book stock data.
- Create book lending data.

## Solution Details

### Architecture

This repository consists of several microservices based on the domain:

- API Gateway
- Book Service
- Lending Service
- User Service

### Tech stack

- Go 1.17
- GraphQL with JWT auth
- gRPC
- MongoDB
- Docker

## Tutorial

1. Build all services and database

``` bash
make env && make run-docker
```

2. Go to [http://localhost:8000](http://localhost:8000) to start request from GraphQL playground.

3. Query example:

    - [User domain query](https://graphqlbin.com/v2/zqzzUw)
    - [Book domain query](https://graphqlbin.com/v2/ypyBfN)
    - [Lending domain query](https://graphqlbin.com/v2/p9mvHO)

## Author

Muhammad Habibullah, 2021

[muhammadhabibullah.id@gmail.com](mailto:muhammadhabibullah.id@gmail.com)
