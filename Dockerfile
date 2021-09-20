## Multi-stage builds
FROM golang:1.17-alpine AS builder
WORKDIR /app
#here use absoluteDir pattern
COPY ./go.mod /app/
COPY ./src /app/src

RUN go mod tidy
RUN go mod verify

WORKDIR /app/src
# CGO has to be disabled for FROM scratch: CGO_ENABLED=0
# https://stackoverflow.com/questions/52640304/standard-init-linux-go190-exec-user-process-caused-no-such-file-or-directory
# https://stackoverflow.com/questions/62817082/how-does-cgo-enabled-affect-dynamic-vs-static-linking
# https://www.geeksforgeeks.org/static-and-dynamic-linking-in-operating-systems/
# Here use relativeDir pattern which binary file inside <WORKDIR/bookstore-oauth-api>
RUN CGO_ENABLED=0 go build -o /app/bookstore-oauth-api

## Deploy and run binary
FROM alpine:latest
WORKDIR /app
# Copied to the location /app/bookstore-oauth-api
COPY --from=builder /app/bookstore-oauth-api .
EXPOSE 8080
# excuted the binary inside /app/bookstore-oauth-api
ENTRYPOINT ["./bookstore-oauth-api"]