package handlers

import (
	"encoding/json"
	"github.com/HeadHardener/tp_lab/internal/app/models/websocket"
	"github.com/HeadHardener/tp_lab/internal/app/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	service   *services.Service
	errLogger *zap.Logger
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service:   service,
		errLogger: newLogger(),
	}
}

type WebSocketHandler struct {
	hub       *ws.Hub
	errLogger *zap.Logger
}

func NewWSHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub:       hub,
		errLogger: newLogger(),
	}
}

func InitRoutes(h *Handler, wsh *WebSocketHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/api", func(r chi.Router) {
		// auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-in", h.signIn)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(h.identifyUser)
			r.Use(h.checkRole)
			r.Route("/worker", func(r chi.Router) {
				r.Post("/sign-up", h.createWorker)
				r.Get("/get-all/", h.getAllWorkers)
				r.Get("/get/{worker_id}", h.getWorkerByID)
				r.Put("/update/{worker_id}", h.updateWorker)
			})
			r.Route("/gsm", func(r chi.Router) {
				r.Put("/{document_id}", h.updateDocument)
				r.Delete("/{document_id}", h.deleteDocument)
			})
		})

		r.Route("/token", func(r chi.Router) {
			r.Use(h.identifyUser)
			// check token for front
			r.Get("/check", h.isValid)
			r.Get("/get-me", h.getMe)
		})

		r.Route("/gsm", func(r chi.Router) {
			r.Use(h.identifyUser)
			r.Post("/", h.createDocument)
			r.Get("/", h.getAllDocuments)
			r.Get("/{document_id}", h.getDocumentByID)
			r.Get("/my", h.getDocumentsWithWorkerID)
		})

		r.Route("/chat", func(r chi.Router) {
			// r.Use(h.identifyUser)
			r.Post("/create-room", wsh.createRoom)
			r.Get("/join-room/{room_id}", wsh.joinRoom)
			r.Get("/get-rooms", wsh.getRooms)
			r.Get("/get-clients/{room_id}", wsh.getClients)
		})
	})

	return r
}

func newLogger() *zap.Logger {
	rawJSON := []byte(`{
	  "level": "error",
	  "encoding": "json",
	  "outputPaths": ["stdout"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	defer logger.Sync()
	return logger
}
