# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /go/src/github.com/fmcarrero/hex-microservices/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goapp 

# final stage
FROM alpine:3.7
COPY --from=builder /go/src/github.com/fmcarrero/hex-microservices/goapp .
EXPOSE 8000
ENTRYPOINT ["./goapp"]