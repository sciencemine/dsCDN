package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

//GetAssets GET /asset
func (c *Controller) GetAssets(w http.ResponseWriter, r *http.Request) {
	assets := c.Repository.GetAssets()
	data, err := json.Marshal(assets)
	if err != nil {
		fmt.Println("Could not marshal data into json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//AddAsset POST /asset
func (c *Controller) AddAsset(w http.ResponseWriter, r *http.Request) {
	var asset Asset

	err := json.NewDecoder(r.Body).Decode(&asset)
	switch {
	case err == io.EOF:
		log.Println("Request body cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		w.WriteHeader(422) // unprocessable entity
		log.Println("Error: ", err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddAsset unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	assetID, err := c.Repository.AddAsset(asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Marshal the assetID into json format
	result := bson.M{
		"_id": assetID.Hex(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Could not Marshal data into json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Send back the response and the id of the inserted dsm
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}
