FROM golang:1.17-alpine AS build

WORKDIR /build

ADD . .

RUN go build -o api-gateway-app

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/api-gateway-app .

CMD ["api-gateway-app"]
