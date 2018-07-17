//requires
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
		bodyParser = require('body-parser')
		path = require('path')
        assert = require('assert')

//local requires
const PATH = require('../classes/path'),
        utils = require('./util')

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

    MongoClient.connect(mongoURL, {useNewUrlParser : true}, (err, client) => {
		if (err) {
			console.error(err)

			res.status(500).send()

			return
        }

        const db = client.db(dbName)

        utils.query(db, pathColName, pathID)
        .then((doc) => {
            let pathObj = doc
            res.status(200).json(pathObj)
        })
        .catch(utils.queryErr)
    })
})
.put((req, res) => {
    let pathID = new mongo.ObjectID(req.params.id)
    let relation = req.body

    MongoClient.connect(mongoURL, {useNewUrlParser: true}, (err, client) => {
        if (err) {
            console.error(err)

            res.status(500).send()

            return
        }

        const db = client.db(dbName)
        utils.addRelation(db, pathColName, relation, pathID)
        .then((results) => {
            res.status(200).send()
        })
        .catch((err) => {
            console.error(err)
        })
        
    })
})
.all((req, res) => {
    res.status(405).send()
})


router.route('/path')
.get((req, res) => {
    MongoClient.connect(mongoURL, {useNewUrlParser: true}, (err, client) => {
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
.post((req, res) => {
    let path = req.body

    MongoClient.connect(mongoURL, {useNewUrlParser: true}, (err, client) => {
        if (err) {
            console.error(err)

            res.status(500).send()

            return
        }

        let db = client.db(dbName)

        pathObj = new PATH(path._id, path.model, path.relations)

        return utils.add(db, pathColName, pathObj)
        .then((results) => {
            console.log(`added a path with the _id ${results.insertedId}`)
            res.status(201).json({"pathId": `${results.insertedId}`})
        })
        .catch((err) => {
            console.error(err)
        })
    })
})
.all((req, res) => {
    res.status(405).send()
})

module.exports = router