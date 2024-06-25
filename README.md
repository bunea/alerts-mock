# Google alerts mock

Run without makefile:

`go run main.go`

Run with makefile:

`make run`

Run with docker:

`make dev`

Build image:

`make docker-build`

Run image:

`make docker-run`
or
`docker run -p 1323:1323 'alerts-mock'`


After the server is running you use the following endpoints:

- `GET http://localhost:1323/alerts/feed/feed-1`
  - This endpoint returns a feed with alerts that follows the google alerts structure
  - 3 feeds available:
    - `feed-1` that has a preset number of alerts, but can be updated using another endpoint
    - `feed-2` that automatically updates every 30 seconds
    - `feed-3` that automatically updates every 2 minutes
- `POST http://localhost:1323/alerts/feed/feed-1/entries`
  - This endpoint allows adding an alert to a feed
  - Use the following body:
```
{
  "id": "feed-1-entry-1",
  "title": "Entry title",
  "href": "https://example.com",
  "content": "Content",
  "author_name": "John Doe"
}
```
