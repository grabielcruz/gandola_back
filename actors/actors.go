package actors

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

func GetActors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	actors := []types.Actor{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM actors ORDER BY id ASC;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		actor := types.Actor{}
		err = rows.Scan(&actor.Id, &actor.Type, &actor.Name, &actor.NationalId, &actor.Address, &actor.Notes, &actor.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		actors = append(actors, actor)
	}
	json_actors, err := json.Marshal(actors)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_actors)
}

func CreateActor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	actor := types.Actor{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &actor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data recibida no es del tipo Actor")
		return
	}
	if actor.Type != "personnel" && actor.Type != "third" && actor.Type != "mine" && actor.Type != "contractee" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el tipo de actor, el cual puede ser 'personal', 'tercero', 'mina' o 'contratante'")
		return
	}
	if actor.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el nombre del actor")
		return
	}
	db := database.ConnectDB()
	defer db.Close()

	insertActorQuery := fmt.Sprintf("INSERT INTO actors (type, name, national_id, address, notes) VALUES ('%v', '%v', '%v', '%v', '%v') RETURNING id, type, name, national_id, address, notes, created_at;", actor.Type, actor.Name, actor.NationalId, actor.Address, actor.Notes)

	rows, err := db.Query(insertActorQuery)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"actors_name_key\"" {
			w.WriteHeader(http.StatusBadRequest)
			rollBackIdQuery := "SELECT setval('actors_id_seq', (SELECT last_value from actors_id_seq) - 1);"
			_, err = db.Query(rollBackIdQuery)
			if err != nil {
				utils.SendInternalServerError(err, w)
				return
			}
			fmt.Fprint(w, "El nombre ya ha sido utilizado")
			return
		}
		utils.SendInternalServerError(err, w)
		return
	}
	for rows.Next() {
		err = rows.Scan(&actor.Id, &actor.Type, &actor.Name, &actor.NationalId, &actor.Address, &actor.Notes, &actor.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	response, err := json.Marshal(actor)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchActor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	actorId := ps.ByName("id")
	actor := types.Actor{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	actorIdNumber, err := strconv.Atoi(actorId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if actorIdNumber <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id de actor no válido")
		return
	}
	if actorIdNumber == 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No puede modificar el actor externo")
		return
	}
	err = json.Unmarshal(body, &actor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada con corresponde con un actor parcial")
		return
	}
	if actor.Type != "personnel" && actor.Type != "third" && actor.Type != "mine" && actor.Type != "contractee" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el tipo de actor, el cual puede ser 'personal', 'tercero', 'mina' o 'contratante'")
		return
	}
	if actor.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el nombre del actor que desea modificar")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var updatedActor types.Actor
	patchActorQuery := fmt.Sprintf("UPDATE actors SET type='%v', name='%v', national_id='%v', address='%v', notes='%v' WHERE id='%v' RETURNING id, type, name, national_id, address, notes, created_at;", actor.Type, actor.Name, actor.NationalId, actor.Address, actor.Notes, actorId)
	actorRow, err := db.Query(patchActorQuery)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"actors_name_key\"" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "El nombre ya ha sido utilizado")
			return
		}
		utils.SendInternalServerError(err, w)
		return
	}
	defer actorRow.Close()
	for actorRow.Next() {
		err = actorRow.Scan(&updatedActor.Id, &updatedActor.Type, &updatedActor.Name, &updatedActor.NationalId, &updatedActor.Address, &updatedActor.Notes, &updatedActor.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if updatedActor.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El actor especificado no existe")
		return
	}

	response, err := json.Marshal(updatedActor)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func GetLastActor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var lastActor types.Actor

	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT * FROM actors ORDER BY id DESC LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&lastActor.Id, &lastActor.Type, &lastActor.Name, &lastActor.NationalId, &lastActor.Address, &lastActor.Notes, &lastActor.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if lastActor.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen más actores")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastActor)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func DeleteActor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	actorId, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if actorId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id de actor no válido")
		return
	}
	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("DELETE FROM actors WHERE id='%v' RETURNING id;", actorId)
	rows, err := db.Query(query)
	if err != nil {
		if err.Error() == "pq: update or delete on table \"actors\" violates foreign key constraint \"transactions_with_balances_actor_fkey\" on table \"transactions_with_balances\"" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "El actor que intenta borrar tiene una o mas transacciones asociadas por lo que no puede ser eliminado")
			return
		}
		utils.SendInternalServerError(err, w)
		return
	}

	defer rows.Close()
	deletedId := types.IdResponse{}
	
	for rows.Next() {
		err = rows.Scan(&deletedId.Id)
		if err != nil {
			utils.SendInternalServerError(err, w)
		return
		}
	}

	if deletedId.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El actor con el id %v no existe", requestedId)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(deletedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}