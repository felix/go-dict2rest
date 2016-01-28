package main

import (
	"golang.org/x/net/dict"
	"log"
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

func getDictClient() (*dict.Client, error) {
	client, err := dict.Dial("tcp", dictServer)
	if err != nil {
		log.Printf("Unable to connect to dict server at %s", dictServer)
		return nil, err
	}
	log.Println("Connected to", dictServer)
	return client, nil
}

func getDictionaries(*dict.Client) ([]dictionary, error) {
	client, err := getDictClient()
	if err != nil {
		log.Printf("Unable to connect to dict server at %s", dictServer)
		return nil, err
	}

	defer client.Close()

	dictArr, err := client.Dicts()
	if err != nil {
		return nil, err
	}

	var dicts []dictionary
	for _, d := range dictArr {
		dicts = append(dicts, dictionary{d.Name, d.Desc})
	}
	return dicts, nil
}
