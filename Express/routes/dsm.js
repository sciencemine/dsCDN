// This route is for defining dsm requests and stuff

//requires
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
		fs = require('fs'),
		bodyParser = require('body-parser')
		path = require('path')
		assert = require('assert')


let router = express.Router()

// constants
const mongoURL = 'mongodb://localhost:27017',
	dbName = 'ds',
	dsmColName = 'dsms'

router.use(bodyParser)

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
	MongoClient.connect(mongoURL, (err, client) => {
		if (err) {
			console.error(err)

			res.status(500)

			return
		}

		const db = client.db(dbName)

		let collection = db.collection('dsms')

		collection.find({}, { 
			projection: {
				 _id: 1,
				 title: 1,
				 description: 1
			} 
		}).toArray((err, docs) => {
			if (err) {
				console.error(err)

				res.status(500)

				return
			}

			res.status(200).json(docs)
		})
	})	
})
.post((req, res) => {
	let body = req.body;
	
	MongoClient.connect(mongoURL, {useNewUrlParser: true},  (err, client) => {
		if (err) {
			console.error(err)

			res.status(500)

			return
		}

		const db = client.db(dbName)

		let dsm = JSON.parse(body) 

		let dsmObj = new DSM(dsm._id, dsm.title, dsm.description,
			dsm.version, dsm.author, dsm.config,
			dsm.stylesheet, dsm.style, dsm.contributors,
			dsm.idle_backgrounds, dsm.video_select_backgrounds,
			dsm.ce_set, dsm.attributes)
		
		return add(db, dsmColName, dsmObj)
		.then((res, result) => {
			console.log(`Inserted a dsm with _id: ${result.insertedID}`)
			res.status(200)
			client.close()
		})
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

		collection.insertOne(obj, opts).then(resolve).catch(reject);
	})
}

module.exports = router