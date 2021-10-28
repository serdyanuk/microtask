package files

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/serdyanuk/microtask/pkg/imgmanager"
)

func uploadImage(imgm *imgmanager.ImgManager) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		file, _, err := r.FormFile("image")
		if err != nil {
			resError(rw, err)
			return
		}
		defer file.Close()
		filename, err := imgm.SaveImage(file)
		if err != nil {
			resError(rw, err)
			return
		}
		err = imgm.DoResize(filename)
		if err != nil {
			resError(rw, err)
			return
		}
		fmt.Fprintf(rw, "file %s was uploaded", filename)
	}
}

func resError(rw http.ResponseWriter, err error) {
	log.Println(err)
	rw.Write([]byte(err.Error()))
}
