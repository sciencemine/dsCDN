// This route is for defining dsm requests and stuff

//requires
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
		bodyParser = require('body-parser')
		path = require('path')
		assert = require('assert')

//local requires
const Asset = require('../classes/asset'),
    CE = require('../classes/ce'),
    DSM = require('../classes/dsm');

let router = express.Router()

// constants
const mongoURL = 'mongodb://localhost:27017',
	dbName = 'ds',
	dsmColName = 'dsms'

router.use(bodyParser.json())
router.use(bodyParser.urlencoded({ extended: true }))

//route returns a dsm of a particular id
router.route('/dsm/:id')
.get((req, res) => {
	let dsmID = new mongo.ObjectID(req.params.id)

	//get the dsm from mongo database
	MongoClient.connect(mongoURL, (err, client) => {
		if (err) {
			console.error(err)

			res.status(500)

			return
		}

		const db = client.db(dbName)

		//query time! woo!
		query(db, 'dsms', dsmID) 
		.then((doc) => { //I *promise* to come back with things :P
			let dsmObj = doc
			let promises = []
			res.status(200).json(dsmObj)
		})
		.catch(queryErr) //I broke my promise guys, I'm sorry
	})
})

.all((req, res) => {
	res.status(405)
})


//a route that returns a list of all the dsms in the database
router.route('/dsm')
.get((req, res) => {
	MongoClient.connect(mongoURL, { useNewUrlParser: true }, (err, client) => {
		if (err) {
			console.error(err)

			res.status(500)

			return
		}

		const db = client.db(dbName)

		let collection = db.collection(dsmColName)

		//find all dsms and grab their _id, title, and desc
		collection.find({}, { 
				 _id: 1,
				 title: 1,
				 description: 1
			}
		).toArray((err, docs) => { //add it to an array named docs
			if (err) { //check for errors
				console.error(err)

				res.status(500).send()

				return
			}

			res.status(200).json(docs) //return the array
		})
	})	
})
.post((req, res) => {
	let dsm = req.body
	
	MongoClient.connect(mongoURL, {useNewUrlParser: true},  (err, client) => {
		if (err) {
			console.error(err)

			res.status(500)

			return
		}

		const db = client.db(dbName)

		//let dsm = JSON.parse(body) 

		let dsmObj = new DSM(dsm._id, dsm.title, dsm.description,
				dsm.version, dsm.author, dsm.config,
				dsm.stylesheet, dsm.style, dsm.contributors,
				dsm.idle_backgrounds, dsm.video_select_backgrounds,
				dsm.ce_set, dsm.attributes)

		//validate that all the ce elements in the set are in the database
		let valid_ce_set = true
		for (let ce in dsmObj.ce_set) {
			db.collectionName.find({_id: ce}, {_id: 1}).limit(1)
			.catch(() => {
				valid_ce_set = false
			})
		}

		if (valid_ce_set) {
			return add(db, dsmColName, dsmObj)
			.then((results) => {
				console.log(`added dsm with _id: ${results.insertedId}`)
				res.status(200).send(`added dsm with _id: ${results.insertedId}`)
				client.close()
			})
			.catch((err) => {
				console.error(err)
				res.status(500)
			})
		}
		res.status(500)
		return
	})
})
.all((req, res) => {
	res.status(405);
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

function add(db, collectionName, obj, opts = { }) {
	return new Promise((resolve, reject) => {
		let collection = db.collection(collectionName)

		collection.insertOne(obj, opts).then(resolve).catch((err) => { console.log(err); reject(err); });
	})
}

module.exports = router