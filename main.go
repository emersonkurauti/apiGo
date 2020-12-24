package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type companyStruct struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Zip     string `json:"zip"`
	WebSite string `json:"website"`
}

func checkErr(err error) {
	if err != nil {
		log.Panic("ERROR: " + err.Error())
	}
}

func main() {
	database, err := sql.Open("sqlite3", "./apiGo.db")
	checkErr(err)
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS companies (id TEXT PRIMARY KEY, name TEXT, zip TEXT, website TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO cmpanies (id, name, zip, website) VALUES (?, ?, ?, ?)")

	csvFile, err := os.Open("Arquivo/q1_catalog.csv")
	checkErr(err)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			checkErr(err)
		}

		query, err := database.Query("SELECT id FROM companies WHERE id = " + line[0])
		checkErr(err)
		exist := query.Next()
		query.Close()

		if !exist {
			statement.Exec(line[0], strings.ToUpper(line[1]), line[2], "")
		}
	}

	database.Close()

	router := mux.NewRouter()
	router.HandleFunc("/companies", GetCompanyByNameAndZip).Methods("GET")
	router.HandleFunc("/companies", PosCompany).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetCompanyByNameAndZip(w http.ResponseWriter, r *http.Request) {
	var company companyStruct

	paramName := r.URL.Query().Get("name")
	paramZip := r.URL.Query().Get("zip")

	database, err := sql.Open("sqlite3", "./apiGo.db")
	checkErr(err)
	query, err := database.Query("SELECT id, name, zip, website FROM companies WHERE name like '%" + strings.ToUpper(paramName) + "%' AND zip = " + paramZip)
	checkErr(err)

	var id string
	var name string
	var zip string
	var website string

	query.Next()
	query.Scan(&id, &name, &zip, &website)
	query.Close()

	company = companyStruct{
		Id:      id,
		Name:    name,
		Zip:     zip,
		WebSite: website,
	}

	database.Close()
	json.NewEncoder(w).Encode(company)
}

func PosCompany(w http.ResponseWriter, r *http.Request) {
	var company companyStruct
	_ = json.NewDecoder(r.Body).Decode(&company)

	database, err := sql.Open("sqlite3", "./apiGo.db")
	checkErr(err)
	query, err := database.Query("SELECT id FROM companies WHERE name = '" + strings.ToUpper(company.Name) + "' AND zip = " + company.Zip)
	checkErr(err)

	var id string

	if !query.Next() {
		return
	}

	query.Scan(&id)
	query.Close()

	statement, err := database.Prepare("UPDATE companies SET website=? WHERE id=?")
	checkErr(err)
	statement.Exec(company.WebSite, id)
	company.Id = id

	database.Close()
	json.NewEncoder(w).Encode(company)
}
