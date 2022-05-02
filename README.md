# orca-webhook-go

Example of how to build an [Orca Scan WebHook](https://orcascan.com/docs/api/webhooks) endpoint and [Orca Scan WebHook In](https://orcascan.com/guides/how-to-update-orca-scan-from-your-system-4b249706) in using [Go](https://go.dev/).

## Install & Run

First ensure you have [Go](https://go.dev/) installed. If not, follow [this guide](https://go.dev/doc/install).

```bash
# start the project
go run server.go
```

Your WebHook receiver will now be running on port 3000.

You can emulate an Orca Scan WebHook using [cURL](https://dev.to/ibmdeveloper/what-is-curl-and-why-is-it-all-over-api-docs-9mh) by running the following:

```bash
curl --location --request POST 'http://127.0.0.1:3000/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "___orca_sheet_name": "Vehicle Checks",
    "___orca_user_email": "hidden@requires.https",
    "Barcode": "orca-scan-test",
    "Date": "2022-04-19T16:45:02.851Z",
    "Name": Orca Scan Validation Example,
}'
```

### Important things to note

1. Only Orca Scan system fields start with `___`
2. Properties in the JSON payload are an exact match to the  field names in your sheet _(case and space)_

## How this example works

This [example](server.go) work as follows:


```go
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

    // dubug purpose: show in console raw data received
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

	// return success message with 200 status code
	w.WriteHeader(http.StatusOK)
}
```

## Test server locally on Orca Cloud

To expose the server securely from localhost and test it easily on the real Orca Cloud environment you can use [Secure Tunnels](https://ngrok.com/docs/secure-tunnels#what-are-ngrok-secure-tunnels). Take a look at [Ngrok](https://ngrok.com/) or [Cloudflare](https://www.cloudflare.com/).

```bash
ngrok http 3000
```

## Troubleshooting

If you run into any issues not listed here, please [open a ticket](https://github.com/orca-scan/orca-validation-go/issues).

## Examples in other langauges
* [orca-validation-dotnet](https://github.com/orca-scan/orca-validation-dotnet)
* [orca-validation-python](https://github.com/orca-scan/orca-validation-python)
* [orca-validation-go](https://github.com/orca-scan/orca-validation-go)
* [orca-validation-java](https://github.com/orca-scan/orca-validation-java)
* [orca-validation-php](https://github.com/orca-scan/orca-validation-php)
* [orca-validation-node](https://github.com/orca-scan/orca-validation-node)

## History

For change-log, check [releases](https://github.com/orca-scan/orca-validation-python/releases).

## License

&copy; Orca Scan, the [Barcode Scanner app for iOS and Android](https://orcascan.com).