package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"sekaitext/backend/internal/config"
	"sekaitext/backend/internal/service"
)

// NewRouter creates and returns a chi router with all routes and middleware configured.
func NewRouter(cfg *config.AppConfig) http.Handler {
	logBuf := service.NewLogBuffer(200)
	h := NewHandler(cfg, logBuf)
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(requestLogger(logBuf))

	// CORS - allow all origins in dev
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(corsHandler.Handler)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Debug
		r.Get("/debug/logs", h.DebugLogs)
		r.Post("/debug/save", h.DebugSaveLogs)

		// Story navigation
		r.Route("/story", func(r chi.Router) {
			r.Get("/types", h.StoryTypes)
			r.Get("/sorts", h.StorySorts)
			r.Get("/index", h.StoryIndex)
			r.Get("/chapter", h.StoryChapter)
			r.Get("/json-path", h.JsonPath)
			r.Post("/load", h.StoryLoad)
			r.Post("/parse-local", h.StoryParseLocal)
			r.Post("/load-local", h.StoryLoadLocal)
			r.Post("/download-json", h.DownloadStoryJSON)
			r.Get("/download-progress", h.DownloadProgress)
		})

		// Translation file operations
		r.Route("/translation", func(r chi.Router) {
			r.Post("/create", h.TranslationCreate)
			r.Post("/load", h.TranslationLoad)
			r.Post("/save", h.TranslationSave)
			r.Post("/check-lines", h.CheckLines)
		})

		// Editor operations
		r.Route("/editor", func(r chi.Router) {
			r.Post("/change-text", h.ChangeText)
			r.Post("/add-line", h.AddLine)
			r.Post("/remove-line", h.RemoveLine)
			r.Post("/compare", h.Compare)
			r.Post("/replace-brackets", h.ReplaceBrackets)
		})

		// Text check
		r.Post("/check/text", h.CheckText)

		// Flashback
		r.Route("/flashback", func(r chi.Router) {
			r.Post("/analyze", h.FlashbackAnalyze)
			r.Get("/clue-hints", h.ClueHints)
			r.Get("/voice-clues", h.VoiceClues)
		})

		// Voice
		r.Get("/voice/url", h.VoiceURL)

		// Speaker
		r.Post("/speaker/count", h.SpeakerCount)

		// Settings
		r.Get("/settings", h.GetSettings)
		r.Put("/settings", h.PutSettings)

		// Update (CDN refresh)
		r.Post("/update", h.Update)
		r.Get("/update/progress", h.UpdateProgress)

		// Assets
		r.Route("/assets", func(r chi.Router) {
			r.Get("/characters", h.Characters)
			r.Get("/character-icon/{index}", h.CharacterIcon)
			r.Get("/units", h.Units)
			r.Get("/areas", h.Areas)
		})
	})

	return r
}

// requestLogger is a custom middleware that logs requests and writes to the log buffer.
func requestLogger(buf *service.LogBuffer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			latency := time.Since(start)
			msg := fmt.Sprintf("%s %s %d %s",
				r.Method, r.URL.Path, ww.Status(), latency.Round(time.Millisecond))
			buf.Write(msg)
			fmt.Println("[server]", msg)
		})
	}
}
