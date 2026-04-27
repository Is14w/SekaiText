package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"sekaitext/backend/internal/config"
)

// NewRouter creates and returns a chi router with all routes and middleware configured.
func NewRouter(cfg *config.AppConfig) http.Handler {
	h := NewHandler(cfg)
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

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
