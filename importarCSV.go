package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type CSV struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	ZipCode string `json:"zipcode"`
}

func checkErr(err error) {
	if err != nil {
		log.Panic("ERROR: " + err.Error())
	}
}

func main() {
	database, _ := sql.Open("sqlite3", "./apiGo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS companies (id INTEGER PRIMARY KEY, name TEXT, zipcode TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO companies (id, name, zipcode) VALUES (?, ?, ?)")

	csvFile, err := os.Open("Arquivo/q1_catalog.csv")
	checkErr(err)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			checkErr(err)
		}

		query, _ := database.Query("SELECT id FROM companies WHERE id = " + line[0])
		exist := query.Next()
		query.Close()

		if !exist {
			statement.Exec(line[0], strings.ToUpper(line[1]), line[2])
		}
	}
}
