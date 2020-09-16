package service

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Gorynychdo/tdligo.git/internal/model"
	"github.com/go-chi/chi"
)

type HTTPServer struct {
	config *model.Config
	router chi.Router
}

func NewHTTPServer(config *model.Config) *HTTPServer {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			log.Println(errors.Wrap(err, "http handling /"))
		}
	})

	return &HTTPServer{
		config: config,
		router: router,
	}
}

func (s *HTTPServer) ServeHTTP() error {
	return http.ListenAndServe(s.config.HTTPPort, s.router)
}
