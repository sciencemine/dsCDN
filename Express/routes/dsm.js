// This route is for defining dsm requests and stuff
const express = require('express'),
        mongo = require('mongodb'),
        MongoClient = mongo.MongoClient,
        mongoURL = 'mongodb://localhost:27017',
        dbName = 'ds';

let router = express.Router()

//This is a test comment

router.route('/dsm/:id')
.get((req, res) => {
	let dsmID = new mongo.ObjectID(req.params.id);

	//get the dsm from mongo database
	MongoClient.connect(mongoURL, (err, client) => {
		if (err) {
			console.error(err)

			res.status(500)

			return
		}

		const db = client.db(dbName)

		//query time! woo!
		query(db, 'dsm', dsmID) 
		.then((doc) => { //I *promise* to come back with things :P
			let dsmObj = doc,
				promises = [];
			
			for (let i = 0; i < dsmObj.ce_set.length; i++ ) {
				let ceId = dsmObj.ce_set[i]
				
			}

		})
		.catch(queryErr) //I broke my promise guys, I'm sorry
	})
})