package server

import (
	"errors"
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//GetCes returns the list of all ces from the db
func (r Repository) GetCes() Ces {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Could not connect to the mongo server: ", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(ceColName)
	results := Ces{}
	if err := c.Find(nil).Select(bson.M{"_id": 1}).All(&results); err != nil {
		fmt.Println("Failed to write results: ", err)
	}

	return results
}

//GetCe fetches a ce with a specified ID
func (r Repository) GetCe(id string) CE {
	results := CE{}
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to the Mongo server:", err)
		return results
	}
	defer session.Close()

	c := session.DB(DBNAME).C(ceColName)
	if err := c.FindId(bson.ObjectIdHex(id)).One(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	a := session.DB(DBNAME).C(assetColName)
	if results.Playlist.Teaser != nil {
		asset := Asset{}
		if assetID, ok := results.Playlist.Teaser.(string); ok {
			if err := a.FindId(bson.ObjectIdHex(assetID)).One(&asset); err != nil {
				fmt.Println("Failed to find the asset related to that id: ", err)
			}
			results.Playlist.Teaser = asset
		}

	}

	for i := range results.Playlist.Queue {
		asset := Asset{}
		if assetID, ok := results.Playlist.Queue[i].Primary.(string); ok {
			if err := a.FindId(bson.ObjectIdHex(assetID)).One(&asset); err != nil {
				fmt.Println("Failed to find the asset related to that id: ", err)
			}
			results.Playlist.Queue[i].Primary = asset
		}
		assetIDReplace(results.Playlist.Queue[i].Backgrounds, a)
		assetIDReplace(results.Playlist.Queue[i].Tracks, a)
		assetIDReplace(results.Playlist.Queue[i].Overlays, a)
	}
	return results
}

//AddCe adds a Ce to the DB
func (r Repository) AddCe(ce CE) (*bson.ObjectId, error) {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	ce.ID = bson.NewObjectId()
	if err := session.DB(DBNAME).C(ceColName).Insert(ce); err != nil {
		fmt.Println("Could not add the ce :", err)
		return nil, errors.New("Adding ce failed")
	}

	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Could not connect to the DB")
	}
	return &ce.ID, nil
}
