package server

import (
	"errors"
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//GetDsms returns the list of all Dsms
func (r Repository) GetDsms() Dsms {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to the Mongo server:", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(dsmColName)
	results := Dsms{}
	if err := c.Find(nil).Select(bson.M{"_id": 1}).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

//GetDsm fetches a dsm with a specified ID
//If the dsm cannot be found an empty dsm object
//will be sent instead and the error will be logged
func (r Repository) GetDsm(id string) DSM {
	dsm := DSM{}
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Println("Failed to establish connection to the Mongo server:", err)
		return dsm
	}
	defer session.Close()

	c := session.DB(DBNAME).C(dsmColName)
	if err := c.FindId(bson.ObjectIdHex(id)).One(&dsm); err != nil {
		log.Println("Failed to write results:", err)
	}

	return dsm
}

//AddDsm inserts a Dsm into the DB
func (r Repository) AddDsm(dsm DSM) (*bson.ObjectId, error) {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Could not connect to the DB")
	}
	defer session.Close()

	dsm.ID = bson.NewObjectId()
	if err := session.DB(DBNAME).C(dsmColName).Insert(dsm); err != nil {
		fmt.Println("Could not add dsm: ", err)
		return nil, errors.New("Adding dsm failed")
	}

	return &dsm.ID, nil
}

//UpdateDsm updates a dsm in the database
func (r Repository) UpdateDsm(dsmID string, dsm DSM) bool {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatalln("Failed to connect to the DB")
		return false
	}
	defer session.Close()
	if err := session.DB(DBNAME).C(dsmColName).UpdateId(bson.ObjectIdHex(dsmID), dsm); err != nil {
		log.Fatalln("Could not update the path: ", err)
		return false
	}
	return true
}
