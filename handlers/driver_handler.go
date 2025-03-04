package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"myproject/config"
	"myproject/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetDrivers(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT * FROM drivers")
	if err != nil {
		http.Error(w, "Error fetching drivers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var drivers []models.Driver

	for rows.Next() {
		var driver models.Driver
		if err := rows.Scan(&driver.ID, &driver.Name, &driver.Phone); err != nil {
			http.Error(w, "Error scanning driver", http.StatusInternalServerError)
			return
		}
		drivers = append(drivers, driver)
	}
	json.NewEncoder(w).Encode(drivers)
}

func GetDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var driver models.Driver
	err = config.DB.QueryRow("SELECT id, name, phone FROM drivers where id=$1", id).
		Scan(&driver.ID, &driver.Name, &driver.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Database Error", http.StatusInternalServerError)
		} else {
			http.Error(w, "Driver not found", http.StatusNotFound)
		}

		return
	}
	json.NewEncoder(w).Encode(driver)

}
func AddDriver(w http.ResponseWriter, r *http.Request) {
	var driver models.Driver

	err := json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	err = config.DB.QueryRow(
		"INSERT INTO drivers (name, phone) VALUES ($1, $2) RETURNING id",
		driver.Name, driver.Phone,
	).Scan(&driver.ID)

	if err != nil {
		http.Error(w, "Failed to create driver", http.StatusBadRequest)
		return
	}
	fmt.Println("decoded driver is: ", driver)
	fmt.Fprintf(w, "Driver added succesfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(driver)
}

func UpdateDriver(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var driver models.Driver
	err = json.NewDecoder(r.Body).Decode(&driver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sqlResultC sql.Result
	sqlResultC, err = config.DB.Exec("UPDATE drivers SET name=$1, phone=$2 WHERE id=$3", driver.Name, driver.Phone, id)
	if err != nil {
		http.Error(w, "Error updating driver", http.StatusInternalServerError)
		return
	}

	var numberOfAffectedRows int64
	numberOfAffectedRows, err = sqlResultC.RowsAffected()

	if err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	if numberOfAffectedRows < 1 {
		http.Error(w, "Driver not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{

		"message": "driver updated successfully",
	})
}

func DeleteDriver(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete call recieved ")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	deleteSql := "DELETE FROM drivers where id=$1"
	var sqlResultC sql.Result

	sqlResultC, err = config.DB.Exec(deleteSql, id)

	if err != nil {
		fmt.Println("error is : ", err)
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	var rowsAffex int64
	rowsAffex, err = sqlResultC.RowsAffected()
	fmt.Println("wga is : ", rowsAffex)

	if err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	if rowsAffex < 1 {
		http.Error(w, "Nothing to delete with this id", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "Driver deleted successfully",
		})
}
