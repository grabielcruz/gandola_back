package handle_uploads

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"example.com/backend_gandola_soft/types"
	"example.com/backend_gandola_soft/utils"
	"github.com/julienschmidt/httprouter"
)

func UploadBill(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer file.Close()
	fmt.Println(header.Filename)

	validImage := false
	extension := strings.ToLower(filepath.Ext(header.Filename))
	for _, v := range types.ImageTypes {
		if extension == v {
			validImage = true
			break
		}
	}
	if !validImage {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El archivo del tipo %v no es una imagen reconocida", extension)
		return
	}

	id := ps.ByName("id")
	date := time.Now()
	year, month, day := date.Local().Date()
	name := fmt.Sprintf("public/bills/factura_%v_%v-%v-%v%v", id, month, day, year, extension)

	// tempFile, err := ioutil.TempFile("public/bills", name)
	// if err != nil {
	// 	utils.SendInternalServerError(err, w)
	// 	return
	// }
	// defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	newFile, err := os.Create(name)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer newFile.Close()

	newFile.Write(fileBytes)
	fileName := newFile.Name()
	fmt.Fprint(w, fileName)
}

func UploadTrucksPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	photos := []string{}
	err := r.ParseMultipartForm(10 << 20) // 10Mb
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	formdata := r.MultipartForm

	files := formdata.File["images"]

	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		defer file.Close()

		extension := strings.ToLower(filepath.Ext(f.Filename))
		id := ps.ByName("id")
		date := time.Now()
		year, month, day := date.Local().Date()
		name := fmt.Sprintf("photo_*_%v_%v-%v-%v%v", id, month, day, year, extension)

		tempFile, err := ioutil.TempFile("public/trucks", name)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		defer tempFile.Close()
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		_, err = io.Copy(tempFile, file)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		photos = append(photos, f.Filename)
	}
	response, err := json.Marshal(photos)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(response))
}
