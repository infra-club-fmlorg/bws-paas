FROM golang:1.18 AS BUILD

ENV GOPATH=

WORKDIR /go
COPY main.go /go
COPY go.mod /go

RUN go mod tidy && GOOS=linux go build -o app main.go

################
FROM ubuntu

COPY --from=BUILD /go/app /app
RUN chmod 111 /app

ENTRYPOINT [ "/app" ]
