# Requirments
- taskfile: https://taskfile.dev/installation/
- air: https://github.com/air-verse/air
- migrate: https://github.com/golang-migrate/migrate

# Setup

```
cp .env.example .env
task setup
task dev
```

# Endpoints

## Create

```
POST - localhost:8085/events

body
{
	"title": "test33",
	"description": "test12334",
	"start_time": "2025-12-01T10:00:00Z",
	"end_time": "2025-12-03T12:00:00Z"
}
```

```
curl --request POST \
  --url http://localhost:8085/events \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/12.1.0' \
  --data '{
	"title": "test33",
	"description": "test12334",
	"start_time": "2025-12-01T10:00:00Z",
	"end_time": "2025-12-03T12:00:00Z"
}'
```

## Get by ID
```
GET - localhost:8085/events/[id]
```

```
curl --request GET \
  --url http://localhost:8085/events/[id] \
  --header 'User-Agent: insomnia/12.1.0'
```

## Get All
```
GET - localhost:8085/events
```

```
curl --request GET \
  --url http://localhost:8085/events \
  --header 'User-Agent: insomnia/12.1.0'
```