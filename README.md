# Notion API cache

A proxy for the [official Notion API](https://developers.notion.com/)

## Main features
💫 Blazing fast responses due to database caching

🚀 Notion query syntax support

✅ Flattened response data structure for easier attribute access

## Usage
### Query Endpoints
In order to query the notion-api-cache:
```
POST http://localhost:8080/v1/databases/<database-id>/query
{
    "page_size": 5,
    "start_cursor": "ef678e4c-54e6-4d71-ad31-93403be247e1",
    "filter": {
        "and": [
            {
                "property": "My Custom Property",
                "text": {
                    "equals": "Hello World"
                }
            }
        ]
    }
}
```
#### Currently supported Notion filter operators
* equals

### Cache Management Endpoints
Update Cache:
```
POST http://localhost:8080/v1/cache/refresh
```
Clear Cache:
```
POST http://localhost:8080/v1/cache/status
```
Get Cache Status:
```
GET http://localhost:8080/v1/cache/clear
```
The Caching Scheduler can be configure using the ``CACHE_SCHEDULER_HOURS``, ``CACHE_SCHEDULER_MINUTES`` and ``CACHE_SCHEDULER_DAYS`` environment variables. ``CACHE_SCHEDULER_DAYS=1`` would mean that the notion databases get synced once every day.

### Docker
```docker
docker build -t notion-api-cache .

docker run -d \
    -e MONGODB_URI=<mongodb-url> \
    -e MONGODB_NAME=notion-api-cache \
    -e NOTION_API_KEY=<notion-api-secret-key> \
    -e NOTION_DATABASES=<comma-separated-list-of-notion-database-ids> \
    -e CACHE_SCHEDULER_DAYS=1 \
    -p 8080:8080 \
    --name notion-api-cache \
    notion-api-cache:latest
```

### Docker Compose
Create a .env file in your root directory. An example file is already present ``.env.sample``.
The docker compose also starts a mongodb container.

```docker
cd docker/
docker compose up -d
```

## Supported caching providers
### Mongodb
Use mongodb to store your notion database data and query it with the native notion query syntax.

## Upcoming features
✅ Cursor-based Pagination support

✅ Scheduler for syncing the notion database data in defined intervals

⬜ Sorting support

⬜ Api-token middleware for authentication

⬜ Better notion query syntax and type support

## Authors
This notion-api-cache was written by [Marc7806](https://github.com/marc7806/)