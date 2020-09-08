# Thaichana (inspiration) API

This example App is pretent to be a Thaichana API

## Prerequisite

- go version 1.4 to later
- VSCode
- REST Client Extension
- Docker
- MySQL

## Spin-up Database

> make db

## Create new Database

> make newdb

## Migrate Database

> make migrate

## Run

> go run main.go

## Testing

> go test ./...

## Acceptance Testing

Open file tests/checkin.md

Short key: command + shift + p

Select `Rest Client: Send Request`

## Acceptance Testing with restcli

### install rest-cli

> npm install -g rest-cli

### run

> restcli tests/suite.http
