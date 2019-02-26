# Stage 1: Build executable
FROM golang:1.11 as buildImage

# We start with migrate so this gets cached most of the time
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/steffenmllr/prometheus-webpagetest-exporter
COPY . .

RUN dep ensure
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build

# Stage 2: Create release image
FROM alpine:latest

RUN mkdir app
WORKDIR app

COPY --from=buildImage /go/src/github.com/steffenmllr/prometheus-webpagetest-exporter/build exporter
EXPOSE 3030
ENTRYPOINT [ "/app/exporter" ]
