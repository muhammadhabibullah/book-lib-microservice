FROM golang:1.15-alpine AS build

WORKDIR /build

ADD . .

RUN go build -o book-service ./cmd/app
RUN go build -o book-service-migration ./cmd/migration

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/book-service .
COPY --from=build /build/book-service-migration .

CMD ["sh", "-c", "book-service-migration && book-service"]
