//requires
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
		bodyParser = require('body-parser')
		path = require('path')
        assert = require('assert')

// constants
const mongoURL = 'mongodb://localhost:27017',
	dbName = 'ds',
	dsmColName = 'dsms'

router.use(bodyParser)

let router = express.Router()

router.route('/asset')
.post((req, res) => {
    let body = req.body

    MongoClient.connect(mongoURL, (err, client) => {
        if (err) {
            console.error(err)

            res.status(500)

            return
        }

        let db = client.db(dbName)

        let asset = JSON.parse(body)

        let assetObj = new Asset(asset.version, asset.url, asset.type, asset)

        return add(db, assetColName, assetObj)
        .then((result) => {
            res.status(200)
            return
        })

    })
})