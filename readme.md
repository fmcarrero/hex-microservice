# hex-microservices

this is a simple example about management a ports and adapters architecture


### Tech


* [GOLANG] - GO!


### Installation

Dillinger requires [golang](https://golang.org/dl/)  to run.

Install the dependencies and devDependencies and start the server.

```sh
$ docker run --name my-redis-container -p 6379:6379 -d redis
$ docker run --name my-mongodb-container -p 27017:27017 -d mongo
$ export URL_DB=redis
$ export REDIS_URL =redis://localhost:6379
$ go run main.go
```
### Run with docker-compose

```sh
$ docker-compose up
