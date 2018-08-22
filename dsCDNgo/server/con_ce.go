package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// GetCe GET /ce/{id}
// Has query parameter type, if type = mime a mime multipart message
// is sent. Otherwise, just the json data of the ce is written
func (c *Controller) GetCe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		log.Println("Not a valid ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respType := r.FormValue("type")

	ce := c.Repository.GetCe(id)
	data, err := json.Marshal(ce)
	if err != nil {
		fmt.Println("Could not Marshal data into json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if respType == "mime" {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		//create a new mime part with proper header
		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", "application/json")
		h.Set("Content-Length", strconv.Itoa(len(data)))
		fw, err := writer.CreatePart(h)
		if err != nil {
			fmt.Println("Could not make field for json data: ", err)
			return
		}

		//Copy the data to the body of the mime part
		_, err = io.Copy(fw, bytes.NewReader(data))
		if err != nil {
			fmt.Println("Could not copy data to the body: ", err)
			return
		}

		//Check if there is a teaser asset to add to the message
		if ce.Playlist.Teaser != nil {
			if assetObj, ok := ce.Playlist.Teaser.(Asset); ok {
				writeAssetToPart(assetObj, writer)
			}
		}

		//for every element in the queue get their asset and add it to the message
		var seen = make(map[string]bool) //a map to ensure no duplicate assets are fetched
		for _, element := range ce.Playlist.Queue {
			if assetObj, ok := element.Primary.(Asset); ok {
				id := assetObj.ID.Hex()
				if visited, _ := seen[id]; !visited {
					if err := writeAssetToPart(assetObj, writer); err != nil {
						fmt.Println("Failed to add asset: ", err)
					}
					seen[id] = true
				}
			}
			//Loops through the arrays and adds all their assets
			writeDescToParts(element.Backgrounds, writer, seen)
			writeDescToParts(element.Tracks, writer, seen)
			writeDescToParts(element.Overlays, writer, seen)
		}

		//Set the headers and write the message
		writer.Close() //Can't be defered. Closing adds the ending boundary
		w.Header().Set("Content-Type", writer.FormDataContentType())
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(body.Bytes())
	} else {
		//Set the headers and write the json data
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
	return
}

// GetCes GET /ce
func (c *Controller) GetCes(w http.ResponseWriter, r *http.Request) {
	ces := c.Repository.GetCes()
	data, err := json.Marshal(ces)
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

// AddCe POST /ce
func (c *Controller) AddCe(w http.ResponseWriter, r *http.Request) {
	var ce CE

	err := json.NewDecoder(r.Body).Decode(&ce)
	switch {
	case err == io.EOF:
		log.Println("Request body cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil:
		fmt.Println("Error:", err)
		w.WriteHeader(422) //unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println("Failed encoding error into response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	ceID, err := c.Repository.AddCe(ce)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := bson.M{
		"_id": ceID.Hex(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Could not Marshal data into json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}
