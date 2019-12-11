package redis_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fmcarrero/hex-microservice/repository/redis"
	"github.com/fmcarrero/hex-microservice/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/teris-io/shortid"
	"github.com/testcontainers/testcontainers-go"
)

var ip string
var port string
var client shortener.RedirectRepository

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379/tcp"},
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer redisC.Terminate(ctx)
	ip, err = redisC.Host(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	redisPort, err := redisC.MappedPort(ctx, "6379/tcp")
	port = redisPort.Port()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	client, _ = redis.NewRedisRepository(fmt.Sprintf("redis://%s:%s", ip, port))
	m.Run()

}

func TestPost(t *testing.T) {

	redirect := &shortener.Redirect{
		Code:      shortid.MustGenerate(),
		URL:       shortid.MustGenerate(),
		CreatedAt: time.Now().UTC().Unix(),
	}

	errStore := client.Store(redirect)
	if errStore != nil {
		assert.Error(t, errStore)
	}

}

func TestFind(t *testing.T) {

	redirect := &shortener.Redirect{
		Code:      shortid.MustGenerate(),
		URL:       shortid.MustGenerate(),
		CreatedAt: time.Now().UTC().Unix(),
	}

	errStore := client.Store(redirect)
	if errStore != nil {
		assert.Error(t, errStore)
	}

	response, err := client.Find(redirect.Code)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, response, redirect)
}
