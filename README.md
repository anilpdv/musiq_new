# Musiq API

A Go REST API for streaming YouTube videos as MP3/MP4.

## Features

- **MP3 Audio Streaming** - Stream any YouTube video as MP3 audio
- **MP4 Video Streaming** - Stream videos with audio in MP4 format
- **Video Search** - Search YouTube videos
- **Related Videos** - Get related videos for discovery
- **Playlist Support** - Browse and stream playlists

## Requirements

- Go 1.20+
- FFmpeg installed and in PATH

## Quick Start

```bash
# Build and run
go build -o musiq-server . && ./musiq-server

# Or run directly
go run main.go
```

Server starts on `http://localhost:8080`

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Health check, list all routes |
| `GET /api/search/:q` | Search YouTube videos |
| `GET /api/listen/:id/:name` | Stream MP3 audio |
| `GET /api/watch/:id/:name` | Stream MP4 video |
| `GET /api/info/:id` | Get video metadata |
| `GET /api/getvideo/:id` | Get related videos |
| `GET /api/related/:id` | Get video details + related |
| `GET /api/playlist/search/:q` | Search playlists |
| `GET /api/getplaylist/:id` | Get playlist videos |

## Usage Examples

```bash
# Search for videos
curl "http://localhost:8080/api/search/lofi"

# Get video info
curl "http://localhost:8080/api/info/dQw4w9WgXcQ"

# Download MP3
curl "http://localhost:8080/api/listen/dQw4w9WgXcQ/song.mp3" --output song.mp3

# Stream video
curl "http://localhost:8080/api/watch/dQw4w9WgXcQ/video.mp4" --output video.mp4

# Get playlist videos
curl "http://localhost:8080/api/getplaylist/PLrAXtmErZgOeiKm4sgNOknGvNjby9efdf"
```

## Project Structure

```
musiq/
├── main.go              # Server entry point
├── handlers/            # HTTP route handlers
│   ├── search.go
│   ├── listen.go        # MP3 streaming
│   ├── watch.go         # MP4 streaming
│   ├── info.go
│   ├── related.go
│   └── playlist.go
├── services/            # Business logic
│   ├── youtube.go       # YouTube client
│   └── ffmpeg.go        # FFmpeg operations
├── middleware/          # HTTP middleware
│   └── cors.go
└── models/              # Data structures
    └── types.go
```

## Tech Stack

- **Gin** - Fast HTTP web framework
- **kkdai/youtube** - YouTube video downloading
- **ffmpeg-go** - FFmpeg wrapper for transcoding

## License

MIT
