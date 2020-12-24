package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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
	csvFile, err := os.Open("Arquivo/q1_catalog.csv")
	checkErr(err)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'

	var empresa []CSV

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			checkErr(err)
		}
		empresa = append(empresa, CSV{
			Id:      line[0],
			Name:    line[1],
			ZipCode: line[2],
		})
	}

	empresaJSON, err := json.Marshal(empresa)
	checkErr(err)
	fmt.Println(string(empresaJSON))
}
