package server

//GetAssets returns the list of all Assets
import (
	"errors"
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//GetAssets gets a list of all asset IDs from the database
func (r Repository) GetAssets() Assets {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to the Mongo server:", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(assetColName)
	results := Assets{}
	if err := c.Find(nil).Select(bson.M{"_id": 1}).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

//AddAsset inserts a asset into the DB
func (r Repository) AddAsset(asset Asset) (*bson.ObjectId, error) {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	asset.ID = bson.NewObjectId()
	if err := session.DB(DBNAME).C(assetColName).Insert(asset); err != nil {
		fmt.Println("Could not add asset: ", err)
		return nil, errors.New("Adding asset failed")
	}

	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Could not connect to the DB")
	}
	return &asset.ID, nil
}

//UpdatePath updates a path by adding a relation to the relations set
func (r Repository) UpdatePath(pathID string, change interface{}) bool {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer session.Close()
	if err := session.DB(DBNAME).C(pathColName).UpdateId(bson.ObjectIdHex(pathID), change); err != nil {
		fmt.Println("Could not update the path: ", err)
		return false
	}
	return true
}
