package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	gotl "github.com/panyam/goutils/template"
	tmplr "github.com/panyam/templar"
)

// TemplateData holds data for template rendering
type TemplateData struct {
	GameID  string
	GameURL string
}

// Server holds our application state
type Server struct {
	templates *tmplr.TemplateGroup
	mux       *http.ServeMux
}

// NewServer creates a new server instance
func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.loadTemplates()
	s.setupRoutes()
	return s
}

// Load HTML templates
func (s *Server) loadTemplates() {
	templates := tmplr.NewTemplateGroup()
	templates.Loader = (&tmplr.LoaderList{}).AddLoader(tmplr.NewFileSystemLoader("./cmd/serve/templates"))
	templates.AddFuncs(gotl.DefaultFuncMap())
	s.templates = templates
}

// Setup HTTP routes
func (s *Server) setupRoutes() {
	// Static files (CSS, JS from static/)
	staticDir := http.Dir("./web/static/")
	staticHandler := http.StripPrefix("/static/", http.FileServer(staticDir))
	s.mux.Handle("/static/", staticHandler)

	// Health check
	s.mux.HandleFunc("/health", s.handleHealth)

	// Main router - handles both home and game pages
	s.mux.HandleFunc("/", s.handleRouting)
}

// Handle routing based on path
func (s *Server) handleRouting(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(r.URL.Path, "/")

	// Home page - shows games list
	if path == "" {
		s.handleGamesPage(w, r)
		return
	}

	// Game page - shows individual game
	if isValidGameID(path) {
		s.handleGamePage(w, r, path)
		return
	}

	// Not found
	http.NotFound(w, r)
}

// Handle games list page
func (s *Server) handleGamesPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	templateFile := "games.html"
	tmpl, err := s.templates.Loader.Load(templateFile, "")
	if err != nil {
		log.Println("Template Load Error: ", templateFile, err)
		fmt.Fprint(w, "Error rendering: ", err.Error())
	} else {
		log.Printf("DEBUG: Successfully loaded template, rendering...")
		err = s.templates.RenderHtmlTemplate(w, tmpl[0], "", nil, nil)
		if err != nil {
			log.Printf("DEBUG: Template render error: %v", err)
			fmt.Fprint(w, "Template render error: ", err.Error())
		} else {
			log.Printf("DEBUG: Template rendered successfully")
		}
	}
}

// Handle individual game page
func (s *Server) handleGamePage(w http.ResponseWriter, r *http.Request, gameID string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := TemplateData{
		GameID:  gameID,
		GameURL: fmt.Sprintf("%s://%s/%s", getScheme(r), r.Host, gameID),
	}

	templateFile := "game.html"
	tmpl, err := s.templates.Loader.Load(templateFile, "")
	if err != nil {
		log.Println("Template Load Error: ", templateFile, err)
		fmt.Fprint(w, "Error rendering: ", err.Error())
	} else {
		log.Printf("DEBUG: Successfully loaded template, rendering...")
		err = s.templates.RenderHtmlTemplate(w, tmpl[0], "", data, nil)
		if err != nil {
			log.Printf("DEBUG: Template render error: %v", err)
			fmt.Fprint(w, "Template render error: ", err.Error())
		} else {
			log.Printf("DEBUG: Template rendered successfully")
		}
	}
}

// Health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Validate game ID format
func isValidGameID(gameID string) bool {
	if len(gameID) == 0 || len(gameID) > 50 {
		return false
	}

	// Allow alphanumeric characters and hyphens
	for _, r := range gameID {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}

	return true
}

// Get request scheme (http/https)
func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		log.Printf("%s %s %d %v",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORS middleware for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Security headers middleware
func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add WASM security headers for all requests (required for WASM)
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}

func main() {
	port := "8080"

	server := NewServer()

	// Add middleware
	handler := loggingMiddleware(
		corsMiddleware(
			securityMiddleware(server),
		),
	)

	log.Printf("🚀 Connect4 server starting on port %s", port)
	log.Printf("📋 Games list: http://localhost:%s/", port)
	log.Printf("🎮 Example game: http://localhost:%s/my-game", port)
	log.Printf("💚 Health check: http://localhost:%s/health", port)
	log.Printf("📁 Serving:")
	log.Printf("   - Static files from: ./static/")
	log.Printf("   - WASM files from: ./web/ -> /static/wasm/")
	log.Printf("   - Stateful proxies in: ./web/gen/")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
