package trucks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"example.com/backend_gandola_soft/database"
	"example.com/backend_gandola_soft/types"
	"example.com/backend_gandola_soft/utils"
	"github.com/julienschmidt/httprouter"
)

func GetTrucks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	trucks := []types.Truck{}
	photos_array := []string{}
	var photos string
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM trucks ORDER BY id;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		truck := types.Truck{}
		err := rows.Scan(&truck.Id, &truck.Name, &truck.Data, &photos, &truck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		if len(photos) > 0 || photos == "null" {
			err = json.Unmarshal([]byte(photos), &photos_array)
			if err != nil {
				utils.SendInternalServerError(err, w)
				return
			}
		} else {
			photos_array = []string{}
		}
		truck.Photos = photos_array
		trucks = append(trucks, truck)
	}
	json_trucks, err := json.Marshal(trucks)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_trucks)
}

func CreateTruck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	truck := types.Truck{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &truck)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data recibida no es del tipo Camión")
		return
	}
	if truck.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el nombre del camión")
		return
	}
	if truck.Data == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la data del camión")
		return
	}

	if len(truck.Photos) == 0 {
		truck.Photos = []string{}
	}

	json_photos, err := json.Marshal(truck.Photos)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	insertTruckQuery := fmt.Sprintf("INSERT INTO trucks (name, data, photos) VALUES ('%v', '%v', '%v') RETURNING id, name, data, photos, created_at", truck.Name, truck.Data, string(json_photos))
	insertedTruck := types.Truck{}

	rows, err := db.Query(insertTruckQuery)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"trucks_name_key\"" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "El nombre del camión ya ha sido utilizado")
			return
		}
		utils.SendInternalServerError(err, w)
		return
	}
	var newPhotos string
	photos_array := []string{}
	for rows.Next() {
		err = rows.Scan(&insertedTruck.Id, &insertedTruck.Name, &insertedTruck.Data, &newPhotos, &insertedTruck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if len(newPhotos) > 0 || newPhotos == "null" {
		err = json.Unmarshal([]byte(newPhotos), &photos_array)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	} else {
		photos_array = []string{}
	}
	insertedTruck.Photos = photos_array

	response, err := json.Marshal(insertedTruck)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchTruck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	truckId := ps.ByName("id")
	truck := types.Truck{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	truckIdNumber, err := strconv.Atoi(truckId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El parametro id debe ser un número")
		return
	}
	if truckIdNumber <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Id de camión no válido")
		return
	}
	err = json.Unmarshal(body, &truck)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada no corresponde con un camión")
		return
	}

	if truck.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el nombre del camión")
		return
	}
	if truck.Data == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la data del camión")
		return
	}

	if len(truck.Photos) == 0 {
		truck.Photos = []string{}
	}

	json_photos, err := json.Marshal(truck.Photos)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var updatedTruck types.Truck
	patchTruckQuery := fmt.Sprintf("UPDATE trucks SET name='%v', data='%v', photos='%v' WHERE id='%v' RETURNING Id, name, data, photos, created_at;", truck.Name, truck.Data, string(json_photos), truckIdNumber)
	truckRow, err := db.Query(patchTruckQuery)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"trucks_name_key\"" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "El nombre del camión ya ha sido utilizado")
			return
		}
		utils.SendInternalServerError(err, w)
		return
	}
	defer truckRow.Close()
	var newPhotos string
	photos_array := []string{}
	for truckRow.Next() {
		err = truckRow.Scan(&updatedTruck.Id, &updatedTruck.Name, &updatedTruck.Data, &newPhotos, &updatedTruck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if updatedTruck.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El camión especificado no existe")
		return
	}

	if len(newPhotos) > 0 || newPhotos == "null" {
		err = json.Unmarshal([]byte(newPhotos), &photos_array)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	} else {
		photos_array = []string{}
	}
	updatedTruck.Photos = photos_array

	response, err := json.Marshal(updatedTruck)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeleteTruck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	truckId := ps.ByName("id")
	if truckId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	truckIdNumber, err := strconv.Atoi(truckId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El parametro id debe ser un número")
		return
	}
	if truckIdNumber <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Id de camión no válido")
		return
	}

	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("DELETE FROM trucks WHERE id='%v' RETURNING id;", truckIdNumber)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	deletedTruckId := types.IdResponse{}

	for rows.Next() {
		err := rows.Scan(&deletedTruckId.Id)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if deletedTruckId.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El camión con el id %v no existe", truckIdNumber)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(deletedTruckId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func GetLastTruck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var lastTruck types.Truck
	photos_array := []string{}
	var photos string
	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT * FROM trucks ORDER BY id DESC LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&lastTruck.Id, &lastTruck.Name, &lastTruck.Data, &photos, &lastTruck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		if len(photos) > 0 || photos == "null" {
			err = json.Unmarshal([]byte(photos), &photos_array)
			if err != nil {
				utils.SendInternalServerError(err, w)
				return
			}
		} else {
			photos_array = []string{}
		}
		lastTruck.Photos = photos_array
	}

	if lastTruck.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen mas camiones")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastTruck)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}