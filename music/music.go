package music

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"  // Import Oto for audio playback
	"github.com/hajimehoshi/go-mp3" // Import MP3 decoder for decoding MP3 files
	"letsquiz/common"
)

var player *oto.Player         // Player instance for audio playback
var context *oto.Context       // Audio context for managing audio settings
var appModel *common.Model     // Application model that controls playback state
var shutdownChan chan struct{} // Channel to signal when the application is shutting down

// SetModel initializes the music package with the provided application model
func SetModel(model *common.Model) {
	appModel = model
}

// SetShutdownChannel initializes the shutdown channel for handling graceful shutdown
func SetShutdownChannel(ch chan struct{}) {
	shutdownChan = ch
}

// PlayBackgroundMusic starts playing the specified MP3 file as background music
func PlayBackgroundMusic(filePath string) {
	// Initialize the Oto audio context with specified options
	options := &oto.NewContextOptions{
		SampleRate:   44100,                   // Sample rate for audio playback
		ChannelCount: 2,                       // Number of audio channels (stereo)
		Format:       oto.FormatSignedInt16LE, // Audio format (mp3)
		BufferSize:   8192,                    // Buffer size for audio playback
	}
	var err error
	context, readyChan, err := oto.NewContext(options)
	if err != nil {
		log.Fatal("Failed to create audio context:", err) // Log error and exit if context creation fails
	}
	<-readyChan // Wait until the context is ready

	// Read the MP3 file into memory
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Failed to read audio file:", err) // Log error and exit if file reading fails
	}

	// Convert the file bytes into a reader object
	fileBytesReader := bytes.NewReader(fileBytes)

	// Create a new MP3 decoder to decode the MP3 file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		log.Fatal("Failed to create MP3 decoder:", err) // Log error and exit if decoder creation fails
	}

	for {
		select {
		case <-shutdownChan: // Handle shutdown signal
			if player != nil {
				player.Close() // Close the player if it is open
			}
			return // Exit the function when shutdown signal is received
		default:
			// Reset the reader to the beginning of the MP3 file
			fileBytesReader.Seek(0, 0)

			// Create a new player for audio playback
			player = context.NewPlayer(decodedMp3)
			player.Play() // Start playing the audio

			// Start a goroutine to monitor the playback state
			go func() {
				for {
					if appModel == nil {
						time.Sleep(100 * time.Millisecond) // Wait if appModel is not initialized
						continue
					}
					// Pause or play the audio based on the application's playback state
					if !appModel.IsPlaying {
						player.Pause()
					} else {
						player.Play()
					}
					time.Sleep(100 * time.Millisecond) // Sleep to prevent busy loop
				}
			}()

			// Wait for the audio playback to finish
			for player.IsPlaying() {
				time.Sleep(time.Millisecond)
			}

			player.Close() // Close the player after playback finishes
		}
	}
}

// ToggleMusicMuteUnmute toggles the music playback state between playing and paused
func ToggleMusicMuteUnmute() {
	if appModel != nil {
		appModel.IsPlaying = !appModel.IsPlaying // Toggle playback state
		if player != nil {
			// Play or pause the audio based on the new playback state
			if appModel.IsPlaying {
				player.Play()
			} else {
				player.Pause()
			}
		}
	}
}
