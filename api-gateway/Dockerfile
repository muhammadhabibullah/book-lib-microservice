FROM golang:1.15-alpine AS build

WORKDIR /build

ADD . .

RUN go build -o api-gateway

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/api-gateway .

CMD ["api-gateway"]
