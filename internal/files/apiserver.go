package files

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
)

type ApiServer struct {
	cfg        config.FilesService
	imgmanager *imgmanager.ImgManager
	publisher  *rabbitmq.ProcessingPublisher
	logger     *logger.Logger
}

func NewApiServer(cfg config.FilesService, imgm *imgmanager.ImgManager, publisher *rabbitmq.ProcessingPublisher, logger *logger.Logger) *ApiServer {
	return &ApiServer{
		cfg:        cfg,
		imgmanager: imgm,
		publisher:  publisher,
		logger:     logger,
	}
}

func (s *ApiServer) Run() error {
	r := httprouter.New()
	r.POST("/api/v1/image", uploadImage(s.imgmanager, s.publisher, s.logger))

	return http.ListenAndServe(s.cfg.Port, r)
}
