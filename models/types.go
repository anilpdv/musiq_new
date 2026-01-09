package models

// RootResponse represents the response for the root endpoint
type RootResponse struct {
	Status int    `json:"status"`
	Routes Routes `json:"routes"`
}

// Routes contains all available API routes
type Routes struct {
	SearchRoute       string `json:"searchRoute"`
	ListenRoute       string `json:"listenRoute"`
	WatchRoute        string `json:"watchRoute"`
	InfoRoute         string `json:"infoRoute"`
	RelatedRoute      string `json:"relatedRoute"`
	PlaylistRoute     string `json:"playlistRoute"`
	PlaylistRouteByID string `json:"playlistRouteById"`
}

// VideoResult represents a video in search results
type VideoResult struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Author      string      `json:"author"`
	Duration    string      `json:"duration"`
	DurationSec int         `json:"durationSec"`
	Views       string      `json:"views"`
	Thumbnails  []Thumbnail `json:"thumbnails"`
}

// Thumbnail represents a video thumbnail
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// VideoInfo represents detailed video information
type VideoInfo struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Author       string        `json:"author"`
	Duration     string        `json:"duration"`
	DurationSec  int           `json:"durationSec"`
	Views        string        `json:"views"`
	Description  string        `json:"description"`
	Thumbnails   []Thumbnail   `json:"thumbnails"`
	Formats      []VideoFormat `json:"formats"`
	RelatedSongs []VideoResult `json:"relatedSongs,omitempty"`
}

// VideoFormat represents a video/audio format
type VideoFormat struct {
	ItagNo       int    `json:"itag"`
	MimeType     string `json:"mimeType"`
	Quality      string `json:"quality"`
	QualityLabel string `json:"qualityLabel,omitempty"`
	Bitrate      int    `json:"bitrate"`
	AudioBitrate int    `json:"audioBitrate,omitempty"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	FPS          int    `json:"fps,omitempty"`
	AudioOnly    bool   `json:"audioOnly"`
	VideoOnly    bool   `json:"videoOnly"`
}

// PlaylistResult represents a playlist in search results
type PlaylistResult struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Author     string      `json:"author"`
	VideoCount int         `json:"videoCount"`
	Thumbnails []Thumbnail `json:"thumbnails"`
}

// RelatedResponse represents the response for the /related endpoint
type RelatedResponse struct {
	VideoDetails VideoInfo     `json:"videoDetails"`
	RelatedSongs []VideoResult `json:"relatedSongs"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
