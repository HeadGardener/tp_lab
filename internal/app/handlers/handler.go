package handlers

import (
	"encoding/json"
	"github.com/HeadHardener/tp_lab/internal/app/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

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
			// r.Delete("/delete/{worker_id}", h.deleteWorker)
		})
		//  переделать базу данных (добавить таблицу удаленных работников и убрать переходную таблицу, также создать
		//	референс от таблицы документов к работникам без каскадного удаления)
		r.Route("/gsm", func(r chi.Router) {
			r.Put("/{document_id}", h.updateDocument)
			r.Delete("/{document_id}", h.deleteDocument)
		})
	})

	r.Route("/gsm", func(r chi.Router) {
		r.Use(h.identifyUser)
		r.Post("/", h.createDocument)
		r.Get("/", h.getAllDocuments)
		r.Get("/{document_id}", h.getDocumentByID)
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
