//requires
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
        bodyParser = require('body-parser'),
        assert = require('assert')

// constants
const mongoURL = 'mongodb://localhost:27017',
	dbName = 'ds',
    assetColName = 'asset',
    dsmColName = 'dsms'
    
const Asset = require('../classes/asset')

let router = express.Router()

router.use(bodyParser.json())
router.use(bodyParser.urlencoded({ extended: true }))

router.route('/asset')
.get((req, res) => {
    MongoClient.connect(mongoURL, { useNewUrlParser: true }, (err, client) => {
        if (err) {
            console.error(error)

            res.status(500).send()

            return
        }

        let db = client.db(dbName)

        let collection = db.collection('asset')

        collection.find({}, { projection: { _id : 1 }}).toArray((err, docs) => {
            if (err) {
                console.error(err);

                res.status(500).send()

                return
            }

            res.status(200).json(docs)
        })
    })
})
.post((req, res) => {
    let asset = req.body

    MongoClient.connect(mongoURL, { useNewUrlParser: true }, (err, client) => {
        if (err) {
            console.error(err)

            res.status(500).send()

            return
        }

        let db = client.db(dbName)

        let assetObj = new Asset(asset.version, asset.url, asset.type, asset.options)

        return add(db, assetColName, assetObj)
        .then((result) => {
            res.status(200).send(`added asset with _id: ${result.insertedId}`)
            return
        })
    })
})
.all((req, res) => {
    res.status(405)
})

//Adds a new item to the database
function add(db, collectionName, obj, opts = { }) {
    return new Promise((resolve, reject) => {
        let collection = db.collection(collectionName)

        collection.insertOne(obj, opts).then(resolve).catch(reject)
    })
}

module.exports = router;