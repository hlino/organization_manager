package api

import (
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

	// TODO: evaluate if I should include these
	originsOk := handlers.AllowedOrigins([]string{"*"})
	s.Router.Use(handlers.CORS(originsOk, methodsOk))

	s.initializeRoutes()
	return nil
}

func (s *Server) Run() {
	log.Info("Listening to port 8084")
	log.Fatal(http.ListenAndServe(":8084", s.Router))
}
