package files

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	doneSignal      = make(chan os.Signal, 1)
	shutdownTimeout = time.Second * 10
)

// ApiServer represents api server.
type ApiServer struct {
	cfg        config.FilesService
	imgmanager *imgmanager.ImgManager
	publisher  *rabbitmq.ProcessingPublisher
	logger     *logger.Logger
	httpServer *http.Server
	router     *httprouter.Router
}

// NewApiServer is used to create api server.
func NewApiServer(cfg config.FilesService, imgm *imgmanager.ImgManager, publisher *rabbitmq.ProcessingPublisher, logger *logger.Logger) *ApiServer {
	router := httprouter.New()
	httpServer := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	return &ApiServer{
		cfg:        cfg,
		imgmanager: imgm,
		publisher:  publisher,
		logger:     logger,
		httpServer: httpServer,
		router:     router,
	}
}

// Run is used to start the application.
func (s *ApiServer) Run() {
	signal.Notify(doneSignal, syscall.SIGINT, syscall.SIGTERM)

	s.router.POST("/api/v1/image", uploadImage(s.imgmanager, s.publisher, s.logger))

	go s.runServer()

	<-doneSignal
	s.shutdown(shutdownTimeout)
}

// runHttpServer is used to start http server.
func (a *ApiServer) runServer() {
	a.logger.WithField("ADDR", a.httpServer.Addr).Info("HTTP running")

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.logger.Fatal(err)
	}
}

// shutdown is used to shutdown http server.
func (s *ApiServer) shutdown(timeout time.Duration) {
	s.logger.Info("Server gracefuly shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		logrus.Fatal("Server shutdown error: ", err)
	}
}
