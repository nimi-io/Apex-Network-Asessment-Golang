# Email API Service

A Go-based REST API service for sending emails with a queue-based processing system.

## Features

- RESTful API for email sending
- Queue-based email processing with multiple workers
- Input validation and error handling
- Dockerized for easy deployment

## API Endpoints

### POST /api/v1/send-email

Send an email by adding it to the processing queue.

**Request Body:**

```json
{
  "to": "user@example.com",
  "subject": "Welcome!",
  "body": "Thanks for signing up."
}
```

**Responses:**

- `202 Accepted` - Email successfully queued
- `422 Unprocessable Entity` - Invalid input
- `503 Service Unavailable` - Queue is full

## Running with Docker

### Prerequisites

- Docker
- Docker Compose

### Build and Run

1. Build and start the service:

```bash
docker-compose up --build
```

2. Run in detached mode:

```bash
docker-compose up -d --build
```

3. View logs:

```bash
docker-compose logs -f email-api
```

4. Stop the service:

```bash
docker-compose down
```

### Testing the API

Once the service is running, you can test it with curl:

```bash
# Test valid email
curl -X POST http://localhost:3000/api/v1/send-email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "subject": "Welcome!",
    "body": "Thanks for signing up."
  }'

# Test invalid email (missing fields)
curl -X POST http://localhost:3000/api/v1/send-email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "invalid-email",
    "subject": "Test"
  }'
```

## Development

### Running locally without Docker

```bash
go mod download
go run main.go
```

The service will start on port 3000 by default.

### Environment Variables

- `PORT` - Port to run the service on (default: 3000)
- `GIN_MODE` - Gin framework mode (development/release)
