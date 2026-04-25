package main

import (
	"context"
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"meeting-planner/backend/internal/db"
	"meeting-planner/backend/internal/handlers"
	"meeting-planner/backend/internal/middleware"
	"meeting-planner/backend/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	backgroundContext := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	database, initializationError := db.Init(backgroundContext, databaseURL)
	if initializationError != nil {
		log.Fatalf("Failed to initialize database: %v", initializationError)
	}
	defer database.Close()

	passwordManager := services.NewPasswordManager(600000, 32, sha256.New)

	handlerInstance := handlers.New(database, passwordManager)
	routeMux := setupRoutes(handlerInstance)

	wrappedHandler := middleware.Recovery(middleware.Logging(middleware.CORS(routeMux)))

	serverAddress := ":8080"
	portEnvVar := os.Getenv("PORT")
	if portEnvVar != "" {
		serverAddress = ":" + portEnvVar
	}

	httpServer := &http.Server{
		Addr:              serverAddress,
		Handler:           wrappedHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("%s", handlers.ToJSONPretty(map[string]any{
			"message": "Server is starting",
			"addr":    httpServer.Addr,
			"time":    time.Now().Format(time.RFC3339),
			"mode":    "development",
			"sdk":     runtime.Version(),
		}))

		if serverError := httpServer.ListenAndServe(); serverError != nil && serverError != http.ErrServerClosed {
			log.Fatalf("Server error: %v", serverError)
		}
	}()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	log.Println("Shutting down server...")

	shutdownContext, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	if shutdownError := httpServer.Shutdown(shutdownContext); shutdownError != nil {
		log.Fatalf("Server forced to shutdown: %v", shutdownError)
	}

	log.Println("Server exited")
}

func setupRoutes(handlerInstance *handlers.Handler) *http.ServeMux {
	routeMux := http.NewServeMux()

	routeMux.HandleFunc("GET /api/health", handlerInstance.HealthcheckEndpoint)
	routeMux.HandleFunc("POST /api/echo/{id}", handlerInstance.EchoEndpoint)

	routeMux.HandleFunc("POST /api/calendars", handlerInstance.CreateCalendarEndpoint)
	routeMux.HandleFunc("POST /api/calendars/{calendar_id}/time-slots", handlerInstance.CreateCalendarTimeSlotsEndpoint)
	routeMux.HandleFunc("GET /api/calendars/{calendar_id}", handlerInstance.GetCalendarEndpoint)
	routeMux.HandleFunc("GET /api/calendars/{calendar_id}/password-protected", handlerInstance.CheckIfCalendarPasswordProtectedEndpoint)
	routeMux.HandleFunc("GET /api/calendars/{calendar_id}/votes", handlerInstance.GetCalendarVotesEndpoint)
	routeMux.HandleFunc("POST /api/calendars/{calendar_id}/time-slots/{time_slot_id}/votes", handlerInstance.VoteOnTimeSlotEndpoint)

	setupStaticFileServer(routeMux)

	return routeMux
}

func setupStaticFileServer(routeMux *http.ServeMux) {
	staticDirectory := "public"

	routeMux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet && request.Method != http.MethodHead {
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if strings.HasPrefix(request.URL.Path, "/api/") {
			http.NotFound(writer, request)
			return
		}

		requestedPath := path.Clean(request.URL.Path)
		if requestedPath == "/" {
			requestedPath = "/index.html"
		}

		requestedPath = strings.TrimPrefix(requestedPath, "/")
		if strings.Contains(requestedPath, "..") {
			http.NotFound(writer, request)
			return
		}

		filePath := path.Join(staticDirectory, requestedPath)
		if _, fileStatError := os.Stat(filePath); fileStatError == nil {
			http.ServeFile(writer, request, filePath)
			return
		}

		indexPath := path.Join(staticDirectory, "index.html")
		if _, indexStatError := os.Stat(indexPath); indexStatError == nil {
			http.ServeFile(writer, request, indexPath)
			return
		}

		http.NotFound(writer, request)
	}))
}
