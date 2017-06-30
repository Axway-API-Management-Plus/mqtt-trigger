package tools

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

/*
func ReadJSONApi(w http.ResponseWriter, r *http.Request, model interface{}) error {
	if err := jsonapi.UnmarshalPayload(r.Body, model); err != nil {
		log.Errorln("Not a JSONApi", err)
		http.Error(w, err.Error(), 422) // unprocessable entity
		return err
	}

	return nil
}
*/

func ReadJSON(w http.ResponseWriter, r *http.Request, model interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(model); err != nil {
		log.Errorln("Not a JSON", err)
		http.Error(w, err.Error(), 422) // unprocessable entity
		return err
	}

	return nil
}

/*
func WriteJSONApi(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	if err := jsonapi.MarshalOnePayload(w, model); err != nil {
		http.Error(w, err.Error(), 500)
		panic(err)
	}
}*/

func WriteJSON(w http.ResponseWriter, model interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var b []byte
	b, err := json.Marshal(model)
	if err != nil {
		http.Error(w, err.Error(), 500)
		panic(err)
	}
	w.Write(b)
}

/*
func WriteJSONApiMany(w http.ResponseWriter, models []interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	if err := jsonapi.MarshalManyPayload(w, models); err != nil {
		http.Error(w, err.Error(), 500)
		panic(err)
	}
}*/

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not Implemented"))
}

func HttpErrorResponse(code int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}
