package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dbacilio88/go-application-api/config"
	"github.com/dbacilio88/go-application-api/config/log"
	"github.com/dbacilio88/go-application-api/pkg/constants"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
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

func NewServer(logger *zap.Logger, config *Config) *Server {
	return &Server{
		router: mux.NewRouter(),
		logger: logger,
		config: config,
	}
}

func (s *Server) SetConnectionHealth() {
	s.logger.Info("setting connection health")
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/actuator/health", s.healthHandler).Methods("GET")
	s.router.HandleFunc("/actuator/health/readiness", s.readinessHandler).Methods("GET")
	s.router.HandleFunc("/actuator/health/liveliness", s.livelinessHandler).Methods("GET")
}

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Reset()
		requestId := uuid.New()
		fmt.Println(requestId)
		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.TraceRequestId, requestId)
		r = r.WithContext(ctx)
		log.LoggerInstance = log.LoggerInstance.With(zap.String(constants.TraceRequestId, requestId.String()))
		r = r.WithContext(log.WithCtx(ctx, log.LoggerInstance))
		next.ServeHTTP(w, r)
	})
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	s.registerHandlers()
	s.handler = s.router
	s.logger.Info("starting server...")
	swaggerUrl := fmt.Sprintf("https://%s:8001/swagger/doc.json", config.Configuration.Microservices.Dns)
	s.logger.Info("swagger url", zap.String("swagger url", swaggerUrl))
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerUrl),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	)).Methods(http.MethodGet)

	go func() {
		if err := http.ListenAndServe(
			fmt.Sprintf("%s:%s", config.Configuration.Microservices.Dns, config.Configuration.Microservices.Port),
			RequestMiddleware(s.handler),
		); err != nil {
			s.logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	s.logger.Info("server started")
	//wait for SIGTERM or SIGINT
	<-stopCh
	s.logger.Info("shutting down server...")
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
