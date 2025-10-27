# HTTP Calendar Server

To start server, run `go run cmd/main.go`

To test handlers, run `go test ./internal/api/handler`

## Project structure

```
cmd/          - main.go
internal/
  api/        - HTTP handlers & middleware
  service/    - business logic
  repository/ - data storage
  model/
  config/
  app/        - server initialization logic
pkg/
  date/       - utility for date parsing
  http/       - utility for http req & resp
```

## API Endpoints

### POST /create_event
```json
{
  "name": "event name",
  "date": "2024-01-15", // YYYY-MM-DD
  "user_id": 1
}
```

### GET /events_for_day
```
/events_for_day?user_id=1&date=2024-01-15
```

### GET /events_for_week
```
/events_for_week?user_id=1&date=2024-01-15
```

### GET /events_for_month
```
/events_for_month?user_id=1&date=2024-01-15
```

### POST /update_event
```json
{
  "id": 1,
  "name": "updated event", // optional
  "date": "2024-01-16"     // optional
}
```

### POST /delete_event
```json
{
  "id": 1
}
```

## Configuration

In the root, create `config.yaml`:
```yaml
port: 8080
```