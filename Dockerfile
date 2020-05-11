FROM golang:1.14.2 as builder
WORKDIR /go/src/github.com/furiatona/pimock/
COPY ./*.go ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pimock *.go 

FROM alpine:3.10.3
WORKDIR app
COPY --from=builder /go/src/github.com/furiatona/pimock/pimock ./
COPY responses/ responses/
EXPOSE 8080
ENTRYPOINT ["./pimock"]
