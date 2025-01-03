# About

Core backend service for Meowhasiswa. Anonymous forum-based social media app targeted for University and School Students.

<img  src="https://i.redd.it/pa049v3v4tfc1.jpeg"  alt="vro"  width="350"/>

## Installation

Ensure you have the following installed
- go >= 1.22.3

Install dependencies by running:
```
go mod download
```

Run development server by using the command below in **cmd** folder :
```
go run main.go
```

## Tech Stack

- Go
- Gin
- Zap
- Viper
- BunDB
- Postgres
- AWS S3
- Resty

## Project Structure

- **/cmd** : This folder contains the entrypoint of this service
- **/infra** : This folder contains dockerfile and infra related files
- **/internal** : This folder contains core logic and code
  - /api: This folder contains the service routes
  - /config: This folder contains the service config like port number, app env, etc.
  - /constants: This folder contains the service constants
  - /middleware: This folder contains the service middleware such as JWT middleware
  - /model: This folder contains the database entity models and DTO 
  - /repository: This folder contains functions that interact with database
  - /request: This folder contains request objects/struct
  - /request: This folder contains response objects/struct
  - /service: This folder contains business logic for route handlers
- **/migrations** : This folder contains migrations for database
- **/pkg** : This folder contains functions to call external services and commonly used functions like logger, utl functions, etc.
- **/postman** : This folder contains postman collection for this service

## Postman Collection
The postman collection json for this service can be found at /postman. You can download it and import it on your local postman application.

