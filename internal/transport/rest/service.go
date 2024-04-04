package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run(port string, ginEngine *gin.Engine) error {
	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: ginEngine,
	}

	return s.server.ListenAndServe()
}
