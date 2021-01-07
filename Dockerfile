#FROM golang:1.11.3-alpine3.8 AS build
FROM golang:1.15-alpine3.12 AS build

COPY . /go/src/github.com/jpweber/cole

WORKDIR /go/src/github.com/jpweber/cole
RUN apk --update upgrade && \
    apk add ca-certificates && \
    update-ca-certificates && \
    apk add dep git 
RUN dep ensure -v
RUN CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o app .

# copy the binary from the build stage to the final stage
FROM alpine:3.8
RUN apk --update upgrade && \
    apk add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# For the case of HTTP insecure server, comment the below two lines
COPY ./cert/https-server.crt /usr/local/share/ca-certificates/https-server.crt
COPY ./cert/https-server.key /usr/local/share/ca-certificates/https-server.key

COPY --from=build /go/src/github.com/jpweber/cole/app /cole
CMD ["/cole"]
