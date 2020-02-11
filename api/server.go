package api

import (
	"context"
	"fmt"
	"sync"

	zl "demo-app/v1/logger"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	once   sync.Once
	server *Server
)

//Server - REST & Websocket server
type Server struct {
	restServerPort int
	restServer     *echo.Echo
	log            zl.ZeroLogger
}

//GetServerInstance - create REST & websocket server.
//This is a singleton object
func GetServerInstance(restPort int, log zl.ZeroLogger) *Server {
	once.Do(func() {
		server = createServer(restPort, log)
	})

	return server
}

//StartServer - start the REST server & start listening for websocket
func (s *Server) StartServer() error {

	s.restServer.POST("/events", s.eventHandler)
	s.restServer.GET("/ws", s.handleWebsocket)

	//start rest server
	go func(port int) {
		addr := fmt.Sprintf(":%d", port)
		s.restServer.Logger.Fatal(s.restServer.Start(addr))
	}(s.restServerPort)

	return nil
}

//Shutdown - stop the server
func (s *Server) Shutdown() {
	s.restServer.Shutdown(context.Background())
}

//start - start the REST server
func createServer(restPort int, log zl.ZeroLogger) *Server {
	apiServer := echo.New()
	apiServer.Use(middleware.Recover())

	return &Server{
		restServerPort: restPort,
		restServer:     apiServer,
		log:            log,
	}
}
