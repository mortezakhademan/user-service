# stage build service
FROM golang:1.24-alpine as build
RUN apk add --no-cache bash git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -mod vendor -o main .

# stage run service from builded executable
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/main .
CMD ["./main", "-env", "/.env"]