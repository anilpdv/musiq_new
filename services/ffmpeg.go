package services

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// FFmpegService handles audio/video conversion
type FFmpegService struct{}

// NewFFmpegService creates a new FFmpeg service
func NewFFmpegService() *FFmpegService {
	return &FFmpegService{}
}

// ConvertToMP3 converts an audio stream to MP3 format
func (s *FFmpegService) ConvertToMP3(input io.Reader, output io.Writer) error {
	err := ffmpeg.Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"acodec": "libmp3lame",
			"q:a":    "0",
			"f":      "mp3",
		}).
		WithInput(input).
		WithOutput(output, os.Stderr).
		Run()

	if err != nil {
		return fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return nil
}

// MuxVideoAudio muxes separate video and audio streams into MP4
// This uses os/exec directly for better control over multiple input pipes
func (s *FFmpegService) MuxVideoAudio(videoStream, audioStream io.Reader, output io.Writer) error {
	// Find ffmpeg path
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	// Prepare ffmpeg command with multiple inputs
	cmd := exec.Command(ffmpegPath,
		"-i", "pipe:3", // Video input
		"-i", "pipe:4", // Audio input
		"-map", "0:v", // Map video from first input
		"-map", "1:a", // Map audio from second input
		"-c:v", "copy", // Copy video codec (no re-encoding)
		"-c:a", "aac", // Convert audio to AAC for mobile compatibility
		"-movflags", "frag_keyframe+empty_moov", // Fragmented MP4 for streaming
		"-f", "mp4", // Output format
		"-loglevel", "error",
		"-", // Output to stdout
	)

	// Set up pipes
	// stdin (0), stdout (1), stderr (2), video (3), audio (4)
	videoPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create video pipe: %w", err)
	}

	// We need to use ExtraFiles for additional file descriptors
	videoReader, videoWriter, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create video pipe: %w", err)
	}
	defer videoReader.Close()

	audioReader, audioWriter, err := os.Pipe()
	if err != nil {
		videoWriter.Close()
		return fmt.Errorf("failed to create audio pipe: %w", err)
	}
	defer audioReader.Close()

	// Close stdin pipe as we're using ExtraFiles
	videoPipe.Close()

	cmd.ExtraFiles = []*os.File{videoReader, audioReader}
	cmd.Stdout = output.(io.Writer)
	cmd.Stderr = os.Stderr

	// Update command to use fd 3 and 4
	cmd.Args = []string{
		ffmpegPath,
		"-i", "pipe:3",
		"-i", "pipe:4",
		"-map", "0:v",
		"-map", "1:a",
		"-c:v", "copy",
		"-c:a", "aac",
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"-f", "mp4",
		"-loglevel", "warning",
		"-",
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		videoWriter.Close()
		audioWriter.Close()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Close readers in parent process (they're being used by child)
	videoReader.Close()
	audioReader.Close()

	// Copy video stream in goroutine
	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(videoWriter, videoStream)
		videoWriter.Close()
		errChan <- err
	}()

	// Copy audio stream in goroutine
	go func() {
		_, err := io.Copy(audioWriter, audioStream)
		audioWriter.Close()
		errChan <- err
	}()

	// Wait for both copies to complete
	var copyErr error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil && copyErr == nil {
			copyErr = err
		}
	}

	// Wait for ffmpeg to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return copyErr
}

// MuxVideoAudioStream uses named pipes (FIFOs) for progressive streaming
// This allows the browser to start playing while data is still being downloaded
func (s *FFmpegService) MuxVideoAudioStream(videoStream, audioStream io.Reader, output io.Writer) error {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}

	// Create unique FIFO paths
	id := uuid.New().String()[:8]
	videoFifo := fmt.Sprintf("/tmp/video_%s.fifo", id)
	audioFifo := fmt.Sprintf("/tmp/audio_%s.fifo", id)

	// Create FIFOs
	if err := syscall.Mkfifo(videoFifo, 0600); err != nil {
		return fmt.Errorf("failed to create video fifo: %w", err)
	}
	defer os.Remove(videoFifo)

	if err := syscall.Mkfifo(audioFifo, 0600); err != nil {
		os.Remove(videoFifo)
		return fmt.Errorf("failed to create audio fifo: %w", err)
	}
	defer os.Remove(audioFifo)

	// Prepare ffmpeg command with flags for progressive streaming
	cmd := exec.Command(ffmpegPath,
		"-i", videoFifo,
		"-i", audioFifo,
		"-c:v", "copy", // Copy H.264 video (no re-encoding needed)
		"-c:a", "aac",
		"-movflags", "frag_keyframe+empty_moov+default_base_moof", // Critical for streaming
		"-frag_duration", "1000000", // 1 second fragments for faster start
		"-f", "mp4",
		"-loglevel", "warning",
		"-",
	)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr

	// Start ffmpeg (it will block waiting for FIFO input)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Channel to collect errors from goroutines
	errChan := make(chan error, 2)

	// Write video stream to FIFO in goroutine
	go func() {
		f, err := os.OpenFile(videoFifo, os.O_WRONLY, 0)
		if err != nil {
			errChan <- fmt.Errorf("failed to open video fifo: %w", err)
			return
		}
		_, err = io.Copy(f, videoStream)
		f.Close()
		errChan <- err
	}()

	// Write audio stream to FIFO in goroutine
	go func() {
		f, err := os.OpenFile(audioFifo, os.O_WRONLY, 0)
		if err != nil {
			errChan <- fmt.Errorf("failed to open audio fifo: %w", err)
			return
		}
		_, err = io.Copy(f, audioStream)
		f.Close()
		errChan <- err
	}()

	// Wait for both writes to complete
	var writeErr error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil && writeErr == nil {
			writeErr = err
		}
	}

	// Wait for ffmpeg to finish
	if err := cmd.Wait(); err != nil {
		// Ignore broken pipe errors (client closed connection)
		if writeErr == nil {
			return fmt.Errorf("ffmpeg failed: %w", err)
		}
	}

	return writeErr
}

// MuxVideoAudioSimple is a simpler version that creates temporary files
// Use this if the pipe-based version has issues
func (s *FFmpegService) MuxVideoAudioSimple(videoStream, audioStream io.Reader, output io.Writer) error {
	// Create temporary files for video and audio
	videoTmp, err := os.CreateTemp("", "video-*.mp4")
	if err != nil {
		return fmt.Errorf("failed to create temp video file: %w", err)
	}
	defer os.Remove(videoTmp.Name())

	audioTmp, err := os.CreateTemp("", "audio-*.m4a")
	if err != nil {
		return fmt.Errorf("failed to create temp audio file: %w", err)
	}
	defer os.Remove(audioTmp.Name())

	outputTmp, err := os.CreateTemp("", "output-*.mp4")
	if err != nil {
		return fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer os.Remove(outputTmp.Name())

	// Write streams to temp files
	if _, err := io.Copy(videoTmp, videoStream); err != nil {
		videoTmp.Close()
		return fmt.Errorf("failed to write video: %w", err)
	}
	videoTmp.Close()

	if _, err := io.Copy(audioTmp, audioStream); err != nil {
		audioTmp.Close()
		return fmt.Errorf("failed to write audio: %w", err)
	}
	audioTmp.Close()

	// Find ffmpeg path
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found: %w", err)
	}

	// Run ffmpeg with proper command structure
	cmd := exec.Command(ffmpegPath,
		"-i", videoTmp.Name(),
		"-i", audioTmp.Name(),
		"-map", "0:v",
		"-map", "1:a",
		"-c:v", "copy",
		"-c:a", "aac",
		"-movflags", "frag_keyframe+empty_moov",
		"-y",
		outputTmp.Name(),
	)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg muxing failed: %w", err)
	}

	// Read output and write to response
	outputFile, err := os.Open(outputTmp.Name())
	if err != nil {
		return fmt.Errorf("failed to open output: %w", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(output, outputFile)
	return err
}

// CheckFFmpegInstalled verifies FFmpeg is available
func CheckFFmpegInstalled() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg not found. Please install ffmpeg and ensure it's in your PATH")
	}
	return nil
}
