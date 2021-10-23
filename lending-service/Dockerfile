FROM golang:1.17-alpine AS build

WORKDIR /build

ADD . .

RUN go build -o lending-service-app ./cmd/app
RUN go build -o lending-service-migration ./cmd/migration

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/lending-service-app .
COPY --from=build /build/lending-service-migration .

CMD ["sh", "-c", "lending-service-migration && lending-service-app"]
