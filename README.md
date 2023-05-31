# url-shortener

A service that shortens URLs and is written in Go.

## Design

Please check the [Design](./docs/Design.md) for further design details.

## Requirement

- [Go](https://go.dev/) > 1.20.4
- [redis](https://redis.io/) > 7.0.0
- [MongoDB](https://www.mongodb.com/) > 6.0.0

## Quick Start

### Requirement

- [docker](https://docs.docker.com/)
- [docker compose](https://docs.docker.com/)

### Clone

```sh
git clone git@github.com:WeiAnAn/url-shortener.git
cd url-shortener
```

### Build Image

```sh
docker compose build
```

### Start Containers

```sh
docker compose up
# or
docker compose up -d # run in the background
```

### Enjoy!

## Run In The Local

### Clone

```sh
git clone git@github.com:WeiAnAn/url-shortener.git
cd url-shortener
```

### Install Dependencies

```sh
go mod download
```

### Run

See [Configuration](#configuration) to set up your config. e.g. Database host, username, password.

```sh
go run cmd/server/main.go
```

## Testing

```sh
go test ./...
```

## Build

Recommend to use docker to build image for deployment.

### Build with Docker

```sh
docker build -t <tag> .
```

### Build in the Local

```sh
# In macOS
CGO_ENABLED=0 GOOS=darwin go build -o ./server ./cmd/server/main.go
# In linux
CGO_ENABLED=0 GOOS=linux go build -o ./server ./cmd/server/main.go
```

Start server

```sh
./server
```

## API

### POST /api/v1/urls

Create the new short url.

**Request Body**

content-type: `application/json`

| field    | type   | constraints                                                                                                                                                  |
| -------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| url      | string | Required. Must be http or https URL scheme and less than 2048 characters                                                                                     |
| expireAt | string | Required. Must be in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339) format. The Date must greater than the current time and less than a year later |

**Response Body**

content-type: `application/json`

| field    | type   | description         |
| -------- | ------ | ------------------- |
| id       | string | short url id        |
| shortUrl | string | generated short url |

**Sample Request and Response**

```sh
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{
  "url": "https://pkg.go.dev",
  "expireAt": "2023-05-31T00:00:00Z"
}'

# Response
{
  "id": "abcdefg",
  "shortUrl": "http://localhost/abcdefg"
}
```

### GET /:url_id

Redirect to the original URL by giving url_id.

If the link not found or expired, the server will response 404.

**Sample Request and Response**

```sh
curl -X GET http://localhost/abcdefg
# Redirect to https://pkg.go.dev
```

## Configuration

You can configure this app by setting below environment variables

### Environment Variables

| Variable    | Description                                                                                                                | Default VALUE                       |
| ----------- | -------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| MONGODB_URI | MongoDB connection string. See [the link](https://www.mongodb.com/docs/manual/reference/connection-string/) for more info. | mongodb://short_url@localhost:27017 |
| REDIS_HOST  | redis connection. format: \<host\>:\<port\>.                                                                               | localhost:6379                      |
| BASE_URL    | short url base url. Generated short url id will append to this base url.                                                   | http://localhost:8080               |
| GIN_MODE    | Gin running mode. Please make sure to set this value to 'release' when you are running in the production environment.      | debug                               |

## Postgres Version

If you want to use postgres as your database, please check [postgres branch](https://github.com/WeiAnAn/url-shortener/tree/postgres)
