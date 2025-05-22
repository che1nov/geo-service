# Geo Service

A service for address lookup and geocoding, implemented using clean architecture principles. The service accesses the [DaData API](https://dadata.ru/) to retrieve address information and supports authentication using JWT. The main functions of the service are:

- **User Registration** - registration with password hashing (bcrypt).
- **User Authorization (login)** - issuance of JWT token.
- **Address Search** - endpoint `/api/address/search` takes a `query` parameter and returns a list of addresses.
- **Geocode** - the `/api/address/geocode` endpoint accepts `lat` and `lng` coordinates and returns location information.
- **Swagger documentation** - all endpoints are documented and available at `/swagger/index.html`.

## Project Structure

- Entities - contains the business objects.
- UseCases - implements the business logic using interfaces.
- Interfaces - defines contracts for the infrastructure.
- Infrastructure - implements external dependencies (DaData, router, logger).
- Adapters - binds HTTP requests to use cases.

## Requirements

- **Go** 1.19 or higher (for local builds).
- **Docker** and **Docker Compose** - for building and running an application in a container.

## Quick Start with Docker Compose
You can start a project with Docker Compose.

## Makefile

The project contains a Makefile with basic commands:
- make build: “Build binary file”

- make run: “Run”

- make docker-build: “Build Docker image”.

- make docker-up: “Start containers via docker-compose”.

- make clean: “Clean up generated files”

Translated with DeepL.com (free version)


