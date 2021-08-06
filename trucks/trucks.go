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
		err := rows.Scan(&truck.Id, &truck.Name, &truck.Data, &truck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		photos_query := fmt.Sprintf("SELECT url FROM truck_photos WHERE truck=%v;", truck.Id)
		photo_rows, err := db.Query(photos_query)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		defer photo_rows.Close()
		for photo_rows.Next() {
			var url string
			err := photo_rows.Scan(&url)
			if err != nil {
				utils.SendInternalServerError(err, w)
				return
			}
			truck.Photos = append(truck.Photos, url)
		}
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

	db := database.ConnectDB()
	defer db.Close()

	insertTruckQuery := fmt.Sprintf("INSERT INTO trucks (name, data) VALUES ('%v', '%v') RETURNING id, name, data, created_at", truck.Name, truck.Data)
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
	for rows.Next() {
		err = rows.Scan(&insertedTruck.Id, &insertedTruck.Name, &insertedTruck.Data, &insertedTruck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
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

	db := database.ConnectDB()
	defer db.Close()

	var updatedTruck types.Truck
	patchTruckQuery := fmt.Sprintf("UPDATE trucks SET name='%v', data='%v' WHERE id='%v' RETURNING Id, name, data, created_at;", truck.Name, truck.Data, truckIdNumber)
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
	for truckRow.Next() {
		err = truckRow.Scan(&updatedTruck.Id, &updatedTruck.Name, &updatedTruck.Data, &updatedTruck.Created_At)
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
		err := rows.Scan(&lastTruck.Id, &lastTruck.Name, &lastTruck.Data, &lastTruck.Created_At)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
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