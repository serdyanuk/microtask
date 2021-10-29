package files

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
	"github.com/serdyanuk/microtask/pkg/logger"
)

// fileSizeLimit 		= 5mb
const fileSizeLimit = 5 << 20

func uploadImage(imgm *imgmanager.ImgManager, publisher *rabbitmq.ProcessingPublisher, logger *logger.Logger) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		rw.Header().Add("Access-Control-Allow-Origin", "*")

		r.Body = http.MaxBytesReader(rw, r.Body, fileSizeLimit)
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}
		defer file.Close()

		stat, err := imgm.ReadAndSaveNewImage(file)
		if err != nil {
			if imgmanager.IsUnknowFormatErr(err) {
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

		logger.Infof("new image saved to disk %s", stat)

		fmt.Fprintf(rw, "OK")
	}
}

func internalError(rw http.ResponseWriter, logger *logger.Logger, err error) {
	logger.Error(err)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
