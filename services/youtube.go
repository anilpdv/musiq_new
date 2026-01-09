package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"musiq/models"

	"github.com/kkdai/youtube/v2"
)

// YouTubeService handles all YouTube operations
type YouTubeService struct {
	client youtube.Client
}

// NewYouTubeService creates a new YouTube service
func NewYouTubeService() *YouTubeService {
	return &YouTubeService{
		client: youtube.Client{},
	}
}

// GetVideo retrieves video information by ID
func (s *YouTubeService) GetVideo(videoID string) (*youtube.Video, error) {
	return s.client.GetVideo(videoID)
}

// GetVideoInfo returns formatted video information
func (s *YouTubeService) GetVideoInfo(videoID string) (*models.VideoInfo, error) {
	video, err := s.client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	info := &models.VideoInfo{
		ID:          video.ID,
		Title:       video.Title,
		Author:      video.Author,
		Duration:    video.Duration.String(),
		DurationSec: int(video.Duration.Seconds()),
		Views:       strconv.FormatInt(int64(video.Views), 10),
		Description: video.Description,
		Thumbnails:  convertThumbnails(video.Thumbnails),
		Formats:     convertFormats(video.Formats),
	}

	return info, nil
}

// GetAudioStream returns the best audio stream for a video
func (s *YouTubeService) GetAudioStream(videoID string) (io.ReadCloser, int64, error) {
	video, err := s.client.GetVideo(videoID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get video: %w", err)
	}

	// Get audio-only formats sorted by bitrate
	formats := video.Formats.WithAudioChannels()
	audioFormats := make([]youtube.Format, 0)
	for _, f := range formats {
		if strings.Contains(f.MimeType, "audio") {
			audioFormats = append(audioFormats, f)
		}
	}

	if len(audioFormats) == 0 {
		// Fallback to any format with audio
		audioFormats = formats
	}

	// Sort by audio bitrate (highest first)
	sort.Slice(audioFormats, func(i, j int) bool {
		return audioFormats[i].AverageBitrate > audioFormats[j].AverageBitrate
	})

	if len(audioFormats) == 0 {
		return nil, 0, fmt.Errorf("no audio formats available")
	}

	format := &audioFormats[0]
	stream, size, err := s.client.GetStream(video, format)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stream: %w", err)
	}

	return stream, size, nil
}

// GetCombinedStream returns a stream that has both video and audio combined
// This is faster than muxing separate streams but may be lower quality (360p/720p)
func (s *YouTubeService) GetCombinedStream(videoID string) (io.ReadCloser, string, int64, error) {
	video, err := s.client.GetVideo(videoID)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get video: %w", err)
	}

	// Find formats with both video and audio (progressive formats)
	// Prefer MP4 formats for browser compatibility
	var combinedFormats []youtube.Format
	for _, f := range video.Formats {
		// Format has both video and audio
		if strings.Contains(f.MimeType, "video") && f.AudioChannels > 0 {
			combinedFormats = append(combinedFormats, f)
		}
	}

	if len(combinedFormats) == 0 {
		return nil, "", 0, fmt.Errorf("no combined video+audio formats available")
	}

	// Sort by quality (prefer higher resolution, then MP4 over WebM)
	sort.Slice(combinedFormats, func(i, j int) bool {
		// Prefer MP4 (avc1) over WebM
		iMP4 := strings.Contains(combinedFormats[i].MimeType, "mp4")
		jMP4 := strings.Contains(combinedFormats[j].MimeType, "mp4")
		if iMP4 != jMP4 {
			return iMP4
		}
		// Then prefer higher resolution
		return combinedFormats[i].Height > combinedFormats[j].Height
	})

	format := &combinedFormats[0]
	stream, size, err := s.client.GetStream(video, format)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get stream: %w", err)
	}

	return stream, format.MimeType, size, nil
}

// GetVideoAndAudioStreams returns separate video and audio streams for muxing
func (s *YouTubeService) GetVideoAndAudioStreams(videoID string) (video io.ReadCloser, audio io.ReadCloser, videoInfo *youtube.Video, err error) {
	videoInfo, err = s.client.GetVideo(videoID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get video: %w", err)
	}

	// Get video-only formats, preferring H.264 (avc1) for browser compatibility
	h264Formats := make([]youtube.Format, 0)
	vp9Formats := make([]youtube.Format, 0)

	for _, f := range videoInfo.Formats {
		if strings.Contains(f.MimeType, "video") && f.AudioChannels == 0 {
			// H.264/AVC formats have "avc1" in MIME type
			if strings.Contains(f.MimeType, "avc1") {
				h264Formats = append(h264Formats, f)
			} else {
				vp9Formats = append(vp9Formats, f)
			}
		}
	}

	// Sort H.264 formats by bitrate (highest first)
	sort.Slice(h264Formats, func(i, j int) bool {
		return h264Formats[i].Bitrate > h264Formats[j].Bitrate
	})

	// Sort VP9 formats as fallback
	sort.Slice(vp9Formats, func(i, j int) bool {
		return vp9Formats[i].Bitrate > vp9Formats[j].Bitrate
	})

	// Prefer H.264 for browser compatibility, fall back to VP9
	var videoFormats []youtube.Format
	if len(h264Formats) > 0 {
		videoFormats = h264Formats
	} else {
		videoFormats = vp9Formats
	}

	if len(videoFormats) == 0 {
		return nil, nil, nil, fmt.Errorf("no video-only formats available")
	}

	// Get best audio format (prefer English)
	audioFormats := make([]youtube.Format, 0)
	for _, f := range videoInfo.Formats {
		if strings.Contains(f.MimeType, "audio") {
			audioFormats = append(audioFormats, f)
		}
	}

	// Sort by bitrate
	sort.Slice(audioFormats, func(i, j int) bool {
		return audioFormats[i].AverageBitrate > audioFormats[j].AverageBitrate
	})

	// Try to find English audio
	var selectedAudio *youtube.Format
	for i := range audioFormats {
		if audioFormats[i].AudioTrack != nil {
			if strings.HasPrefix(audioFormats[i].AudioTrack.ID, "en") {
				selectedAudio = &audioFormats[i]
				break
			}
		}
	}
	if selectedAudio == nil && len(audioFormats) > 0 {
		selectedAudio = &audioFormats[0]
	}

	if selectedAudio == nil {
		return nil, nil, nil, fmt.Errorf("no audio formats available")
	}

	// Get video stream
	videoStream, _, err := s.client.GetStream(videoInfo, &videoFormats[0])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get video stream: %w", err)
	}

	// Get audio stream
	audioStream, _, err := s.client.GetStream(videoInfo, selectedAudio)
	if err != nil {
		videoStream.Close()
		return nil, nil, nil, fmt.Errorf("failed to get audio stream: %w", err)
	}

	return videoStream, audioStream, videoInfo, nil
}

// SearchVideos searches YouTube for videos
func (s *YouTubeService) SearchVideos(query string) ([]models.VideoResult, error) {
	return searchYouTube(query, "video")
}

// SearchPlaylists searches YouTube for playlists
// Note: Due to InnerTube API limitations, this returns video results for "query playlist"
func (s *YouTubeService) SearchPlaylists(query string) ([]models.PlaylistResult, error) {
	// Search for query + playlist to get playlist-related results
	videos, err := s.SearchVideos(query + " playlist")
	if err != nil {
		return nil, err
	}

	// Convert video results to playlist results
	playlists := make([]models.PlaylistResult, 0, len(videos))
	for _, v := range videos {
		playlist := models.PlaylistResult{
			ID:         v.ID,
			Title:      v.Title,
			Author:     v.Author,
			VideoCount: 0, // Unknown for video results
			Thumbnails: v.Thumbnails,
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

// GetPlaylistVideos retrieves videos from a playlist
func (s *YouTubeService) GetPlaylistVideos(playlistID string) ([]models.VideoResult, error) {
	playlist, err := s.client.GetPlaylist(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	videos := make([]models.VideoResult, 0, len(playlist.Videos))
	for _, entry := range playlist.Videos {
		video := models.VideoResult{
			ID:          entry.ID,
			Title:       entry.Title,
			Author:      entry.Author,
			Duration:    entry.Duration.String(),
			DurationSec: int(entry.Duration.Seconds()),
			Thumbnails:  convertThumbnails(entry.Thumbnails),
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// GetRelatedVideos retrieves related videos for a video ID
func (s *YouTubeService) GetRelatedVideos(videoID string) ([]models.VideoResult, error) {
	// Use YouTube search as fallback since kkdai/youtube doesn't have direct related videos
	return s.SearchVideos(videoID)
}

// Helper functions

func convertThumbnails(thumbnails youtube.Thumbnails) []models.Thumbnail {
	result := make([]models.Thumbnail, 0, len(thumbnails))
	for _, t := range thumbnails {
		result = append(result, models.Thumbnail{
			URL:    t.URL,
			Width:  int(t.Width),
			Height: int(t.Height),
		})
	}
	return result
}

func convertFormats(formats youtube.FormatList) []models.VideoFormat {
	result := make([]models.VideoFormat, 0, len(formats))
	for _, f := range formats {
		format := models.VideoFormat{
			ItagNo:       f.ItagNo,
			MimeType:     f.MimeType,
			Quality:      f.Quality,
			QualityLabel: f.QualityLabel,
			Bitrate:      f.Bitrate,
			AudioBitrate: f.AverageBitrate,
			Width:        f.Width,
			Height:       f.Height,
			FPS:          f.FPS,
			AudioOnly:    strings.Contains(f.MimeType, "audio"),
			VideoOnly:    f.AudioChannels == 0 && strings.Contains(f.MimeType, "video"),
		}
		result = append(result, format)
	}
	return result
}

// searchYouTube performs a YouTube search using the InnerTube API
func searchYouTube(query string, searchType string) ([]models.VideoResult, error) {
	results, err := searchYouTubeRaw(query, searchType)
	if err != nil {
		return nil, err
	}

	videos := make([]models.VideoResult, 0)
	for _, item := range results {
		if videoID, ok := item["videoId"].(string); ok {
			video := models.VideoResult{
				ID:    videoID,
				Title: getString(item, "title"),
			}

			// Get author
			if ownerText, ok := item["ownerText"].(map[string]interface{}); ok {
				if runs, ok := ownerText["runs"].([]interface{}); ok && len(runs) > 0 {
					if run, ok := runs[0].(map[string]interface{}); ok {
						video.Author = getString(run, "text")
					}
				}
			}

			// Get duration
			if lengthText, ok := item["lengthText"].(map[string]interface{}); ok {
				video.Duration = getString(lengthText, "simpleText")
			}

			// Get views
			if viewCount, ok := item["viewCountText"].(map[string]interface{}); ok {
				video.Views = getString(viewCount, "simpleText")
			}

			// Get thumbnails
			if thumbnail, ok := item["thumbnail"].(map[string]interface{}); ok {
				video.Thumbnails = extractThumbnails(thumbnail)
			}

			videos = append(videos, video)
		}
	}

	return videos, nil
}

// searchYouTubeRaw performs raw YouTube search using InnerTube API
func searchYouTubeRaw(query string, searchType string) ([]map[string]interface{}, error) {
	ctx := context.Background()

	// Build InnerTube API request
	apiURL := "https://www.youtube.com/youtubei/v1/search?prettyPrint=false"

	payload := map[string]interface{}{
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"clientName":    "WEB",
				"clientVersion": "2.20231219.04.00",
				"hl":            "en",
				"gl":            "US",
			},
		},
		"query": query,
	}

	// Note: Params filtering removed - InnerTube API filter params are unreliable
	// Results are filtered by renderer type in extractSearchResults instead

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract video/playlist results from response
	return extractSearchResults(result), nil
}

func extractSearchResults(data map[string]interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)

	// Navigate through the nested response structure
	contents, ok := data["contents"].(map[string]interface{})
	if !ok {
		return results
	}

	twoColumnResults, ok := contents["twoColumnSearchResultsRenderer"].(map[string]interface{})
	if !ok {
		return results
	}

	primaryContents, ok := twoColumnResults["primaryContents"].(map[string]interface{})
	if !ok {
		return results
	}

	sectionListRenderer, ok := primaryContents["sectionListRenderer"].(map[string]interface{})
	if !ok {
		return results
	}

	sectionContents, ok := sectionListRenderer["contents"].([]interface{})
	if !ok {
		return results
	}

	for _, section := range sectionContents {
		sectionMap, ok := section.(map[string]interface{})
		if !ok {
			continue
		}

		itemSectionRenderer, ok := sectionMap["itemSectionRenderer"].(map[string]interface{})
		if !ok {
			continue
		}

		items, ok := itemSectionRenderer["contents"].([]interface{})
		if !ok {
			continue
		}

		for _, item := range items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			// Check for video renderer
			if videoRenderer, ok := itemMap["videoRenderer"].(map[string]interface{}); ok {
				results = append(results, videoRenderer)
			}

			// Check for playlist renderer (multiple possible keys)
			if playlistRenderer, ok := itemMap["playlistRenderer"].(map[string]interface{}); ok {
				results = append(results, playlistRenderer)
			}

			// Check for compact playlist renderer
			if compactPlaylist, ok := itemMap["compactPlaylistRenderer"].(map[string]interface{}); ok {
				results = append(results, compactPlaylist)
			}

			// Check for grid playlist renderer
			if gridPlaylist, ok := itemMap["gridPlaylistRenderer"].(map[string]interface{}); ok {
				results = append(results, gridPlaylist)
			}
		}
	}

	return results
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case map[string]interface{}:
			// Handle nested text objects like {"simpleText": "..."} or {"runs": [...]}
			if simpleText, ok := v["simpleText"].(string); ok {
				return simpleText
			}
			if runs, ok := v["runs"].([]interface{}); ok && len(runs) > 0 {
				if run, ok := runs[0].(map[string]interface{}); ok {
					if text, ok := run["text"].(string); ok {
						return text
					}
				}
			}
		}
	}
	return ""
}

func extractThumbnails(thumbnail map[string]interface{}) []models.Thumbnail {
	thumbnails := make([]models.Thumbnail, 0)
	if thumbList, ok := thumbnail["thumbnails"].([]interface{}); ok {
		for _, t := range thumbList {
			if thumbMap, ok := t.(map[string]interface{}); ok {
				thumb := models.Thumbnail{}
				if url, ok := thumbMap["url"].(string); ok {
					thumb.URL = url
				}
				if width, ok := thumbMap["width"].(float64); ok {
					thumb.Width = int(width)
				}
				if height, ok := thumbMap["height"].(float64); ok {
					thumb.Height = int(height)
				}
				thumbnails = append(thumbnails, thumb)
			}
		}
	}
	return thumbnails
}

// ExtractVideoID extracts video ID from various YouTube URL formats
func ExtractVideoID(input string) string {
	// If it's already a video ID (11 characters, alphanumeric with - and _)
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{11}$`, input); matched {
		return input
	}

	// Try to parse as URL
	u, err := url.Parse(input)
	if err != nil {
		return input
	}

	// Handle youtu.be URLs
	if u.Host == "youtu.be" {
		return strings.TrimPrefix(u.Path, "/")
	}

	// Handle youtube.com URLs
	if strings.Contains(u.Host, "youtube.com") {
		// Check for /watch?v=ID format
		if v := u.Query().Get("v"); v != "" {
			return v
		}
		// Check for /embed/ID format
		if strings.HasPrefix(u.Path, "/embed/") {
			return strings.TrimPrefix(u.Path, "/embed/")
		}
		// Check for /v/ID format
		if strings.HasPrefix(u.Path, "/v/") {
			return strings.TrimPrefix(u.Path, "/v/")
		}
	}

	return input
}
