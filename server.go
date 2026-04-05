package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// validationHandler is called every time Orca Scan sends a barcode scan to your server.
// It reads the scan data, applies your validation logic, and tells Orca Scan whether
// to save the data, reject it, or change it before saving.
func validationHandler(w http.ResponseWriter, r *http.Request) {

	// Read the raw request body sent by Orca Scan.
	// This is the JSON data containing the barcode scan and all sheet field values.
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the JSON body into a map so you can access any field by name.
	// Fields starting with ___ are Orca system fields (e.g. ___orca_sheet_name).
	// All other fields match your sheet column names exactly (case and spaces matter).
	// For example: data["Barcode"], data["Name"], data["___orca_sheet_name"]
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the value of the Name field from the incoming data.
	// The .(string) part safely converts the value to a string.
	// If the field is missing or not a string, name will be empty ("").
	name, _ := data["Name"].(string)

	// ---------------------------------------------------------------
	// OPTION 1: Reject the scan and show an error dialog in the app.
	// Return HTTP 400 with an ___orca_message to block the save and
	// display the message to the user. They must dismiss the dialog
	// before they can try again.
	// ---------------------------------------------------------------
	if len(name) > 20 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"___orca_message": map[string]interface{}{
				"display": "dialog",
				"type":    "error",
				"title":   "Invalid Name",
				"message": "Name cannot be longer than 20 characters",
			},
		})
		return
	}

	// ---------------------------------------------------------------
	// OPTION 2: Modify the data before it saves.
	// Return HTTP 200 with only the fields you want to change.
	// Orca Scan will update those fields and allow the save.
	// ---------------------------------------------------------------
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"Name": name, // example: you could trim whitespace or reformat the value
	// })
	// return

	// ---------------------------------------------------------------
	// OPTION 3: Show a success notification (green banner in the app).
	// The data still saves - this just gives the user feedback.
	// Return HTTP 200 with an ___orca_message to show the notification.
	// ---------------------------------------------------------------
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"___orca_message": map[string]interface{}{
	// 		"display": "notification",
	// 		"type":    "success",
	// 		"message": "Barcode scanned successfully",
	// 	},
	// })
	// return

	// ---------------------------------------------------------------
	// SECURITY: Verify the request came from your specific Orca sheet.
	// Set a secret in Orca Scan (Integrations > Events API > Secret)
	// then check it matches here before trusting the data.
	// ---------------------------------------------------------------
	// secret := r.Header.Get("orca-secret")
	// if secret != os.Getenv("ORCA_SECRET") {
	// 	http.Error(w, "", http.StatusUnauthorized)
	// 	return
	// }

	// All good - return HTTP 204 to allow the data to save with no changes.
	// HTTP 204 means "success, no content" - Orca Scan will save the data as-is.
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Use the PORT environment variable if set, otherwise default to 8888.
	// This makes the server easy to deploy to cloud platforms that set PORT for you.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	// Register the validationHandler function to handle all incoming POST requests.
	http.HandleFunc("/", validationHandler)

	fmt.Println("Listening on port " + port + ". Ready for Orca Scan requests.")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
