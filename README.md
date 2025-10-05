# GoTodo — Project Overview & Architecture (for AI)

This README explains the project structure, core architecture, and runtime details so an AI (or new developer) can understand and reason about the codebase. Includes short examples and commands for running locally.

## Project summary
GoTodo is a RESTful Go API (Gin) that manages todos with JWT authentication and MongoDB persistence. Main features: user registration/login, JWT-protected endpoints, CRUD for todos, and Swagger docs.

## High-level architecture
- HTTP Server: Gin framework. All routes registered in a central routes package.
- Auth: JWT-based authentication middleware for protected routes.
- Persistence: MongoDB accessed via a database package that returns a mongo.Database instance.
- Config: Centralized config loader that reads environment variables (used by app and DB).
- Middleware: Logging, CORS, Auth lived as reusable Gin middlewares.
- Docs: swaggo-generated Swagger docs served at /swagger/*.

Sequence (request):
1. Client -> NGINX/load balancer (if used) -> app
2. Global middlewares (logger, recovery, CORS)
3. Route-level middlewares (auth for protected routes)
4. Handler -> service/DB -> mongo
5. Response -> client

## Directory structure (important files)
- cmd/api/main.go — application entry; server setup, Swagger, graceful shutdown
- internal/config — config loader (env/defaults)
- internal/database — MongoDB connection helper
- internal/middleware — logger, CORS, auth middleware
- internal/routes — register route groups and feature routes
- internal/features/auth — auth handlers, routes, models
- internal/features/todos — todos handlers, routes, models
- docs — generated Swagger docs
- go.mod / Makefile — build/run/test helpers

## Key components & responsibilities
- config.Load(): returns config struct (Port, MongoURI, MongoDB, AppEnv, FrontendURL, JWT secret etc.)
- database.Connect(mongoURI, dbName): creates mongo client and returns wrapper with Database and Disconnect()
- middleware.Auth(): extracts Authorization header, validates JWT, sets user context
- routes.RegisterRoutes(router, db): wires feature routes (auth, todos) and passes mongo.Database
- Feature handlers receive the mongo.Database (or collection) and perform CRUD operations.

## Data models (examples)
User and Todo shapes (approximate):

    type User struct {
        ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
        Email    string             `bson:"email" json:"email"`
        Password string             `bson:"password,omitempty" json:"-"`
        Created  time.Time          `bson:"created_at" json:"created_at"`
    }

    type Todo struct {
        ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
        Title       string             `bson:"title" json:"title"`
        Description string             `bson:"description,omitempty" json:"description"`
        Completed   bool               `bson:"completed" json:"completed"`
        OwnerID     primitive.ObjectID `bson:"owner_id" json:"owner_id"`
        Created     time.Time          `bson:"created_at" json:"created_at"`
        Updated     time.Time          `bson:"updated_at" json:"updated_at"`
    }

## API endpoints (copy/paste)
- GET    /health
- GET    /swagger/*any
- POST   /api/v1/auth/register
- POST   /api/v1/auth/login
- GET    /api/v1/auth/me         (protected)
- POST   /api/v1/todos/         (protected)
- GET    /api/v1/todos/         (protected)
- GET    /api/v1/todos/:id      (protected)
- PUT    /api/v1/todos/:id      (protected)
- DELETE /api/v1/todos/:id      (protected)

Authentication: protected endpoints require header `Authorization: Bearer <token>`.

## Example usage

Run MongoDB (Docker):
    docker run -d --name mongo -p 27017:27017 -v mongodata:/data/db mongo:6.0
    docker logs -f mongo

Set env and run:
    export MONGO_URI="mongodb://127.0.0.1:27017"
    export MONGO_DB="gotodo"
    export APP_PORT="8080"
    export JWT_SECRET="replace-with-secret"
    make run
or
    go run ./cmd/api

Simple curl flows:
- Register
    curl -X POST http://localhost:8080/api/v1/auth/register \
      -H "Content-Type: application/json" \
      -d '{"email":"you@example.com","password":"secret"}'

- Login (receive token)
    curl -X POST http://localhost:8080/api/v1/auth/login \
      -H "Content-Type: application/json" \
      -d '{"email":"you@example.com","password":"secret"}'

- Use token
    export TOKEN="eyJ..."
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/auth/me

- Create todo
    curl -X POST http://localhost:8080/api/v1/todos \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"title":"Buy milk","description":"2 liters"}'

## How to check server & DB health
- Server: curl http://localhost:8080/health
- Mongo: nc -vz 127.0.0.1 27017 or lsof -i :27017
- Docker containers: docker ps -a --filter name=mongo

## Testing & development notes
- Use `go test ./...` to run unit tests (if implemented).
- Use the Swagger UI at http://localhost:8080/swagger/index.html for interactive API exploration.
- Prefer 127.0.0.1 in MONGO_URI to avoid IPv6 localhost issues.

## For the AI: important files and assumptions
- Entry point: cmd/api/main.go — follow server setup, config.Load, database.Connect, routes.RegisterRoutes
- Database objects passed around are mongo.Database from the official driver.
- Middleware chain: logger -> CORS -> (auth for protected routes) -> handler.
- JWT secret and Mongo URI are environment-configurable.
- The project expects a running MongoDB at startup; the server will exit if DB unreachable.

## Troubleshooting
- "Failed to connect to MongoDB: dial tcp 127.0.0.1:27017: connect: connection refused" — start MongoDB or point MONGO_URI to a running DB.
- If Docker complains about a container name conflict:
    docker ps -a --filter name=mongo
    docker start mongo          # restart an existing container
    or
    docker rm -f mongo
    docker run -d --name mongo -p 27017:27017 -v mongodata:/data/db mongo:6.0

## Useful commands
- Build: `go build ./cmd/api`
- Run (dev): `go run ./cmd/api`
- Run with Makefile: `make run`
- Lint: `golangci-lint run` (if configured)

## Contact points in code (where to change behavior)
- config.Load() — change default ports/secrets
- database.Connect() — change mongo options/retry logic
- middleware/* — change logging/CORS/auth behavior
- internal/features/* — business logic for auth and todos

---
This file is intended for quick onboarding and for an AI to parse program structure, runtime assumptions, and where to find core logic. Adjust environment variable names to match internal/config if they differ.

