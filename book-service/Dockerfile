FROM golang:1.17-alpine AS build

WORKDIR /build

ADD . .

RUN go build -o book-service-app ./cmd/app
RUN go build -o book-service-migration ./cmd/migration

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/book-service-app .
COPY --from=build /build/book-service-migration .

CMD ["sh", "-c", "book-service-migration && book-service-app"]
