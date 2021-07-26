package bills

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"example.com/backend_gandola_soft/database"
	"example.com/backend_gandola_soft/types"
	"example.com/backend_gandola_soft/utils"
	"github.com/julienschmidt/httprouter"
)

func GetBills(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bills :=[]types.Bill{}
	db := database.ConnectDB();
	defer db.Close();
	rows, err := db.Query("SELECT bills.id, url, date, charged, company, name, national_id, bills.created_at FROM bills INNER JOIN actors ON bills.company = actors.id;")
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		bill := types.Bill{}
	  err = rows.Scan(&bill.Id, &bill.Url, &bill.Date, &bill.Charged, &bill.Company.Id, &bill.Company.Name, &bill.Company.NationalId, &bill.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
		bills = append(bills, bill)
	}
	json_bills, err := json.Marshal(bills)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_bills)
}

func CreateBill(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bill := types.Bill{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
		return
	}

	err = json.Unmarshal(body, &bill)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data recibida no corresponde con una factura")
		return
	}

	// if bill.Url == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "Debe especificar el url del archivo de la factura")
	// 	return
	// }

	if bill.Company.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la compañía a la que pertenece la factura")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var companyId int
	getCompanyIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id=%v;", bill.Company.Id)
	companyIdRow, err := db.Query(getCompanyIdQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer companyIdRow.Close()
	for companyIdRow.Next() {
		err = companyIdRow.Scan(&companyId)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if companyId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La compañía especificada no existe")
		return
	}

	var insertedId int
	var insertBillQuery string
	if  bill.Date == "" {
		insertBillQuery = fmt.Sprintf("INSERT INTO bills (url, company, charged) VALUES ('%v', '%v', '%v') RETURNING id;", bill.Url, bill.Company.Id, bill.Charged)
	} else {
		_, err := time.Parse(time.RFC3339, bill.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "La fecha de la factura no tiene un formato válido")
			return
		}
		insertBillQuery = fmt.Sprintf("INSERT INTO bills (url, date, company, charged) VALUES ('%v', '%v', '%v', '%v') RETURNING id;", bill.Url, bill.Date, bill.Company.Id, bill.Charged)
	}

	rowsInsertedId, err := db.Query(insertBillQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsInsertedId.Close()
	for rowsInsertedId.Next() {
		err = rowsInsertedId.Scan(&insertedId)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	insertedBill := types.Bill{}
	retrieveBillQuery := fmt.Sprintf("SELECT bills.id, url, date, charged, company, name, national_id, bills.created_at FROM bills INNER JOIN actors ON bills.company = actors.id WHERE bills.id='%v';", insertedId)
	rowsRetreivedBill, err := db.Query(retrieveBillQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsRetreivedBill.Close()
	for rowsRetreivedBill.Next() {
		err = rowsRetreivedBill.Scan(&insertedBill.Id, &insertedBill.Url, &insertedBill.Date, &insertedBill.Charged, &insertedBill.Company.Id, &insertedBill.Company.Name, &insertedBill.Company.NationalId, &insertedBill.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	response, err := json.Marshal(insertedBill)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func PatchBill(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	billsId, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	newBill := types.Bill{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se pudo leer el cuerpo de la petición")
	}
	
	err = json.Unmarshal(body, &newBill)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La data enviada no corresponde con una factura")
		return
	}

	if billsId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El id de la factura debe ser mayor a cero")
		return
	}

	if newBill.Company.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar la compañía a la que pertenece la factura")
		return
	}

	db := database.ConnectDB()
	defer db.Close()

	var companyId int
	getCompanyIdQuery := fmt.Sprintf("SELECT id FROM actors WHERE id=%v;", newBill.Company.Id)
	companyIdRow, err := db.Query(getCompanyIdQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer companyIdRow.Close()
	for companyIdRow.Next() {
		err = companyIdRow.Scan(&companyId)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if companyId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La compañía especificada no existe")
		return
	}

	var updatedId int
	var updateBillQuery string
	if newBill.Date == "" {
		updateBillQuery = fmt.Sprintf("UPDATE bills SET url='%v', company='%v', charged='%v' WHERE id='%v' RETURNING id;", newBill.Url, newBill.Company.Id, newBill.Charged, billsId)
	} else {
		_, err := time.Parse(time.RFC3339, newBill.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "La fecha de la factura no tiene un formato válido")
			return
		}
		updateBillQuery = fmt.Sprintf("UPDATE bills SET url='%v', date='%v', company='%v', charged='%v' WHERE id='%v' RETURNING id;", newBill.Url, newBill.Date, newBill.Company.Id, newBill.Charged, billsId)
	}

	rowsUpdatedId, err := db.Query(updateBillQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rowsUpdatedId.Close()
	for rowsUpdatedId.Next() {
		err = rowsUpdatedId.Scan(&updatedId)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if updatedId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "La factura solicitada no existe")
		return
	}

	updatedBill := types.Bill{}
	retrieveBillQuery := fmt.Sprintf("SELECT bills.id, url, date, charged, company, name, national_id, bills.created_at FROM bills INNER JOIN actors ON bills.company = actors.id WHERE bills.id='%v';", updatedId)
	rowsRetreivedBill, err := db.Query(retrieveBillQuery)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	defer rowsRetreivedBill.Close()
	for rowsRetreivedBill.Next() {
		err = rowsRetreivedBill.Scan(&updatedBill.Id, &updatedBill.Url, &updatedBill.Date, &updatedBill.Charged, &updatedBill.Company.Id, &updatedBill.Company.Name, &updatedBill.Company.NationalId, &updatedBill.CreatedAt)
		if err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	response, err := json.Marshal(updatedBill)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeleteBill(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedId := ps.ByName("id")
	if requestedId == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Debe especificar el parametro id en la petición de borrado")
		return
	}

	id, err := strconv.Atoi(requestedId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}

	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "El id de la factura debe ser mayor a cero")
		return
	}

	db := database.ConnectDB()
	defer db.Close()
	query := fmt.Sprintf("DELETE FROM bills WHERE id='%v' RETURNING id;", id)
	rows, err := db.Query(query)
	if err != nil {
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
		fmt.Fprintf(w, "La factura con el id %v no existe", requestedId)
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

func GetLastBillId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lastBillId := types.IdResponse{
		Id: -1,
	}
	db := database.ConnectDB()
	defer db.Close()
	query := "SELECT id FROM bills ORDER BY id DESC LIMIT 1;"
	rows, err := db.Query(query)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&lastBillId.Id); err != nil {
			utils.SendInternalServerError(err, w)
			return
		}
	}

	if lastBillId.Id == -1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No existen más facturas")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(lastBillId)
	if err != nil {
		utils.SendInternalServerError(err, w)
		return
	}
	w.Write(response)
}