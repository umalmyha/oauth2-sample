FROM golang:1.22.0-alpine AS build

WORKDIR /build
COPY . .

RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/server .

FROM alpine:latest

ENV HTTP_SERVER_PORT=8080

COPY --from=build /bin/server /bin/server

EXPOSE ${HTTP_SERVER_PORT}

WORKDIR /bin

CMD ["server"]