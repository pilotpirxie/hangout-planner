# hangout-planner

A simple and open source hangout planner.

## How to run

```sh
# backend server
cd backend
go run ./cmd/api

# to run client in a separate terminal
cd client
yarn dev
```

## How to build client

```sh
cd client
yarn build:copy
```

## How to prepare migrations
```sh
# to install goose
# https://github.com/pressly/goose
brew install goose

# to create a new migration
goose create MIGRATION_NAME sql

# to apply migrations
goose up
```

## How to generate query handlers
```sh
sqlc generate
```

## How to format files
```sh
go fmt ./...
```