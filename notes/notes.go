package notes

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

func GetNotes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	notes := []types.Note{}
	db := database.ConnectDB()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM notes ORDER BY id;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		note := types.Note{}
		err := rows.Scan(&note.Id, &note.Description, &note.Urgency, &note.Attended, &note.CreatedAt, &note.AttendedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		notes = append(notes, note)
	}
	json_notes, err := json.Marshal(notes)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_notes)
}

func CreateNote(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	note := types.Note{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	err = json.Unmarshal(body, &note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data recibida no es del tipo Nota")
		return
	}
	if note.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la descripción de la nota")
		return
	}
	if note.Urgency != "low" && note.Urgency != "medium" && note.Urgency != "high" && note.Urgency != "critical" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la urgencia de la nota, la cual puede ser baja, media, alta o crítica")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	insertNoteQuery := fmt.Sprintf("INSERT INTO notes (description, urgency) VALUES ('%v', '%v') RETURNING id, description, urgency, attended, created_at, attended_at;", note.Description, note.Urgency)

	rows, err := db.Query(insertNoteQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	for rows.Next() {
		err = rows.Scan(&note.Id, &note.Description, &note.Urgency, &note.Attended, &note.CreatedAt, &note.AttendedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	response, err := json.Marshal(note)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	noteId := ps.ByName("id")
	note := types.Note{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}
	noteIdNumber, err := strconv.Atoi(noteId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El parametro id debe ser un número")
		return
	}
	if noteIdNumber <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id de nota no válido")
		return
	}
	err = json.Unmarshal(body, &note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada con corresponde a una nota")
		return
	}

	if note.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la descripción de la nota que desea modificar")
		return
	}
	if note.Urgency != "low" && note.Urgency != "medium" && note.Urgency != "high" && note.Urgency != "critical" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la urgencia de la nota, la cual puede ser baja, media, alta o crítica")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var updatedNote types.Note
	patchNoteQuery := fmt.Sprintf("UPDATE notes SET description='%v', urgency='%v' WHERE id='%v' RETURNING id, description, urgency, attended, created_at, attended_at;", note.Description, note.Urgency, noteIdNumber)
	noteRow, err := db.Query(patchNoteQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer noteRow.Close()
	for noteRow.Next() {
		err = noteRow.Scan(&updatedNote.Id, &updatedNote.Description, &updatedNote.Urgency, &updatedNote.Attended, &updatedNote.CreatedAt, &updatedNote.AttendedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if updatedNote.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La nota especificada no existe")
		return
	}

	response, err := json.Marshal(updatedNote)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeleteNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	noteId, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if noteId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id de nota no válido")
		return
	}
	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("DELETE FROM notes WHERE id='%v' RETURNING id;", noteId)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	deletedId := types.IdResponse{}

	for rows.Next() {
		err := rows.Scan(&deletedId.Id)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if deletedId.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La nota con el id %v no existe", requestedId)
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

func AttendNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	noteId, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if noteId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id de nota no válido")
		return
	}
	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("UPDATE notes SET attended='TRUE', attended_at=CURRENT_TIMESTAMP WHERE id='%v' RETURNING id, description, urgency, attended, created_at, attended_at;", noteId)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	attendedNote := types.Note{}

	for rows.Next() {
		err := rows.Scan(&attendedNote.Id, &attendedNote.Description, &attendedNote.Urgency, &attendedNote.Attended, &attendedNote.CreatedAt, &attendedNote.AttendedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if attendedNote.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La nota con el id %v no existe", requestedId)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(attendedNote)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func UnattendNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}
	noteId, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	if noteId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id de nota no válido")
		return
	}
	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("UPDATE notes SET attended='FALSE' WHERE id='%v' RETURNING id, description, urgency, attended, created_at, attended_at;", noteId)
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	unattendedNote := types.Note{}

	for rows.Next() {
		err := rows.Scan(&unattendedNote.Id, &unattendedNote.Description, &unattendedNote.Urgency, &unattendedNote.Attended, &unattendedNote.CreatedAt, &unattendedNote.AttendedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if unattendedNote.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La nota con el id %v no existe", requestedId)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(unattendedNote)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}

func GetLastNoteId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lastNoteId := types.IdResponse{
		Id: -1,
	}
	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT id FROM notes LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&lastNoteId.Id); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}
	if lastNoteId.Id == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen más notas")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastNoteId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}