package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/dict"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func render(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(body)
}

func formatError(e error) jsonError {
	parts := strings.SplitN(e.Error(), "[", 2)
	log.Printf("Error %s", parts[0])
	parts = strings.SplitN(parts[0], " ", 2)
	code, err := strconv.Atoi(parts[0])
	if err != nil {
		return jsonError{
			Code:    500,
			Message: "Error parsing error o_0",
		}
	}
	return jsonError{
		Code:    code,
		Message: strings.TrimSpace(parts[1]),
	}
}

func formatDefinitions(defs []*dict.Defn) ([]definition, error) {
	definitions := make([]definition, len(defs))
	for i, def := range defs {
		definitions[i] = definition{
			Dictionary: def.Dict.Desc,
			Word:       def.Word,
			Definition: string(def.Text[:]),
		}
	}
	return definitions, nil
}

func dictDatabases(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, err := getDictClient()
	if err != nil {
		log.Printf("Unable to connect to dict server at %s", dictServer)
		render(w, 500, jsonError{420, "Server temporarily unavailable"})
		return
	}

	defer client.Close()

	dicts, err := getDictionaries(client)
	if err != nil {
		render(w, 400, formatError(err))
		return
	}
	render(w, 200, dicts)
}

func dictDefine(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, err := getDictClient()
	if err != nil {
		render(w, 500, jsonError{420, "Server temporarily unavailable"})
		return
	}
	defer client.Close()

	word := ps.ByName("word")
	queryValues := r.URL.Query()
	d := queryValues.Get("dict")

	var dict string

	if d != "" {
		dict = d
	} else {
		dict = "*"
	}

	defs, err := client.Define(dict, word)
	if err != nil {
		render(w, 400, formatError(err))
		return
	}

	definitions, err := formatDefinitions(defs)
	if err != nil {
		render(w, 500, err)
		return
	}
	render(w, 200, definitions)
}
