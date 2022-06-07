package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
)

type OrcaBarcode struct {
	Barcode					string
    Date 					string
    Name					string
    ___orca_sheet_name		string
    ___orca_user_email		string
}

func validationHandler(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse JSON data
	var barcode OrcaBarcode
	jsonErr := json.Unmarshal([]byte(body), &barcode)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		http.Error(w, jsonErr.Error(), 500)
		return
	}

    // debug purpose: show in console raw data received
	fmt.Println(barcode)

	// NOTE:
	// orca system fields start with ___
	// you can access the value of each field using the field name (data.Name, data.Barcode, data.Location)
	name := barcode.Name

	// validation example
	if(len(name) > 20){
		// return error message with json format
		w.Write([]byte(`{
			"title": "Invalid Name",
			"message": "Name must be less than 20 characters"}
			`))
		return
	}

	// return HTTP Status 200 with no body
	w.Write([]byte(""))
}

func main() {
    http.HandleFunc("/", validationHandler)

    fmt.Println("Server started at port 3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

