package server

import (
	"errors"
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//GetPaths returns a list of ids of all the paths in the db
func (r Repository) GetPaths() Paths {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Could not connect to the mongo server: ", err)
	}

	defer session.Close()
	c := session.DB(DBNAME).C(pathColName)
	results := Paths{}
	if err := c.Find(nil).Select(bson.M{"_id": 1}).All(&results); err != nil {
		fmt.Println("Failed to write results: ", err)
	}

	return results
}

//GetPath returns a path of a specified ID
func (r Repository) GetPath(id string) PATH {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Could not connect to the mongo server: ", err)
	}
	defer session.Close()

	c := session.DB(DBNAME).C(pathColName)
	results := PATH{}
	if err := c.FindId(bson.ObjectIdHex(id)).One(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}

	return results
}

//AddPath inserts a Dsm into the DB
func (r Repository) AddPath(path PATH) (*bson.ObjectId, error) {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	path.ID = bson.NewObjectId()
	if err := session.DB(DBNAME).C(pathColName).Insert(path); err != nil {
		fmt.Println("Could not add the path: ", err)
		return nil, errors.New("Adding path failed")
	}

	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Could not connect to the DB")
	}
	return &path.ID, nil
}
