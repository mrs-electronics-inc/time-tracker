# Quickstart: 002-docker-container

## Building the Docker Image

```bash
docker build -t time-tracker .
```

## Running the Container

```bash
docker run --rm time-tracker [command] [args]
```

## Testing in Container

The LLM agent can use the container for safe testing without affecting the host system.