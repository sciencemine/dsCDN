package server

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

//GetDsms GET /dsm
func (c *Controller) GetDsms(w http.ResponseWriter, r *http.Request) {
	dsms := c.Repository.GetDsms() // list of all DSMs
	data, err := json.Marshal(dsms)
	if err != nil {
		fmt.Println("Could not Marshal data into json.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//GetDsm GET /dsm/{id}
func (c *Controller) GetDsm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		log.Println("Not a valid ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dsm := c.Repository.GetDsm(id)
	data, err := json.Marshal(dsm)
	if err != nil {
		fmt.Println("Could not Marshal data into json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// UpdateDsm PUT /dsm/{id}
func (c *Controller) UpdateDsm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dsmID := vars["id"]

	dsm := c.Repository.GetDsm(dsmID)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error UpdateDsm", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Could not close response body: ", err)
	}

	if err := json.Unmarshal(body, &dsm); err != nil {
		w.Header().Set("Conent-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error UpdateDsm unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if ok := c.Repository.UpdateDsm(dsmID, dsm); !ok {
		log.Fatalln("Error updating dsm in the db")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	return
}

// AddDsm POST /dsm
func (c *Controller) AddDsm(w http.ResponseWriter, r *http.Request) {
	var dsm DSM
	err := json.NewDecoder(r.Body).Decode(&dsm)
	switch {
	case err == io.EOF:
		log.Println("Request body cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		w.WriteHeader(422) // unprocessable entity
		log.Println("Error: ", err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error decoding dsm data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	dsmID, err := c.Repository.AddDsm(dsm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Marshal the dsmID into json format
	result := bson.M{
		"_id": dsmID.Hex(),
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
