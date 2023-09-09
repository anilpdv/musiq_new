# Musiq New Documentation

This documentation provides details about the endpoints and usage of the My API.

## Base URL

The base URL for all API endpoints is `/api`. You can access the API using both HTTP and HTTPS.

## Endpoints

### Search Videos

- **URL**: `/search/{q}`
- **Method**: GET
- **Description**: Search for videos.
- **Parameters**:
  - `q` (path) - Search query (string, required)
- **Response**:
  - Status 200: OK

### Get Video

- **URL**: `/listen/{id}/{name}`
- **Method**: GET
- **Description**: Get information about a video.
- **Parameters**:
  - `id` (path) - Video ID (string, required)
  - `name` (path) - Video name (string, required)
- **Response**:
  - Status 200: OK

### Get Related Videos

- **URL**: `/getvideo/{id}`
- **Method**: GET
- **Description**: Get related videos for a given video ID.
- **Parameters**:
  - `id` (path) - Video ID (string, required)
- **Response**:
  - Status 200: OK

### Search Playlists

- **URL**: `/playlist/search/{q}`
- **Method**: GET
- **Description**: Search for playlists.
- **Parameters**:
  - `q` (path) - Search query (string, required)
- **Response**:
  - Status 200: OK

### Get Playlist

- **URL**: `/getplaylist/{id}`
- **Method**: GET
- **Description**: Get information about a playlist.
- **Parameters**:
  - `id` (path) - Playlist ID (string, required)
- **Response**:
  - Status 200: OK

## Version

This documentation is based on API version 1.0.0.

## How to Use

You can make requests to these endpoints using your preferred HTTP client.

Example using cURL:

```shell
# Search for videos
curl -X GET "/api/search/{query}"

# Get video by ID and name
curl -X GET "/api/listen/{id}/{name}"

# Get related videos by ID
curl -X GET "/api/getvideo/{id}"

# Search for playlists
curl -X GET "/api/playlist/search/{query}"

# Get playlist by ID
curl -X GET "/api/getplaylist/{id}"
