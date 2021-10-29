package files

import (
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/internal/rabbitmq"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

func uploadImage(imgm *imgmanager.ImgManager, publisher *rabbitmq.ProcessingPublisher) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		file, _, err := r.FormFile("image")
		if err != nil {
			internalError(rw, err)
			return
		}
		defer file.Close()
		stat, err := imgm.ReadAndSaveNewImage(file)
		if err != nil {
			if errors.Is(err, imgmanager.ErrUnsupportedFormat) {
				http.Error(rw, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
				return
			}
			log.Println(err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = publisher.Publish(stat)
		if err != nil {
			internalError(rw, err)
			return
		}

		fmt.Fprintf(rw, "file %s was uploaded", stat.ID)
	}
}

func internalError(rw http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
