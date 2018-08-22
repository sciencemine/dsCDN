package server

/*
 * The Respository is responsible for directly interacting with the database
 * it's methods will query the database for data and format it as needed
 * before returning back to the controller
 */

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Repository ...
type Repository struct{}

//SERVER the DB server
const SERVER = "localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "ds"

const dsmColName = "dsms"
const ceColName = "ces"
const pathColName = "paths"
const assetColName = "assets"

//assetIDReplace replaces all the asset ids in a description array with the
//json data related to that id
func assetIDReplace(queueList []desc, a *mgo.Collection) {
	for i := range queueList {
		asset := Asset{}
		if assetID, ok := queueList[i].Asset.(string); ok {
			if err := a.FindId(bson.ObjectIdHex(assetID)).One(&asset); err != nil {
				fmt.Println("Failed to find the asset related to that id: ", err)
			}
		}
		queueList[i].Asset = asset
	}
}
