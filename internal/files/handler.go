package files

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
)

func uploadImage(imgm *imgmanager.ImgManager, publisher *rabbitmq.ProcessingPublisher, logger *logger.Logger) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		file, _, err := r.FormFile("image")
		if err != nil {
			internalError(rw, logger, err)
			return
		}
		defer file.Close()
		stat, err := imgm.ReadAndSaveNewImage(file)
		if err != nil {
			if errors.Is(err, imgmanager.ErrUnsupportedFormat) {
				http.Error(rw, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
				return
			}
			internalError(rw, logger, err)
			return
		}

		err = publisher.Publish(stat)
		if err != nil {
			internalError(rw, logger, err)
			return
		}

		fmt.Fprintf(rw, "file %s was uploaded", stat.ID)
	}
}

func internalError(rw http.ResponseWriter, logger *logger.Logger, err error) {
	logger.Println(err)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
