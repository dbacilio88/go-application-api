package controllers

import (
	"fmt"
	"github.com/dbacilio88/go-application-api/pkg/models/response"
	"net/http"
)

// @router /actuator/health
// @method [GET]
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/actuator/health")
	s.JsonResponse(w, r, &response.AppHealthState)
	return
}

// @router /actuator/health/readiness
// @method [GET]
func (s *Server) readinessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/actuator/health/readiness")
	s.JsonResponse(w, r, &response.AppHealthState.Readiness)
	return
}

// @router /actuator/health/liveliness
// @method [GET]
func (s *Server) livelinessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/actuator/health/liveliness")
	s.JsonResponse(w, r, &response.AppHealthState.Liveliness)
	return
}
