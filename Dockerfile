# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /go/src/github.com/fmcarrero/hex-microservices/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go get -u github.com/gobuffalo/packr/... 
RUN CGO_ENABLED=0 GOOS=linux packr build -a -installsuffix cgo -o  goapp

# final stage
FROM alpine:3.7

COPY --from=builder /go/src/github.com/fmcarrero/hex-microservices/goapp .
ENTRYPOINT ["./goapp"]