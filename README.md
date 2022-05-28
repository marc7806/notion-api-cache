# Notion API cache

A proxy for the [official Notion API](https://developers.notion.com/)

## Main features
💫 Blazing fast responses due to database caching

🚀 Notion query syntax support

✅ Flattened response data structure for easier attribute access

## Usage
### Docker
```docker
docker build -t notion-api-cache .

docker run -d \
    -e MONGODB_URI=<mongodb-url> \
    -e MONGODB_NAME=notion-api-cache \
    -e NOTION_API_KEY=<notion-api-secret-key> \
    -e NOTION_DATABASES=<comma-separated-list-of-notion-database-ids> \
    -p 8090:8080 \
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

## Supported Notion Query Operators
* Equals

## Upcoming features
✅ Cursor-based Pagination support

⬜ Scheduler for syncing the notion database data in defined intervals

⬜ Sorting support

⬜ Api-token middleware for authentication

⬜ Better notion query syntax and type support

## Authors
This notion-api-cache was written by [Marc7806](https://github.com/marc7806/)