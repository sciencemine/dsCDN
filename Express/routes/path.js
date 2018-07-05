//requires
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
		bodyParser = require('body-parser')
		path = require('path')
		assert = require('assert')

//local requires
const PATH = require('../classes/path')

let router = express.Router()

// constants
const mongoURL = 'mongodb://localhost:27017',
    dbName = 'ds',
    pathColName = 'paths'

router.use(bodyParser.json())
router.use(bodyParser.urlencoded({ extended: true }))

router.route('/path/:id') 
.get((req, res) => {
    let pathID = new mongo.ObjectID(req.params.id) 

    MongoClient.connect(mongoURL, (err, client) => {
		if (err) {
			console.error(err)

			res.status(500).send()

			return
        }

        const db = client.db(dbName)

        query(db, 'paths', pathID)
        .then((doc) => {
            let pathObj = doc
            let promises = []
            res.status(200).json(pathObj)
        })
        .catch(queryErr)
    })
})
.all((req, res) => {
    res.status(405).send()
})

router.route('/path')
.get((req, res) => {
    
    MongoClient.connect(mongoURL, (err, client) => {
        if (err) {
            console.error(err)

            res.status(500).send()

            return
        }

        let db = client.db(dbName)

        let collection = db.collection(pathColName)

        collection.find({}, {_id: 1}).toArray((err, doc) => {
            if (err){
                console.error(err)

                res.status(500).send()

                return
            }

            res.status(200).json(doc)
        })
    })
})
.all((req, res) => {
    res.status(405).send()
})

function query(db, collectionName, id) {
	return new Promise((resolve, reject) => {
		let collection = db.collection(collectionName)

		collection.find({_id: id }).toArray((err, doc) => {
			if (err) {
				reject(err)
			}

			resolve(doc[0])
		})
	})
}

function queryErr(err) {
	console.error(err)
}

module.exports = router