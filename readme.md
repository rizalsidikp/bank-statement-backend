# Bank Statement Backend

This repository contains the backend service for processing and managing bank statements.

## Features

- Upload and parse bank statements
- Store transaction data securely
- Returns total balance
- Returns non-successful transactions

## Getting Started

### Prerequisites

- Go (version 1.20 or higher)

### Installation

```bash
git clone https://github.com/rizalsidikp/bank-statement-backend.git
cd bank-statement-backend
go mod tidy
go run cmd/main.go
```

## Unit Test

### Running Unit Test

```go test ./... -tags=test -coverprofile=coverage.out```
or
```go test ./... -tags=test "-coverprofile=coverage.out"```

### Show Result Of Unit Test

```go tool cover -html=coverage.out```
or
```go tool cover "-html=coverage.out"```

**Notes:** If the command does not work, take out ".out" from the command

## Docker

### Build Docker

```docker build -t bank-statement-backend .```

### Run Docker Container

```docker run --rm -p 8080:8080 bank-statement-backend```

## Github Workflow
The Github workflow run unit test and build docker
