package customjson

import (
	"encoding/json"
	"net/http"
)

// Encode will encode data. If status == 0, it will automatically write StatusOk
func Encode(w http.ResponseWriter, data interface{}, status int) {
	// Write content type header
	w.Header().Add("content-type", "application/json")

	if status != 0 {
		w.WriteHeader(status)
	}

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes response
	err := encoder.Encode(&data)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}

// Decode will decode both request and response.
func Decode(i interface{}, data interface{}) error {
	//declare the decoder
	var dec *json.Decoder

	//create decoder based on interface type
	switch i.(type) {
	case *http.Request:
		dec = json.NewDecoder(i.(*http.Request).Body)
	case *http.Response:
		dec = json.NewDecoder(i.(*http.Response).Body)
	}

	//decodes interface
	if err := dec.Decode(data); err != nil {
		return err
	}

	return nil
}
