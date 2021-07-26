package handle_uploads

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example.com/backend_gandola_soft/utils"
	"github.com/julienschmidt/httprouter"
)

func UploadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file, header, err := r.FormFile("image")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer file.Close()

	extension := filepath.Ext(header.Filename)
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