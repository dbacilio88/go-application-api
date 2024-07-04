package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dbacilio88/go-application-api/pkg/models/response"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Config struct {
}

type Server struct {
	router  *mux.Router
	logger  *zap.Logger
	config  *Config
	handler http.Handler
}

func (s *Server) SetConnectionHealth() {
	s.logger.Info("setting connection health")
}

func NewServer(logger *zap.Logger, config *Config) *Server {
	return &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/actuator/health")
	s.JsonResponse(w, r, &response.AppHealthState)
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/actuator/health", s.healthHandler).Methods("GET")
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	s.registerHandlers()
	s.handler = s.router
	s.logger.Info("starting server")
}

func (s *Server) JsonResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Error("json marshal error", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if data, errorConversion := prettyJson(body); errorConversion != nil {
		s.logger.Error("error converting to pretty json error", zap.Error(errorConversion))
	} else {
		_, _ = w.Write(data)
	}

}

func prettyJson(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
