hasUniversity
=============

API middleware for hasCorp's online learning platform

## Running locally

Use `docker compose` to start. The `.local` compose file creates and inserts
dummy data into the Cassandra instance.
```bash
docker compose -f ./docker-compose.local.yml up --build --force-recreate
```

Clean up the docker network:
```bash
docker compose down
```

### Building
Simply done with:
```bash
go build .
```
Or, you can build the docker image by itself:
```bash
docker build -t hascorp/hasuniversity -f Dockerfile .
```

### Running
Recommended to run with `docker compose`, though it is possible to run
locally:
```bash
go run ./main.go
```

### Testing
Unit tests can be run by executing:
```bash
go test -v ./...
```

The API is exposed to port `:8200`, so you can curl it:
```bash
# ping healthcheck endpoint
curl localhost:8200/

# add a flashcard set to the Cassandra database
curl localhost:8200/flashcard/ -d '{"author": "Murat Piker", "name": "Being a Gigachad", "tags": ["Murat"], "cards": [{"front": "engineer", "back": "NODDERS"}, {"front": "twitch streamer", "back": "NOPERS"}]}'
```
