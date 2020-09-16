package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Gorynychdo/tdligo.git/internal/model"
	"github.com/Gorynychdo/tdligo.git/internal/tdclient"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

const InvalidBodyError = "invalid body"

type HTTPServer struct {
	config   *model.Config
	router   chi.Router
	tdClient *tdclient.TDClient
}

func NewHTTPServer(config *model.Config, tc *tdclient.TDClient) *HTTPServer {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	return &HTTPServer{
		config:   config,
		router:   router,
		tdClient: tc,
	}
}

func (s *HTTPServer) ServeHTTP() error {
	s.router.Post("/send", s.sendMessage)
	return http.ListenAndServe(s.config.HTTPPort, s.router)
}

func (s *HTTPServer) sendMessage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(errors.Wrap(err, "reading request body"))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var message model.OutgoingMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println(errors.Wrapf(err, "unmarshall message: %s", string(body)))
		http.Error(w, InvalidBodyError, http.StatusBadRequest)
		return
	}
	if message.ChatID == 0 || message.Text == "" {
		log.Printf("%s: %s", InvalidBodyError, string(body))
		http.Error(w, InvalidBodyError, http.StatusBadRequest)
		return
	}
	if err = s.tdClient.SendMessage(message); err != nil {
		log.Println(errors.Wrapf(err, "send message"))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
