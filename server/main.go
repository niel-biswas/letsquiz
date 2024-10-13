package main

import (
	"context"
	"letsquiz/config"
	"letsquiz/logger"
	"letsquiz/server/database"
	"letsquiz/server/middleware"
	"letsquiz/server/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// start backend HTTP server
	log.Println("Starting backend HTTP server")

	// Load the database configuration from a JSON file
	log.Println("loading dbconfig.json")
	if err := config.LoadConfig("dbconfig.json", true); err != nil {
		logger.Error("Error loading dbconfig.json", "error", err) // Log error if configuration loading fails
		log.Fatalf("could not load config: %v\n", err)            // Exit the program with an error message
	}

	// Initialize the logger
	logger.InitLogger(true)

	// Connect to the database using the loaded configuration
	logger.Info("connecting to database")
	database.ConnectDatabase(config.DbConfig.DbType, config.DbConfig.DbDsn)

	// Set up the router for handling HTTP requests
	logger.Info("setting up router")
	router := routes.NewRouter()
	logger.Info("registering routes to handle requests")
	routes.RegisterRoutes(router)

	// Apply middleware to the router
	logger.Info("Applying middleware: rate limiter")
	var handler http.Handler = router
	handler = middleware.RateLimiter(middleware.Logger(handler)) // Apply rate limiting and logging middleware

	// Apply Okta authentication middleware if enabled in configuration
	if config.AppConfig.EnableOktaAuth {
		logger.Info("Applying middleware: Okta authentication")
		handler = middleware.OktaAuth(handler)
	}

	// Create an HTTP server with the specified address and handler
	srv := &http.Server{
		Addr:    ":8086", // Server will listen on port 8086
		Handler: handler,
	}

	// Start the server in a separate goroutine to avoid blocking
	go func() {
		logger.Info("Starting server on port:8086")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Could not start server", "error", err) // Log error if server fails to start
			log.Fatalf("could not start server: %v\n", err)      // Exit the program with an error message
		}
	}()

	// Set up a channel to catch OS signals for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM) // Notify on interrupt or termination signals

	<-stopChan // Wait for a signal to stop the server

	logger.Info("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err) // Log error if server fails to shut down gracefully
		log.Fatalf("server forced to shutdown: %v", err)        // Exit the program with an error message
	}

	// server stopped, exiting program
	logger.Info("Server exiting")
}
