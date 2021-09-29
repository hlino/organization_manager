package api

import (
	"organization_manager/pkg/api/controllers"
)

func (s *Server) initializeRoutes() {
	router := s.Router.PathPrefix("/api/v1").Subrouter()
	router.HandleFunc("/organizations", controllers.CreateOrganization).Methods("POST")
	router.HandleFunc("/organizations", controllers.GetOrganizations).Methods("GET")
}

