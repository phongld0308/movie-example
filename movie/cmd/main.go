package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/phongld0308/movie-example/movie/internal/controller/movie"
	metadatagateway "github.com/phongld0308/movie-example/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/phongld0308/movie-example/movie/internal/gateway/rating/grpc"
	httphandler "github.com/phongld0308/movie-example/movie/internal/handler/http"
	"github.com/phongld0308/movie-example/movie/internal/repository/postgres"
	"github.com/phongld0308/movie-example/pkg/discovery"
	"github.com/phongld0308/movie-example/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)

	// Get configuration from environment
	consulAddr := getEnvOrDefault("CONSUL_ADDR", "consul:8500")
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("movie:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}

			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// Get database configuration from environment
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "password")
	dbName := getEnvOrDefault("DB_NAME", "movieexample")

	// Parse port number
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic(fmt.Sprintf("invalid port number: %v", err))
	}

	// Initialize PostgreSQL repository
	repo, err := postgres.New(
		dbHost,
		dbPortInt,
		dbUser,
		dbPassword,
		dbName,
	)
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	// Initialize other dependencies
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	// Initialize controller with both repository and gateways
	ctrl := movie.NewWithRepo(repo, ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
