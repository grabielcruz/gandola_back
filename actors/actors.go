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
	Id int
	Name string
	Description string
	CreatedAt string
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
	defer rows.Close();
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