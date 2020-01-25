FROM golang:latest

WORKDIR /go/src/github.com/dcb9/docker-cron/
COPY . .
RUN go get ./... \
    && CGO_ENABLED=0 GOOS=linux go build -o /usr/bin/cron

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=0 /usr/bin/cron /usr/bin/cron

ENTRYPOINT ["cron"]

