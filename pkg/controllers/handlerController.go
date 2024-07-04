package controllers

import (
	"fmt"
	"github.com/dbacilio88/go-application-api/pkg/models/response"
	"github.com/dbacilio88/go-application-api/pkg/node"
	"net/http"
)

func (s *node.Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/actuator/health")
	s.JsonResponse(w, r, &response.AppHealthState)
	return
}
