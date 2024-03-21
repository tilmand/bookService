FROM golang:latest AS builder

WORKDIR /go/src/app

COPY . .

RUN test -f go.mod || go mod init app
RUN go mod tidy

RUN go test ./...

RUN go build -o app .

FROM debian:buster AS mongodb

RUN apt-get update && apt-get install -y --allow-unauthenticated mongo-tools

FROM golang:latest

COPY --from=builder /go/src/app/app /usr/local/bin/app

COPY --from=mongodb /usr/bin/mongodump /usr/local/bin/mongodump
COPY --from=mongodb /usr/bin/mongorestore /usr/local/bin/mongorestore

CMD ["app"]
