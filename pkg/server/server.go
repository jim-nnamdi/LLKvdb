package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type GracefulShutdownServer struct {
	HTTPListenAddr      string
	PutHandler          http.Handler
	ReadHandler         http.Handler
	ReadKeyRangeHandler http.Handler
	BatchPutHandler     http.Handler
	DeleteHandler       http.Handler

	httpServer     *http.Server
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	IdleTimeout    time.Duration
	HandlerTimeout time.Duration
}

func (server *GracefulShutdownServer) getRouter() *mux.Router {
	router := mux.NewRouter()

	router.SkipClean(true)
	router.Handle("/put", server.PutHandler)
	router.Handle("/read/{key}", server.ReadHandler)
	router.Handle("/readkeyrange", server.ReadKeyRangeHandler)
	router.Handle("/batchput", server.BatchPutHandler)
	router.Handle("/delete", server.DeleteHandler)
	return router
}

func (server *GracefulShutdownServer) Start() {
	router := server.getRouter()
	server.httpServer = &http.Server{
		Addr:         server.HTTPListenAddr,
		WriteTimeout: server.WriteTimeout,
		ReadTimeout:  server.ReadTimeout,
		IdleTimeout:  server.IdleTimeout,
		Handler:      router,
	}
	server.httpServer.ListenAndServe()
}
