package files

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/config"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

type ApiServer struct {
	cfg        config.FilesService
	imgmanager *imgmanager.ImgManager
}

func NewApiServer(cfg config.FilesService, imgm *imgmanager.ImgManager) *ApiServer {
	return &ApiServer{
		cfg:        cfg,
		imgmanager: imgm,
	}
}

func (s *ApiServer) Run() error {
	r := httprouter.New()
	r.POST("/api/v1/image", uploadImage(s.imgmanager))

	return http.ListenAndServe(s.cfg.Port, r)
}
