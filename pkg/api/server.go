package api

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize() error {
	s.Router = mux.NewRouter()
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})

	originsOk := handlers.AllowedOrigins([]string{"*"})
	s.Router.Use(handlers.CORS(originsOk, methodsOk))

	s.initializeRoutes()
	return nil
}

func (s *Server) Run(port int) {
	log.Infof("Listening to port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router))
}
