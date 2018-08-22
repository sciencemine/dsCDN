package server

/*
 * The controller methods for the /path endpoint
 */

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetPaths GET /path
func (c *Controller) GetPaths(w http.ResponseWriter, r *http.Request) {
	//Get a list of the paths from the repository
	paths := c.Repository.GetPaths()
	data, err := json.Marshal(paths) //turn the path struct into valid json
	if err != nil {
		fmt.Println("Could not marshal data into json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//set the headers and write the data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//GetPath GET /path/{id}
func (c *Controller) GetPath(w http.ResponseWriter, r *http.Request) {
	//get the path ID from the route
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		log.Printf("%s is not a valid ID", id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Get the path from the repository by ID
	path := c.Repository.GetPath(id)
	data, err := json.Marshal(path)
	if err != nil {
		fmt.Println("Could not Marshal data into json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Set the headers and write the data
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// AddPath POST /path
func (c *Controller) AddPath(w http.ResponseWriter, r *http.Request) {
	var path PATH

	err := json.NewDecoder(r.Body).Decode(&path)
	switch {
	case err == io.EOF:
		log.Println("Response Body cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddPath unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	pathID, err := c.Repository.AddPath(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Marshal the pathID into json format
	result := bson.M{
		"_id": pathID.Hex(),
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

//UpdatePath PUT /path/{id}
func (c *Controller) UpdatePath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathID := vars["id"]

	var relation Relation

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read teh body of the request
	if err != nil {
		log.Fatalln("Error UpdatePath", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error UpdatePath", err)
	}
	if err := json.Unmarshal(body, &relation); err != nil {
		w.Header().Set("Conent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdatePath unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	change := bson.M{
		"$addToSet": bson.M{
			"relations": relation,
		},
	}
	if success := c.Repository.UpdatePath(pathID, change); !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}
