FROM golang:1.17-alpine AS build

WORKDIR /build

ADD . .

RUN go build -o user-service-app ./cmd/app
RUN go build -o user-service-migration ./cmd/migration

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/user-service-app .
COPY --from=build /build/user-service-migration .

CMD ["sh", "-c", "user-service-migration && user-service-app"]
