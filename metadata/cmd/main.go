package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/phongld0308/movie-example/gen"
	"github.com/phongld0308/movie-example/metadata/internal/controller/metadata"
	grpchandler "github.com/phongld0308/movie-example/metadata/internal/handler/grpc"
	"github.com/phongld0308/movie-example/metadata/internal/repository/postgres"
	"github.com/phongld0308/movie-example/pkg/discovery"
	"github.com/phongld0308/movie-example/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)

	// Get configuration from environment
	consulAddr := getEnvOrDefault("CONSUL_ADDR", "consul:8500")
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("metadata:%d", port)); err != nil {
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

	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
	// h := httphandler.New(ctrl)

	// http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	// if err := http.ListenAndServe(":8081", nil); err != nil {
	// 	panic(err)
	// }
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
