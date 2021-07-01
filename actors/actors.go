package actors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"example.com/backend_gandola_soft/database"
	"github.com/julienschmidt/httprouter"
)

type Actor struct {
	Id          int
	Name        string
	Description string
	CreatedAt   string
}

type PartialActor struct {
	Name        string
	Description string
}

type IdResponse struct {
	Id int
}

func GetActors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	actors := []Actor{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM actors ORDER BY id ASC;")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		actor := Actor{}
		err = rows.Scan(&actor.Id, &actor.Name, &actor.Description, &actor.CreatedAt)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		actors = append(actors, actor)
	}
	json_actors, err := json.Marshal(actors)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_actors)
}

func CreateActor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	actor := Actor{}
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
	if actor.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el nombre del actor")
		return
	}
	if actor.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la descripción del actor")
		return
	}
	db := database.ConnectDB()
	defer db.Close()

	insertActorQuery := fmt.Sprintf("INSERT INTO actors (name, description) VALUES ('%v', '%v') RETURNING id, name, description, created_at;", actor.Name, actor.Description)

	rows, err := db.Query(insertActorQuery)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"actors_name_key\"" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "El nombre ya ha sido utilizado")
			return
		}
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		err = rows.Scan(&actor.Id, &actor.Name, &actor.Description, &actor.CreatedAt)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	response, err := json.Marshal(actor)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchActor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	actorId := ps.ByName("id")
	partialActor := PartialActor{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &partialActor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada con corresponde con un actor parcial")
		return
	}
	if partialActor.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el nombre del actor que desea modificar")
		return
	}
	if partialActor.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la descripción del actor que desea modificar")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var updatedActor Actor
	patchActorQuery := fmt.Sprintf("UPDATE actors SET name='%v', description='%v' WHERE id='%v' RETURNING id, name, description, created_at;", partialActor.Name, partialActor.Description, actorId)
	actorRow, err := db.Query(patchActorQuery)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"actors_name_key\"" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "El nombre ya ha sido utilizado")
			return
		}
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error del servidor")
		return
	}
	defer actorRow.Close()
	for actorRow.Next() {
		err = actorRow.Scan(&updatedActor.Id, &updatedActor.Name, &updatedActor.Description, &updatedActor.CreatedAt)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error del servidor")
		}
	}

	if updatedActor.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El actor especificado no existe")
		return
	}

	response, err := json.Marshal(updatedActor)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error del servidor")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func GetLastActorId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var lastActorId IdResponse

	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT id FROM actors ORDER BY id DESC LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error del servidor")
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&lastActorId.Id)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error del servidor")
			return
		}
	}

	if lastActorId.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen más actores")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastActorId)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error del servidor")
		return
	}
	w.Write(response)
}