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

type definition struct {
	Dictionary string `json:"dictionary"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}
type dictionary struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Render(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(body)
}

func FormatError(e error) jsonError {
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

func FormatDefinitions(defs []*dict.Defn) ([]definition, error) {
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

func Databases(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var dicts []dictionary
	for _, d := range dictMap {
		dicts = append(dicts, dictionary{d.Name, d.Desc})
	}
	if len(dicts) == 0 {
		Render(w, 200, jsonError{554, "No databases present"})
		return
	}
	Render(w, 200, dicts)
}

func Define(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	word := ps.ByName("word")
	queryValues := r.URL.Query()
	d := queryValues.Get("dict")

	var dict string

	_, ok := dictMap[d]
	if d != "" && !ok {
		Render(w, 400, jsonError{500, "Invalid database"})
		return
	} else if d != "" && ok {
		dict = d
	} else {
		dict = "*"
	}
	defs, err := client.Define(dict, word)
	if err != nil {
		Render(w, 400, FormatError(err))
		return
	}
	//log.Printf("DEFINE '%s' from '%s' found %d definitions", word, dict, len(defs))

	definitions, err := FormatDefinitions(defs)
	if err != nil {
		Render(w, 500, err)
		return
	}
	Render(w, 200, definitions)
}
