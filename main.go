package main

import (
	"fmt"
	"letsquiz/models"
	"letsquiz/music"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea" // Import Bubble Tea framework
	"letsquiz/config"
	"letsquiz/logger"
	"letsquiz/screens"
)

func main() {
	// Load the application configuration from a JSON file
	err := config.LoadConfig("appconfig.json", false)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err) // Print error and exit if configuration loading fails
		os.Exit(1)
	}

	// Initialize the logger
	logger.InitLogger(false)
	logger.Info("Application started") // Log the application startup

	// Create a shutdown channel to signal when the app is closing
	shutdownChan := make(chan struct{})
	music.SetShutdownChannel(shutdownChan) // Set the shutdown channel for the music package

	// Initialize the model for login screen
	model := models.InitialLoginModel()
	music.SetModel(&model.Model) // Set the model for the music package

	// Start playing background music in a separate goroutine
	go music.PlayBackgroundMusic(fmt.Sprintf("%s.mp3", config.AppConfig.MainMp3Track))

	// Set up a channel to catch OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // Notify on interrupt or termination signals

	// Create a new Bubble Tea program with initial model and settings
	p := tea.NewProgram(screens.InitialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())

	// Start a goroutine to listen for shutdown signals and quit the program
	go func() {
		<-signalChan // Wait for a signal
		p.Quit()     // Quit the Bubble Tea program
	}()

	// Run the Bubble Tea program and handle any errors
	m, err := p.Run()
	if err != nil {
		logger.Error("Error running program", "error", err) // Log error if the program fails to run
		fmt.Printf("Alas, there's been an error: %v", err)  // Print error and exit
		os.Exit(1)
	}

	// Close the shutdown channel to signal the music package to stop
	close(shutdownChan)

	// Type switch to handle different screen models and print information
	switch model := m.(type) {
	case *screens.Login:
		logger.Info(fmt.Sprintf("Selected option (LoginModel): %s\n", model))
		fmt.Printf("Selected option (LoginModel): %s\n", model)
	case *screens.Menu:
		logger.Info(fmt.Sprintf("Selected option (Menu): %v\n", model))
		fmt.Printf("Entering Menu with model: %v\n", model)
	case *screens.EditQuestionnaire:
		logger.Info(fmt.Sprintf("Selected option (EditQuestionnaire): %v\n", model))
		fmt.Printf("Editing questionnaire with model: %v\n", model)
	case *screens.QuizMetadata:
		logger.Info(fmt.Sprintf("Selected option (QuizMetadata): %v\n", model))
		fmt.Printf("Entering QuizMetaData with model: %v\n", model)
	case *screens.DynamicQuizForms:
		logger.Info(fmt.Sprintf("Selected option (DynamicQuizForms): %v\n", model))
		fmt.Printf("Entering DynamicQuizForms with model: %v\n", model)
	default:
		fmt.Printf("Unknown model type: %T\n", m) // Handle unknown model types
	}

	// Exiting program, Application ended
	logger.Info("Application ended")
}
