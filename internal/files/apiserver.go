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
	errCh           = make(chan error)
)

type ApiServer struct {
	cfg        config.FilesService
	imgmanager *imgmanager.ImgManager
	publisher  *rabbitmq.ProcessingPublisher
	logger     *logger.Logger
	httpServer *http.Server
	router     *httprouter.Router
}

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

func (s *ApiServer) Run() {
	signal.Notify(doneSignal, syscall.SIGINT, syscall.SIGTERM)

	s.router.POST("/api/v1/image", uploadImage(s.imgmanager, s.publisher, s.logger))

	go s.runHttpServer()

	select {
	case err := <-errCh:
		s.logger.Fatal(err)
		return
	case <-doneSignal:
		s.shutdown(shutdownTimeout)
	}
}

func (a *ApiServer) runHttpServer() {
	a.logger.WithField("ADDR", a.httpServer.Addr).Info("HTTP running")
	errCh <- a.httpServer.ListenAndServe()
}

func (s *ApiServer) shutdown(timeout time.Duration) {
	s.logger.Info("App gracefuly shutdown...")

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		logrus.Error(err)
	}
}
