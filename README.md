# customer-playground

.
├── .config.toml
├── Dockerfile
├── docker-compose.yaml
├── add-swagger.txt
├── go.mod
├── go.sum
├── main.go
│
├── app/
│   └── app.go                 # Application setup and server start
│
├── database/
│   └── database.go            # DB connection pool setup
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml           # Auto-generated Swagger docs
│
├── domain/
│   ├── customer.go            # Customer entity
│   ├── customernote.go        # Customer note entity
│   └── response.go            # Shared response types
│
├── init/
│   └── init.sql               # SQL script for DB initialization
│
├── services/
│   ├── customer/
│   │   ├── delivery/
│   │   │   └── handler.customer.go
│   │   ├── repository/
│   │   │   └── repository.customer.go
│   │   └── usecase/
│   │       └── usecase.customer.go
│   │
│   └── customernote/
│       ├── delivery/
│       │   └── handler.customernote.go
│       ├── repository/
│       │   └── repository.customernote.go
│       └── usecase/
│           └── usecase.customernote.go
│
├── types/
│   └── time.go                # Custom nullable time type (NullTime)


how to run:
- open terminal in project base dir
- then run "docker-compose up -d"
- swagger url at "http://localhost:8080/swagger/index.html"
