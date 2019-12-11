package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"

	h "github.com/fmcarrero/hex-microservice/api"
	"github.com/fmcarrero/hex-microservice/config"
	rr "github.com/fmcarrero/hex-microservice/repository/redis"

	"github.com/fmcarrero/hex-microservice/shortener"
	"github.com/gobuffalo/packr"
)

// https://www.google.com -> 98sj1-293
// http://localhost:8000/98sj1-293 -> https://www.google.com

// repo <- service -> serializer  -> http

func main() {

	var cfg config.Configuration

	readFile(&cfg)
	readEnv(&cfg)
	repo := chooseRepo(&cfg)
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :" + cfg.Server.Port)
		errs <- http.ListenAndServe(":"+cfg.Server.Port, r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)

}

func chooseRepo(cfg *config.Configuration) shortener.RedirectRepository {
	host := fmt.Sprintf("redis://%s:%s", cfg.Database.Host, cfg.Database.Port)
	repo, err := rr.NewRedisRepository(host)
	if err != nil {
		log.Fatal(err)
	}
	return repo

}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *config.Configuration) {
	box := packr.NewBox("./config")
	f, err := box.Find("config-local.yml")
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(f))
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *config.Configuration) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
