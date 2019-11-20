# pubgolf

## Setup Background Services (for Frontend Dev)

Install the following:
* Docker
* Docker Compose

Get a `.env` file (`cp .env.example .env`).

Start background services:
```
docker-compose up -d db envoy
bin/migrate
docker-compose up --build -d api
```

Reset database (clears all data):
```
bin/reset-db
```

Accessing logs (`SERVICE_NAME` is either `api` or `envoy`):
```
docker-compose logs SERVICE_NAME
```

Shut down background services:
```
docker-compose down
```

## Running the API Locally (for Backend Dev)

Install the following:
* Docker
* Docker Compose

Get a `.env` file (`cp .env.example .env`), but set `API_HOST=docker.for.mac.localhost` to pass API requests through to the local (non-dockerized) version of the API.

Start background services (database and proxy):
```
docker-compose up -d envoy db
```

Start local instance of API:
```
cd api
bin/run
```

Shut down background services:
```
docker-compose down
```


## Deployment

TODO
