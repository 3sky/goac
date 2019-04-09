FROM golang:alpine as builder
RUN mkdir /server
RUN apk add --no-cache git gcc libc-dev
ADD  server/ /server
WORKDIR /server
ENV GOBIN /go/bin
RUN go test -v ./...
RUN go build -o appServer ./...

FROM alpine
RUN apk --no-cache add ca-certificates
RUN adduser -S -D -H -h /app q
USER q
COPY --from=builder /server/appServer /app/appServer
WORKDIR /app
EXPOSE 5000
CMD ["./appServer"]