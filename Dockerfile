# gabe/Dockerfile

### Builder Image
FROM golang:latest AS builder
WORKDIR /go/src/github.com/dicurrio/gabe
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

### Application Image
FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/dicurrio/gabe/gabe .
CMD ["./gabe"]