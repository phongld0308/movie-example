# Movie Rating System

A microservices-based movie rating system built with Go, gRPC, and PostgreSQL. The system consists of three main services: Movie, Metadata, and Rating services, orchestrated using Docker Compose.

## Architecture Overview

The system is composed of the following services:

- **Movie Service (Port 8083)**: Aggregates data from metadata and rating services
- **Metadata Service (Port 8081)**: Manages movie metadata (title, description, director)
- **Rating Service (Port 8082)**: Handles movie ratings
- **Consul (Port 8500)**: Service discovery and registration
- **PostgreSQL**: Database for storing movie metadata and ratings

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)
- grpcurl (for testing gRPC endpoints)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/movie-rating
cd movie-rating
```

2. Start the services:
```bash
docker-compose up --build
```

3. The following services will be available:
- Movie Service: http://localhost:8083
- Metadata Service: localhost:8081 (gRPC)
- Rating Service: localhost:8082 (gRPC)
- Consul UI: http://localhost:8500

## Testing the Services

### Metadata Service (gRPC)

```bash
# Add movie metadata
grpcurl -plaintext -d '{
  "metadata": {
    "id": "1",
    "title": "The Matrix",
    "description": "A computer programmer discovers a mysterious world",
    "director": "Lana Wachowski"
  }
}' localhost:8081 MetadataService/PutMetadata

# Get movie metadata
grpcurl -plaintext -d '{"movie_id": "1"}' localhost:8081 MetadataService/GetMetadata
```

### Rating Service (gRPC)

```bash
# Add rating
grpcurl -plaintext -d '{
  "user_id": "user1",
  "record_id": "1",
  "record_type": "movie",
  "rating_value": 5
}' localhost:8082 rating.RatingService/PutRating

# Get aggregated rating
grpcurl -plaintext -d '{
  "record_id": "1",
  "record_type": "movie"
}' localhost:8082 rating.RatingService/GetAggregatedRating
```

### Movie Service (HTTP)

```bash
# Get movie details (combines metadata and rating)
curl -X GET "http://localhost:8083/movie?id=1"
```

## Project Structure

```
.
├── api/                    # Protocol Buffers definitions
├── gen/                    # Generated Protocol Buffers code
├── metadata/              # Metadata service
│   ├── cmd/              # Service entry point
│   ├── internal/         # Internal packages
│   └── pkg/              # Public packages
├── movie/                # Movie service
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── rating/               # Rating service
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── pkg/                  # Shared packages
├── schema/               # Database schemas
└── docker-compose.yml    # Service orchestration
```

## Database Schema

The system uses PostgreSQL with the following main tables:

- `movies`: Stores movie metadata
  - id (Primary Key)
  - title
  - description
  - director
  - created_at
  - updated_at

- `ratings`: Stores user ratings
  - record_id
  - record_type
  - user_id
  - value
  - created_at
  - updated_at

## Development

For local development:

1. Install dependencies:
```bash
go mod download
```

2. Generate Protocol Buffers code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/*.proto
```

3. Run services individually:
```bash
# Run metadata service
go run metadata/cmd/main.go

# Run rating service
go run rating/cmd/main.go

# Run movie service
go run movie/cmd/main.go
```

## Environment Variables

Each service can be configured using the following environment variables:

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=movieexample
CONSUL_ADDR=consul:8500

